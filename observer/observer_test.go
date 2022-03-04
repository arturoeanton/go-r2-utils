package observer_test

import (
	"fmt"
	"testing"

	"github.com/arturoeanton/go-r2-utils/observer"
)

type ElementObservable struct {
	observer.BasicObservable
}

type ElementObserver struct {
	observer.BasicObserver
	Data string
}

func (e *ElementObserver) Notify(data interface{}) {
	e.Data = data.(string)
	fmt.Println("State1:", data.(string))
}

func TestObserver(t *testing.T) {
	observable1 := ElementObservable{}

	observer1 := ElementObserver{}
	observer1.GetID()

	observable1.AddObserver(&observer1)

	observable1.ChangeState("hola")
	if observer1.Data != "hola" {
		t.Error("Observer failed")
	}

	observable1.ChangeState("hola1")

	if observer1.Data != "hola1" {
		t.Error("Observer failed")
	}

}
