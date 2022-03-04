package notify_test

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2-utils/notify"
)

func TestNotify(t *testing.T) {
	filename := "/tmp/test.txt"
	dirname := "/tmp"
	flagCreate := false
	flagRemove := false
	flagRename := false
	flagWrite := false
	flagChmod := false
	fx := func(observer *notify.ObserverNotify) {
		if strings.HasSuffix(observer.CurrentEvent.Name, filename) && observer.CurrentEvent.Op.String() == "CREATE" {
			flagCreate = true
		}
		if strings.HasSuffix(observer.CurrentEvent.Name, filename) && observer.CurrentEvent.Op.String() == "REMOVE" {
			flagRemove = true
		}
		if strings.HasSuffix(observer.CurrentEvent.Name, filename) && observer.CurrentEvent.Op.String() == "RENAME" {
			flagRename = true
		}
		if strings.HasSuffix(observer.CurrentEvent.Name, filename) && observer.CurrentEvent.Op.String() == "CHMOD" {
			flagChmod = true
		}
		if strings.HasSuffix(observer.CurrentEvent.Name, filename) && observer.CurrentEvent.Op.String() == "WRITE" {
			flagWrite = true
		}
	}
	notify.NewNotify(dirname, "*").
		FxCreate(fx).
		FxWrite(fx).
		FxChmod(fx).
		FxRemove(fx).
		FxRename(fx).
		Run()

	exec.Command("touch", filename).Run()
	// TODO: Generate Write event on file
	exec.Command("mv", filename, filename+"1").Run()
	exec.Command("mv", filename+"1", filename).Run()
	exec.Command("chmod", "777", filename).Run()
	exec.Command("rm", filename).Run()

	time.Sleep(2 * time.Second)
	if !flagRemove {
		t.Error("file was not removed")
		return
	}
	if !flagCreate {
		t.Error("file was not created")
		return
	}
	if !flagRename {
		t.Error("file was not renamed")
		return
	}
	if !flagChmod {
		t.Error("file was not chmod")
		return
	}
	if !flagWrite {
		// TODO: fix this test
		//t.Error("file was not written")
		//return
	}

}
