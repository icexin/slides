package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

func handShake(conn net.Conn) error {
	var version int8
	binary.Read(conn, binary.BigEndian, &version)

	var methodlen int8
	binary.Read(conn, binary.BigEndian, &methodlen)

	method := make([]byte, methodlen)
	io.ReadFull(conn, method)

	if version != 5 {
		return errors.New("bad version")
	}

	resp := []byte{5, 0}
	conn.Write(resp)

	return nil
}

func readAddr(conn net.Conn, addrtype int8) (string, error) {
	switch addrtype {
	case 1:
		buf := make([]byte, 4)
		io.ReadFull(conn, buf)
		return net.IP(buf).String(), nil
	case 4:
		buf := make([]byte, 16)
		io.ReadFull(conn, buf)
		return net.IP(buf).String(), nil
	case 3:
		var domainlen int8
		binary.Read(conn, binary.BigEndian, &domainlen)
		buf := make([]byte, domainlen)
		io.ReadFull(conn, buf)
		return string(buf), nil
	}
	return "", errors.New("bad type")
}

func parseRequest(conn net.Conn) (string, int16, error) {
	var version int8
	binary.Read(conn, binary.BigEndian, &version)

	var cmd int8
	binary.Read(conn, binary.BigEndian, &cmd)

	var reserve int8
	binary.Read(conn, binary.BigEndian, &reserve)

	var addrtype int8
	binary.Read(conn, binary.BigEndian, &addrtype)

	addr, err := readAddr(conn, addrtype)
	if err != nil {
		return "", 0, err
	}

	var port int16
	binary.Read(conn, binary.BigEndian, &port)

	if version != 5 {
		return "", 0, errors.New("bad version")
	}

	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	return addr, port, nil
}

func passwordToKey(password string) []byte {
	md5sum := md5.Sum([]byte(password))
	return md5sum[:]
}

type cryptoConn struct {
	key  []byte
	iv   []byte
	conn net.Conn
	dec  cipher.Stream
	enc  cipher.Stream
}

func NewCryptoConn(conn net.Conn, password string) *cryptoConn {
	return &cryptoConn{
		conn: conn,
		key:  passwordToKey(password),
	}
}

func (c *cryptoConn) initEnc() {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		panic(err)
	}

	if c.iv == nil {
		c.iv = make([]byte, 16)
		io.ReadFull(rand.Reader, c.iv)
	}
	c.conn.Write(c.iv)
	c.enc = cipher.NewCFBEncrypter(block, c.iv)
}

func (c *cryptoConn) initDec() {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		panic(err)
	}

	iv := make([]byte, 16)
	io.ReadFull(c.conn, iv)
	if c.iv == nil {
		c.iv = iv
	}
	if !bytes.Equal(iv, c.iv) {
		panic("")
	}

	c.dec = cipher.NewCFBDecrypter(block, c.iv)
}

func (c *cryptoConn) Write(b []byte) (int, error) {
	if c.enc == nil {
		c.initEnc()
	}
	buf := make([]byte, len(b))
	c.enc.XORKeyStream(buf, b)
	return c.conn.Write(buf)
}

func (c *cryptoConn) Read(b []byte) (int, error) {
	if c.dec == nil {
		c.initDec()
	}
	n, err := c.conn.Read(b)
	c.dec.XORKeyStream(b[:n], b[:n])
	return n, err
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	err := handShake(conn)
	if err != nil {
		log.Print(err)
		return
	}

	addr, port, err := parseRequest(conn)
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(addr)
	remote, err := net.Dial("tcp", "")
	if err != nil {
		log.Print(err)
		return
	}
	defer remote.Close()

	rremote := NewCryptoConn(remote, "")

	addrreq := []byte{3, byte(len(addr))}
	addrreq = append(addrreq, []byte(addr)...)
	rremote.Write(addrreq)
	binary.Write(rremote, binary.BigEndian, port)

	go io.Copy(rremote, conn)
	io.Copy(conn, rremote)
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
