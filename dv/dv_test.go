package dv_test

import (
	"testing"

	"github.com/arturoeanton/go-r2-utils/dv"
)

type User struct {
	Name    string   `dv:"Pedro"`
	Email   string   `dv:"pedro@test.com"`
	Age     int      `dv:"17"`
	Salary  *float32 `dv:"1234.56"`
	Size    float64  `dv:"1.2"`
	Address struct {
		Calle     string  `dv:"Av Caseros"`
		Direccion int     `dv:"1001"`
		Dpto      *string `dv:"2 B"`
	}
}

func TestFill(t *testing.T) {
	user := User{}
	err := dv.Fill(&user)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if user.Name != "Pedro" {
		t.Error("Name is not Pedro")
	}
	if user.Email != "pedro@test.com" {
		t.Error("Email is not  pedro@test.com")
	}
	if user.Age != 17 {
		t.Error("Age is not 17")
	}
	if user.Address.Calle != "Av Caseros" {
		t.Error("Address.Calle is not Caseros")
	}
	if user.Address.Direccion != 1001 {
		t.Error("Address.Direccion is not 1001")
	}
	if *user.Address.Dpto != "2 B" {
		t.Error("Address.Dpto is not 2 B", user.Address)
	}

	if *user.Salary != 1234.56 {
		t.Error("Salary is not 1234.56")
	}

	if user.Size != 1.2 {
		t.Error("Size is not 1.2")
	}
}
