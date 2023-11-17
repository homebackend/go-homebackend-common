package homecommon

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

type Nothing int

type Pid int

func (s *Pid) GetStatus(args *Nothing, pid *int) error {
	p := os.Getpid()
	*pid = p
	return nil
}

func socketFile(progName string) string {
	return fmt.Sprintf("/tmp/%s.sock", progName)
}

func delSocketFile(progName string) error {
	sockFile := socketFile(progName)
	if _, err := os.Stat(sockFile); err == nil {
		log.Printf("Removing existing socket file: %s", sockFile)
		if e := os.Remove(sockFile); err != nil {
			return e
		}
	}

	return nil
}

func StartIpc(progName string, rcvrs ...any) error {
	if err := delSocketFile(progName); err != nil {
		return err
	}

	p := new(Pid)
	rpc.Register(p)

	for _, rcvr := range rcvrs {
		rpc.Register(rcvr)
	}

	l, err := net.Listen("unix", socketFile(progName))
	if err != nil {
		log.Printf("Error during listening: %s", err)
		return err
	}

	// Setting permission to 777 so that non root users can also access
	os.Chmod(socketFile(progName), 0777)

	go rpc.Accept(l)

	return nil
}

func StopIpc(progName string) error {
	return delSocketFile(progName)
}

func IpcGetData[T any](progName string, funcName string, args any) (T, error) {
	var reply T
	client, err := rpc.Dial("unix", socketFile(progName))
	if err != nil {
		log.Printf("Error connecting to UNIX socket: %s", err)
		return reply, err
	}

	err = client.Call(funcName, args, &reply)
	if err != nil {
		log.Printf("Error while getting status: %s", err)
		return reply, err
	}

	return reply, nil
}

func IpcGetStatus(progName string) (int, error) {
	return IpcGetData[int](progName, "Pid.GetStatus", 0)
}
