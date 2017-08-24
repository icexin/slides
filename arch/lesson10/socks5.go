package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func handshake(r *bufio.Reader, conn net.Conn) error {
	// handshake
	version, _ := r.ReadByte()
	if version != 5 {
		return errors.New("bad version")
	}
	methodlen, _ := r.ReadByte()
	method := make([]byte, methodlen)
	io.ReadFull(r, method)

	resp := []byte{5, 0}
	conn.Write(resp)

	return nil
}

func readAddr(r *bufio.Reader) (string, error) {
	// version
	r.ReadByte()
	// cmd
	r.ReadByte()
	// reserve
	r.ReadByte()
	addrtype, _ := r.ReadByte()
	if addrtype != 3 {
		return "", errors.New("bad addrtype")
	}

	addrlen, _ := r.ReadByte()
	buf := make([]byte, addrlen)
	io.ReadFull(r, buf)
	var port int16
	binary.Read(r, binary.BigEndian, &port)
	return fmt.Sprintf("%s:%d", buf, port), nil
}

func doproxy(r *bufio.Reader, conn net.Conn, addr string) error {
	remote, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go io.Copy(remote, r)
	io.Copy(conn, remote)
	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	err := handshake(r, conn)
	if err != nil {
		log.Print(err)
		return
	}

	addr, err := readAddr(r)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(addr)
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	err = doproxy(r, conn, addr)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
