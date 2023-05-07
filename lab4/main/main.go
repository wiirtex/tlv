package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func parent() {
	pid := os.Getpid()
	fmt.Println("Parent process ID:", pid)

	attr := &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	attr.Sys = &syscall.SysProcAttr{}
	attr.Sys.Cloneflags = syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWUTS
	childPID, err := os.StartProcess("/bin/sh", []string{"", "test.sh"}, attr)
	if childPID == nil || err != nil {
		if err != nil {
			panic(err)
		}
		fmt.Println("Error: could not fork process: ", err)
		return
	}

	cmd := exec.Command("ip", strings.Split(fmt.Sprintf("link add name veth0 type veth peer name veth1 netns %d", childPID.Pid), " ")...)
	cmd.Run()

	fmt.Println("Child process ID:", childPID.Pid)

	cmd = exec.Command("ifconfig", strings.Split("veth0 10.1.1.1/24 up", " ")...)
	cmd.Run()

	_, err = childPID.Wait()
	if err != nil {
		fmt.Println("Error waiting for child process to finish:", err)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Starting root")
		parent()
	} else {
		panic("how on Earth is this possible")
	}
}
