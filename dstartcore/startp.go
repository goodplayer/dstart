package dstartcore

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func startFromConfig(c *Config) error {
	if c.Global.Envs == nil {
		c.Global.Envs = make(map[string]string)
	}

	log.Println("verifying all services")
	for s, v := range c.Services {
		if err := VerifyService(v); err != nil {
			log.Println("verify service:", s, "failed.", err)
			return err
		}
	}

	log.Println("starting all services")
	for s, v := range c.Services {
		if err := StartService(v, c.Global.Envs); err != nil {
			log.Println("start service:", s, "failed.", err)
			return err
		}
		log.Println("service:", s, "started.")
	}

	return nil
}

func VerifyService(service Service) error {
	//FIXME
	return nil
}

func StartService(service Service, globalEnv map[string]string) error {

	attr := new(os.ProcAttr)

	// working directory
	if service.WorkingDirectory != "" {
		attr.Dir = service.WorkingDirectory
	}

	// envs
	var envs []string
	for k, v := range globalEnv {
		envs = append(envs, k+"="+v)
	}
	for k, v := range service.Envs {
		envs = append(envs, k+"="+v)
	}
	if len(service.Envs) > 0 {
		attr.Env = append(os.Environ(), envs...)
	} else {
		attr.Env = append(os.Environ())
	}

	// input file
	in := os.Stdin
	if service.InputFile != "" {
		f, err := os.Open(service.InputFile)
		if err != nil {
			log.Println("open input file error: ", err)
			return errors.New("open input file error")
		}
		in = f
	}

	// output file
	out := os.Stdout
	if service.OutputFile != "" {
		fp := os.O_WRONLY | os.O_CREATE
		if service.OutputFileAppend {
			fp |= os.O_APPEND
		} else {
			fp |= os.O_TRUNC
		}
		f, err := os.OpenFile(service.OutputFile, fp, 0666)
		if err != nil {
			log.Println("open output file error: ", err)
			return errors.New("open output file error")
		}
		out = f
	}

	// error output file
	erroutput := os.Stderr
	if service.RedirectErrorOutputToOutput {
		erroutput = out
	} else {
		if service.ErrorOutputFile != "" {
			fp := os.O_WRONLY | os.O_CREATE
			if service.ErrorOutputFileAppend {
				fp |= os.O_APPEND
			} else {
				fp |= os.O_TRUNC
			}
			f, err := os.OpenFile(service.ErrorOutputFile, fp, 0666)
			if err != nil {
				log.Println("open error output file error: ", err)
				return errors.New("open error output file error")
			}
			erroutput = f
		}
	}

	attr.Files = []*os.File{in, out, erroutput} // input/output/error

	// user
	if service.Username != "" {
		u, err := user.Lookup(service.Username)
		if err != nil {
			log.Println(err)
			return errors.New("lookup user error")
		}
		if attr.Sys == nil {
			attr.Sys = new(syscall.SysProcAttr)
			attr.Sys.Credential = new(syscall.Credential)
		} else {
			if attr.Sys.Credential == nil {
				attr.Sys.Credential = new(syscall.Credential)
			}
		}
		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			log.Println("strconv uid error: ", err)
			return errors.New("convert uid error")
		}
		gid, err := strconv.Atoi(u.Gid)
		if err != nil {
			log.Println("strconv gid error: ", err)
			return errors.New("convert gid error")
		}
		attr.Sys.Credential.Uid = uint32(uid)
		attr.Sys.Credential.Gid = uint32(gid)
	}
	// group
	if service.GroupName != "" {
		g, err := user.LookupGroup(service.GroupName)
		if err != nil {
			log.Println(err)
			return errors.New("lookup group error")
		}
		if attr.Sys == nil {
			attr.Sys = new(syscall.SysProcAttr)
			attr.Sys.Credential = new(syscall.Credential)
		} else {
			if attr.Sys.Credential == nil {
				attr.Sys.Credential = new(syscall.Credential)
			}
		}
		gid, err := strconv.Atoi(g.Gid)
		if err != nil {
			log.Println("strconv gid error: ", err)
			return errors.New("convert gid error")
		}
		attr.Sys.Credential.Gid = uint32(gid)
	}

	// exec
	cmds := strings.Split(strings.TrimSpace(service.ExecCommand), " ")
	var args []string
	if len(cmds) > 1 {
		args = cmds[1:]
	}
	p, err := os.StartProcess(cmds[0], args, attr)
	signal.Ignore(syscall.SIGCHLD)

	if err != nil {
		log.Println("start process:[", cmds[0], "] error: ", err)
		return errors.New("start process error")
	}

	log.Println("start process:[", cmds[0], "] pid: ", p.Pid)

	return nil
}
