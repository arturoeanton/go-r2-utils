package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/arturoeanton/go-r2-utils/notify"
)

var (
	dir        string
	extensions []string
	m          = sync.Mutex{}
	command    string
	args       []string
	cmd        *exec.Cmd
)

func event(observer *notify.ObserverNotify) {
	path := observer.CurrentEvent.Name
	log.Println(observer.CurrentEvent.String())
	log.Println(path)

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	if fileInfo.IsDir() {
		notify.NewNotify(path, "*").FxAll(event).Run()
	} else {
		for _, ext := range extensions {
			if strings.HasSuffix(path, "."+ext) {
				go runCommand()
			}
		}
	}
}

func runCommand() {
	m.Lock()
	defer m.Unlock()

	if cmd != nil {
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	}

	cmd = exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	go func() {
		for {
			tmp := make([]byte, 1024)
			_, err := stdout.Read(tmp)
			fmt.Print(string(tmp))
			if err != nil {
				return
			}
		}
	}()
}

func main() {
	c := make(chan bool)
	dir = "."
	command = "go"
	args = []string{"run", dir}
	extensions = []string{"go"}

	files, err := FilePathWalkDir(dir)
	if err != nil {
		panic(err)
	}

	for _, currentDir := range files {
		log.Println("Watch on", currentDir)
		notify.NewNotify(currentDir, "*").FxAll(event).Run()
	}
	runCommand()

	<-c
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			files = append(files, path)
			if path != root {
				subFiles, err := FilePathWalkDir(path)
				if err != nil {
					log.Panic(err)
				}
				files = append(files, subFiles...)
			}
		}
		return nil
	})

	m := make(map[string]bool)
	for _, file := range files {
		m[file] = true
	}
	setFiles := make([]string, 0)
	for k := range m {
		setFiles = append(setFiles, k)
	}

	return setFiles, err
}
