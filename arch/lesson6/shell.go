package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	line := "cat shell.go|grep main"
	cmds := strings.Split(line, "|")
	s1 := strings.Fields(cmds[0])
	s2 := strings.Fields(cmds[1])
	cmd1 := exec.Command(s1[0], s1[1:]...)
	cmd1.Stdin = os.Stdin
	out, _ := cmd1.StdoutPipe()
	cmd2 := exec.Command(s2[0], s2[1:]...)
	cmd2.Stdin = out
	cmd2.Stdout = os.Stdout
	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	cmd2.Wait()
}
