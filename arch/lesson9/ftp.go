package main

import (
	"log"
	"net"
)

func parseArgs(conn net.Conn) ([]string, error) {
	return nil, nil
}

func doAction(conn net.Conn, args []string) error {
	return nil
}

func writeError(conn net.Conn, err error) {
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// 获取命令行参数
	// 执行动作获取结果
	// 返回结果
	args, err := parseArgs(conn)
	if err != nil {
		log.Print(err)
		writeError(conn, err)
		return
	}

	err = doAction(conn, args)
	if err != nil {
		writeError(conn, err)
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
