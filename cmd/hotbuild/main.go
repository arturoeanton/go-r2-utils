package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/arturoeanton/go-r2-utils/notify"
)

var (
	dir        *string
	extensions []string
	m          = sync.Mutex{}
	command    string
	args       []string

	terminateCommand string
	terminateArgs    []string
	cmd              *exec.Cmd
)

func event(observer *notify.ObserverNotify) {
	path := observer.CurrentEvent.Name

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
	log.Println("runCommand")

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

	defer func() {
		if terminateCommand != "" {
			tcmd := exec.Command(terminateCommand, terminateArgs...)
			tcmd.Run()
			log.Println("hotbuild>", "terminateCommand >", terminateCommand, terminateArgs)
		}
	}()

	done := make(chan struct{})

	go func() {
		log.Println("Listening signals...")
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		close(done)
	}()

	dir = flag.String("dir", ".", "default .")
	extensionsParam := flag.String("ext", "go", "default go")
	commandParam := flag.String("cmd", "go run .", "default \"go run .\"")
	terminateCommandParam := flag.String("end", "", "example \"killall hello\"")
	flag.Parse()

	commandArray := strings.Split(*commandParam, " ")
	terminateCommandArray := strings.Split(*terminateCommandParam, " ")
	command = commandArray[0]
	args = []string{}
	if len(commandArray) > 1 {
		args = commandArray[1:]
	}
	terminateCommand = terminateCommandArray[0]
	terminateArgs = []string{}
	if len(terminateCommandArray) > 1 {
		terminateArgs = terminateCommandArray[1:]
	}
	extensions = strings.Split(*extensionsParam, ",")

	files, err := FilePathWalkDir(*dir)
	if err != nil {
		panic(err)
	}

	for _, currentDir := range files {
		log.Println("Watch on", currentDir)
		notify.NewNotify(currentDir, "*").FxAll(event).Run()
	}
	runCommand()

	<-done
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
