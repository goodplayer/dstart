package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// user
	// param
	// env
	// input/output/error
	// exec

	fmt.Println(user.Lookup("root"))
	u, _ := user.Lookup("root")

	attr := new(os.ProcAttr)
	attr.Dir = "/tmp"
	attr.Env = []string{} // env
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr} // input/output/error
	attr.Sys = new(syscall.SysProcAttr)
	attr.Sys.Credential = new(syscall.Credential)
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)
	attr.Sys.Credential.Uid = uint32(uid)
	attr.Sys.Credential.Gid = uint32(gid)

	p, err := os.StartProcess("go-import-server", []string{"20"}, attr)

	fmt.Println(p, err)

	time.Sleep(50 * time.Second)
}
