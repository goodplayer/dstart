package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"
)

var username string
var wd string
var input string
var output string
var erroroutput string
var etoo bool

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "array flags."
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var env arrayFlags

func init() {
	flag.StringVar(&username, "u", "", "-u username")
	flag.StringVar(&wd, "wd", "", "-wd working_directory")
	flag.Var(&env, "env", "-env key=value [-env key=value ...]")
	flag.StringVar(&input, "in", "", "-in input_file")
	flag.StringVar(&output, "out", "", "-out output_file (currently overwriten file not create a new one)")
	flag.StringVar(&erroroutput, "err", "", "-err error_output_file (currently overwriten file not create a new one)")
	flag.BoolVar(&etoo, "etoo", false, "-etoo : redirect error output to standard output")
}

// dstart [options] executable [params...]
func main() {
	flag.Parse()

	if len(flag.Args()) <= 0 {
		log.Fatalln("should have executable file. illegal argument. usage: dstart [options] executable [params...]")
	}

	attr := new(os.ProcAttr)
	attr.Dir = wd
	if len(env) > 0 {
		attr.Env = append(os.Environ(), env...)
	} else {
		attr.Env = append(os.Environ())
	}

	in := os.Stdin
	if input != "" {
		f, err := os.Open(input)
		if err != nil {
			log.Fatalln("open input file error: ", err)
		}
		in = f
	}

	out := os.Stdout
	if output != "" {
		f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalln("open output file error: ", err)
		}
		out = f
	}

	erroutput := os.Stderr
	if erroroutput != "" {
		f, err := os.OpenFile(erroroutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalln("open error output file error: ", err)
		}
		erroutput = f
	}

	if etoo {
		erroutput = out
	}

	attr.Files = []*os.File{in, out, erroutput} // input/output/error

	// user
	if username != "" {
		u, err := user.Lookup(username)
		if err != nil {
			log.Fatalln(err)
		}
		attr.Sys = new(syscall.SysProcAttr)
		attr.Sys.Credential = new(syscall.Credential)
		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			log.Fatalln("strconv uid error: ", err)
		}
		gid, err := strconv.Atoi(u.Gid)
		if err != nil {
			log.Fatalln("strconv gid error: ", err)
		}
		attr.Sys.Credential.Uid = uint32(uid)
		attr.Sys.Credential.Gid = uint32(gid)
	}

	p, err := os.StartProcess(flag.Args()[0], flag.Args(), attr)
	signal.Ignore(syscall.SIGCHLD)

	if err != nil {
		log.Fatalln("start process error: ", err)
	} else {
		log.Println("start process. pid: ", p.Pid)
	}
}
