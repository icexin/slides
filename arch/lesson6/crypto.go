package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Xor struct {
	w io.Writer
	x byte
}

func (x *Xor) Write(p []byte) (int, error) {
	p1 := make([]byte, len(p))
	copy(p1, p)
	for i, b := range p1 {
		p1[i] = b ^ x.x
	}
	return x.w.Write(p1)
}

func NewXor(w io.Writer, x byte) *Xor {
	return &Xor{
		w: w,
		x: x,
	}
}

func main() {
	buf := new(bytes.Buffer)
	x := NewXor(buf, 'a')
	io.WriteString(x, "hello")
	fmt.Println(buf.Bytes())

	x1 := NewXor(os.Stdout, 'a')
	io.WriteString(x1, buf.String())
}
