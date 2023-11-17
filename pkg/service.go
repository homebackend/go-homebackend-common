package homecommon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"runtime"
	"strings"
	"syscall"

	"github.com/juju/fslock"
)

const (
	O_ANY          = 0
	O_LINUX        = 1
	USER_RUN_DIR   = "/run/user"
	GLOBAL_RUN_DIR = "/var/run"
)

func CreatePidFile() *fslock.Lock {
	var pidFilePath string

	if u, err := user.Current(); err == nil {
		if u.Uid != "0" {
			pidFilePath = fmt.Sprintf("%s/%s/goifs.pid", USER_RUN_DIR, u.Uid)
		} else {
			pidFilePath = fmt.Sprintf("%s/goifs.pid", GLOBAL_RUN_DIR)
		}
	} else {
		log.Fatalf("Unable to determine the current user: %v", err)
	}

	lock := fslock.New(pidFilePath)

	if err := lock.TryLock(); err != nil {
		log.Fatalf("An instance of the internet failover service is already running. If you are sure that is not the case please delete the file: %s", pidFilePath)
	}

	return lock
}

func CheckPrerequisites(flag int, confFilePath string, requiredCommands []string) {
	if flag&O_LINUX == 1 {
		if !strings.Contains(runtime.GOOS, "linux") {
			log.Fatalf("Unsupported OS: %s", runtime.GOOS)
		}
	}

	for _, c := range requiredCommands {
		if _, err := exec.LookPath(c); err != nil {
			log.Fatalf("Required command not found in PATH: %s", c)
		}
	}

	if confFilePath != "" {
		if _, err := os.Stat(confFilePath); err != nil {
			log.Fatalf("Unable to load configuration from file: %s : %v", confFilePath, err)
		}
	}
}

func GetPid(progName string) int {
	pid, err := IpcGetStatus(progName)
	if err != nil {
		log.Fatal("Error communicating with goifs. Are you sure it is running?")
	}

	return pid
}

func Signal() chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	return sigc
}

func Stop(progName string) {
	pid := GetPid(progName)

	if p, err := os.FindProcess(pid); err != nil {
		log.Fatalf("Process does not exist: %s", err)
	} else {
		if err := p.Signal(os.Interrupt); err != nil {
			log.Fatalf("Error sending signal: %s", err)
		}
	}
}
