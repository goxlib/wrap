package os

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sh "github.com/codeskyblue/go-sh"
)

// GetPidByName returns app's pid by given name,
func GetPidByName(name string) (pid int, err error) {
	err = filepath.Walk("/proc", func(path string, info os.FileInfo, err error) error {
		// using the truth that file `/proc/[pid]/status` store information for
		// a running app with pid
		if strings.Contains(path, "/status") && strings.Count(path, "/") == 3 {
			f, err := os.Open(path)
			if err != nil {
				return err // maybe just return a nil
			}
			defer f.Close()

			rd := bufio.NewReader(f)
			// the first line contains name, so just read this line
			for i := 0; i < 1; i++ {
				line, err := rd.ReadString('\n')
				if err != nil || io.EOF == err {
					break
				}
				if strings.Contains(line, name) {
					// get the pid from file path
					target := strings.Split(path, "/")[2]
					var err error
					pid, err = strconv.Atoi(target)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		return nil
	})

	return pid, err
}

// GetNameByPid returns app's name by given pid.
func GetNameByPid(pid int) (name string, err error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return name, fmt.Errorf("no app running with pid[%d] found", pid)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	// the first line contains name, so just read this line
	for i := 0; i < 1; i++ {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			return "", err
		}
		if strings.Contains(line, "Name") {
			target := strings.Split(line, ":")[1]
			return strings.TrimSpace(target), nil
		}
	}
	return name, fmt.Errorf("no app running with pid[%d] found", pid)
}

// IsRunningByPort check whether a server is running by port.
func IsRunningByPort(port int) bool {
	o, err := sh.Command("lsof", fmt.Sprintf("-i:%d", port)).
		Command("wc", "-l").
		Output()
	if err != nil || strings.TrimSpace(string(o)) == "0" {
		return false
	}

	return true
}

// GetPidByPort returns app's pid by its port.
func GetPidByPort(port int) (pid string) {
	// lsof -i:5050 |grep -v PID| awk '{print $2}'
	o, err := sh.Command("lsof", fmt.Sprintf("-i:%d", port)).
		Command("grep", "-v", "PID").
		Command("awk", `{print $2}`).
		Output()
	if err != nil {
		return
	}

	return strings.TrimSpace(string(o))
}
