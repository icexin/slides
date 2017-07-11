package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/rakyll/pb"
)

var (
	works = flag.Int("n", 5, "works")
)

type Block struct {
	Id    int
	Begin int64
	End   int64
	Name  string
}

func GetSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.ContentLength, nil
}

func GenBlocks(total int64, works int) []*Block {
	var blocks []*Block
	n := int64(math.Ceil(float64(total) / float64(works)))
	for i := int64(0); i < int64(works)-1; i++ {
		block := &Block{
			Id:    int(i),
			Begin: i * n,
			End:   (i + 1) * n,
		}
		blocks = append(blocks, block)
	}
	blocks = append(blocks, &Block{
		Id:    works - 1,
		Begin: n * (int64(works) - 1),
		End:   total,
	})
	return blocks
}

func Download(url string, b *Block, bar *pb.ProgressBar) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", b.Begin, b.End-1))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	name := fmt.Sprintf("%s.%d", path.Base(url), b.Id)
	b.Name = name
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	w := io.MultiWriter(bar, f)
	_, err = io.Copy(w, resp.Body)
	return err
}

func Merge(name string, blocks []*Block) error {
	readers := make([]io.Reader, len(blocks))
	for i, b := range blocks {
		f, err := os.Open(b.Name)
		if err != nil {
			return err
		}
		defer f.Close()
		readers[i] = f
	}
	r := io.MultiReader(readers...)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	for _, b := range blocks {
		os.Remove(b.Name)
	}
	return nil
}

func main() {
	flag.Parse()
	url := flag.Arg(0)
	total, err := GetSize(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("total:%d", total)

	blocks := GenBlocks(total, *works)

	bar := pb.New(int(total)).SetUnits(pb.U_BYTES)
	bar.Start()

	group := new(sync.WaitGroup)
	for _, b := range blocks {
		log.Printf("%v", b)
		b := b
		group.Add(1)
		go func() {
			defer group.Done()
			err := Download(url, b, bar)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	group.Wait()
	bar.Finish()

	log.Printf("group done")
	err = Merge(path.Base(url), blocks)
	if err != nil {
		log.Fatal(err)
	}
}
