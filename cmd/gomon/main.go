package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/arturoeanton/go-r2-utils/commons"
	"github.com/arturoeanton/go-r2-utils/dv"
	"github.com/arturoeanton/go-r2-utils/notify"
)

var (
	config     GoMonConfig
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

type GoMonConfig struct {
	Name     string   `yaml:"name" dv:""`
	Include  []string `yaml:"include" dv:"."`
	Exclude  []string `yaml:"exclude" dv:""`
	Commands struct {
		Command   string `yaml:"command" dv:"go run ."`
		Terminate string `yaml:"terminate" dv:""`
	} `yaml:"commands"`
	Extensions []string `yaml:"extensions" dv:"go,html,js"`
	Log        bool     `yaml:"log" dv:"false"`
	Version    bool     `yaml:"version" dv:"false"`
}

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
		watchOnDir(path)
	} else {
		for _, ext := range extensions {
			if strings.HasSuffix(path, "."+ext) {
				go runCommand(path)
			}
		}
	}
}

func runCommand(filenameOfEvent string) {
	m.Lock()
	defer m.Unlock()
	if *flagLog {
		log.Println("gomom>", "RunCommand", filenameOfEvent)
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
			log.Println("gomom>", err)
		}
	}

	if err := cmd.Start(); err != nil {
		if *flagLog {
			log.Println("gomom>", err)
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
				log.Println(name, "[", len(str), "]>"+str)
			} else {
				fmt.Println(str)
			}

			if err != nil {
				return
			}
		}
	}()
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
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

	return commons.RemoveRepeat(files), err
}

func watchOnDir(currentDir string) {
	if !commons.ContainsString(config.Exclude, currentDir) {
		if *flagLog {
			log.Println("gomom>", "Watch on", currentDir)
		}
		notify.NewNotify(currentDir, "*").FxAll(event).Run()
	} else {
		if *flagLog {
			log.Println("gomom>", "Exclude", currentDir)
		}
	}
}

func main() {

	configFile := flag.String("config", "gomon.yaml", "Config file")
	flagLog = flag.Bool("log", false, "print log")
	v := flag.Bool("v", false, "print version")
	flag.Parse()

	err := dv.Fill(&config)
	if err != nil {
		log.Printf("Fill err   #%v ", err)
		return
	}
	if commons.Exists(*configFile) {
		yamlFile, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
			return
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
			return
		}
		if config.Version {
			*v = true
		}
		if config.Log {
			*flagLog = true
		}
	}
	name = config.Name
	extensions = config.Extensions
	commandArray := strings.Split(config.Commands.Command, " ")
	command = commandArray[0]
	args = []string{}
	if len(commandArray) > 1 {
		args = commandArray[1:]
	}
	terminateCommandArray := strings.Split(config.Commands.Terminate, " ")
	terminateCommand = terminateCommandArray[0]
	terminateArgs = []string{}
	if len(terminateCommandArray) > 1 {
		terminateArgs = terminateCommandArray[1:]
	}

	if *v {
		fmt.Println("gomom> v0.0.3")
	}

	defer func() {
		if terminateCommand != "" {
			tcmd := exec.Command(terminateCommand, terminateArgs...)
			tcmd.Run()
			if *flagLog {
				log.Println("gomom>", "terminateCommand >", terminateCommand, terminateArgs)
			}
		}
		if cmd != nil {
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				if *flagLog {
					log.Println("gomom>", "kill >", cmd.Process.Pid, cmd)
				}
			}
		}
	}()

	done := make(chan struct{})

	go func() {
		if *flagLog {
			log.Println("gomom>", "Listening signals...")
		}
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		close(done)
	}()

	config.Include = commons.RemoveRepeat(config.Include)
	config.Exclude = commons.RemoveRepeat(config.Exclude)
	directories := make([]string, 0)
	for _, directory := range config.Include {
		directoriesTmp, err := FilePathWalkDir(directory)
		if err != nil {
			panic(err)
		}
		directories = append(directories, directoriesTmp...)
	}

	directories = commons.RemoveRepeat(directories)

	for _, currentDir := range directories {

		watchOnDir(currentDir)

	}
	runCommand("[start]")

	<-done
}
