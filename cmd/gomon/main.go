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
	name       string
	args       []string
	flagLog    *bool

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
	if *flagLog {
		log.Println("hotbuild>", "RunCommand")
	}

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
		if *flagLog {
			log.Println("hotbuild>", err)
		}
	}

	if err := cmd.Start(); err != nil {
		if *flagLog {
			log.Println("hotbuild>", err)
		}
	}

	go func() {
		for {
			tmp := make([]byte, 1024*1024)
			i, err := stdout.Read(tmp)
			str := string(tmp)[:i]

			if i == 0 {
				continue
			}

			str = strings.TrimSuffix(str, "\n")

			if *flagLog {
				fmt.Println("build>", name, "><", len(str), ">"+str)
			} else {
				fmt.Println(str)
			}

			if err != nil {
				return
			}
		}
	}()
}

// gomon -cmd "go run ." -dir . -ext go,html,js -log
// gomon  -dir .
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

func main() {

	dir = flag.String("dir", ".", "default .")
	flagLog = flag.Bool("log", false, "print log")
	ov := flag.Bool("ov", false, "only version")
	v := flag.Bool("v", false, "print version")
	nameParam := flag.String("name", "", "Example example1")
	extensionsParam := flag.String("ext", "go", "default go")
	commandParam := flag.String("cmd", "go run .", "default \"go run .\"")
	terminateCommandParam := flag.String("end", "", "example \"killall hello\"")
	flag.Parse()

	if *ov {
		fmt.Println("hotbuild> 1.0.0")
		return
	}
	if *v {
		fmt.Println("hotbuild> 1.0.0")
	}

	defer func() {
		if terminateCommand != "" {
			tcmd := exec.Command(terminateCommand, terminateArgs...)
			tcmd.Run()
			if *flagLog {
				log.Println("hotbuild>", "terminateCommand >", terminateCommand, terminateArgs)
			}
		}
		if cmd != nil {
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				if *flagLog {
					log.Println("hotbuild>", "kill >", cmd.Process.Pid, cmd)
				}
			}
		}
	}()

	done := make(chan struct{})

	go func() {
		if *flagLog {
			log.Println("hotbuild>", "Listening signals...")
		}
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		close(done)
	}()

	name = *nameParam
	if *nameParam == "" {
		name, _ = filepath.Abs(*dir)
	}

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

		if *flagLog {
			log.Println("hotbuild>", "Watch on", currentDir)
		}

		notify.NewNotify(currentDir, "*").FxAll(event).Run()
	}
	runCommand()

	<-done
}
