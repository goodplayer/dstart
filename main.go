package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func main() {
	fmt.Println(user.Lookup("root"))
	u, _ := user.Lookup("root")

	attr := new(os.ProcAttr)
	attr.Dir = "/sdf/sdf"
	attr.Env = []string{} // env
	//attr.Files = []*os.File{} // input/output/error
	attr.Sys = new(syscall.SysProcAttr)
	attr.Sys.Credential = new(syscall.Credential)
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)
	attr.Sys.Credential.Uid = uint32(uid)
	attr.Sys.Credential.Gid = uint32(gid)

	p, err := os.StartProcess("xxx", []string{"xxx"}, attr)

	fmt.Println(p, err)
}
