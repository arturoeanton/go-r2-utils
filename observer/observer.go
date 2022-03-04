package observer

import (
	"fmt"

	"github.com/google/uuid"
)

type Observer interface {
	GetID() string
	Notify(data interface{})
}

type BasicObserver struct {
	id string
}

func (e *BasicObserver) GetID() string {
	if e.id == "" {
		e.id = uuid.New().String()
	}
	return e.id
}
func (e *BasicObserver) Notify(data interface{}) {
	fmt.Println("Basic Observer see   State:", data.(string))
}
