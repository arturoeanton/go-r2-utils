package notify

import (
	"log"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type ObserverNotify struct {
	Filename     string
	Directory    string
	Watcher      *fsnotify.Watcher
	CurrentEvent *fsnotify.Event
	fxWrite      func(observer *ObserverNotify)
	fxCreate     func(observer *ObserverNotify)
	fxRemove     func(observer *ObserverNotify)
	fxRename     func(observer *ObserverNotify)
	fxChmod      func(observer *ObserverNotify)
}

func (o *ObserverNotify) FxCreate(fxCreate func(observer *ObserverNotify)) *ObserverNotify {
	o.fxCreate = fxCreate
	return o
}
func (o *ObserverNotify) FxWrite(fxWrite func(observer *ObserverNotify)) *ObserverNotify {
	o.fxWrite = fxWrite
	return o
}
func (o *ObserverNotify) FxRemove(fxRemove func(observer *ObserverNotify)) *ObserverNotify {
	o.fxRemove = fxRemove
	return o
}
func (o *ObserverNotify) FxRename(fxRename func(observer *ObserverNotify)) *ObserverNotify {
	o.fxRename = fxRename
	return o
}
func (o *ObserverNotify) FxChmod(fxChmod func(observer *ObserverNotify)) *ObserverNotify {
	o.fxChmod = fxChmod
	return o
}

func NewNotify(directory string, filename string) *ObserverNotify {
	observer := &ObserverNotify{
		Filename:  filename,
		Directory: directory,
		fxWrite:   func(observer *ObserverNotify) {},
		fxCreate:  func(observer *ObserverNotify) {},
		fxRemove:  func(observer *ObserverNotify) {},
		fxRename:  func(observer *ObserverNotify) {},
		fxChmod:   func(observer *ObserverNotify) {},
	}
	return observer
}

func (o *ObserverNotify) Run() {
	wg := &sync.WaitGroup{}
	go func() {
		var err error
		o.Watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer o.Watcher.Close()

		err = o.Watcher.Add(o.Directory)
		if err != nil {
			log.Fatal(err)
		}

		done := make(chan bool)
		go func() {
			for {
				wg.Add(1)
				select {
				case event, ok := <-o.Watcher.Events:
					if !ok {
						return
					}
					if !strings.HasSuffix(event.Name, o.Filename) && o.Filename != "*" {
						continue
					}
					o.CurrentEvent = &event
					switch {
					case event.Op&fsnotify.Write == fsnotify.Write:
						o.fxWrite(o)
					case event.Op&fsnotify.Create == fsnotify.Create:
						o.fxCreate(o)
					case event.Op&fsnotify.Remove == fsnotify.Remove:
						o.fxRemove(o)
					case event.Op&fsnotify.Rename == fsnotify.Rename:
						o.fxRename(o)
					case event.Op&fsnotify.Chmod == fsnotify.Chmod:
						o.fxChmod(o)
					}
				case err, ok := <-o.Watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
				wg.Done()
			}
		}()
		<-done
	}()
}
