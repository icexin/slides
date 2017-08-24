package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	globalRoom = NewRoom()
)

type User struct {
	id   string
	conn net.Conn
	r    *bufio.Reader
}

func NewUser(conn net.Conn) *User {
	return &User{
		conn: conn,
		r:    bufio.NewReader(conn),
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) HandShake() error {
	msg, err := u.ReadMsg()
	if err != nil {
		return err
	}
	fields := strings.Fields(msg)
	if len(fields) != 2 {
		io.WriteString(u.conn, "bad user or password")
		return errors.New("bad user or password")
	}

	name, password := fields[0], fields[1]
	if password != "reboot123" {
		io.WriteString(u.conn, "bad password")
		return errors.New("bad password")
	}
	u.id = name
	return nil
}

func (u *User) ReadMsg() (string, error) {
	msg, err := u.r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(msg), nil
}

func (u *User) SendMsg(s string) {
	io.WriteString(u.conn, s)
}

func (u *User) Offline() {
	u.conn.Close()
}

type Room struct {
	lock  sync.Mutex
	users map[string]*User
}

func NewRoom() *Room {
	return &Room{
		users: make(map[string]*User),
	}
}

func (r *Room) BroadCast(who string, s string) {
	timestr := time.Now().Format("20060102 15:04:05")
	msg := fmt.Sprintf("%s %s: %s\n", timestr, who, s)
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, u := range r.users {
		if u.Id() == who {
			continue
		}
		u.SendMsg(msg)
	}
}

func (r *Room) Join(u *User) {
	r.BroadCast("system", u.Id()+" join room")
	r.lock.Lock()
	defer r.lock.Unlock()
	if u, ok := r.users[u.Id()]; ok {
		u.Offline()
	}
	r.users[u.Id()] = u
}

func (r *Room) Leave(u *User) {
	r.BroadCast("system", u.Id()+" leave room")
	r.lock.Lock()
	defer r.lock.Unlock()
	if u, ok := r.users[u.Id()]; ok {
		u.Offline()
		delete(r.users, u.Id())
	}
}

func handleConn(conn net.Conn) {
	u := NewUser(conn)
	defer u.Offline()

	err := u.HandShake()
	if err != nil {
		log.Print(err)
		return
	}

	globalRoom.Join(u)
	defer globalRoom.Leave(u)

	for {
		msg, err := u.ReadMsg()
		if err != nil {
			log.Print(err)
			return
		}
		globalRoom.BroadCast(u.Id(), msg)
	}
}

var (
	listenAddr = flag.String("addr", ":8080", "listen address")
)

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", *listenAddr)
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
