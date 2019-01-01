package main

import (
	"os"
	"strings"

	"github.com/BaichuanWu/myDocker/container"
	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, cmds []string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("new parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Fatal(err)
	}
	sendInitCommand(cmds, writePipe)
	parent.Wait()
	os.Exit(0)
}

func sendInitCommand(cmds []string, writePipe *os.File) {
	cmd := strings.Join(cmds, " ")
	log.Infof("commands all is %s", cmd)
	writePipe.WriteString(cmd)
	writePipe.Close()
}
