package shellout

import (
	"fmt"
	"github.com/golang/glog"
	"os/exec"
)

const RETRIES = 10

func Cmd(cmd string, shell bool) (string, error) {
	var out []byte
	var err error
	if shell {
		if out, err = exec.Command("bash", "-c", cmd).Output(); err != nil {
			return "", err
		}
	} else {
		if out, err = exec.Command(cmd).Output(); err != nil {
			return "", err
		}
	}
	return string(out), nil
}

func RetryCmd(cmd string) string {
	var err error
	var out string
	for i := 0; i < RETRIES; i++ {
		if out, err = Cmd(cmd, true); err == nil {
			return out
		}
		glog.Errorf("Error executing %v", cmd)
	}
	glog.Infof(out)
	panic(err)
}
