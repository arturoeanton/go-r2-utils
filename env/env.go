package env

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Getenv struct {
	key          string
	defaultValue interface{}
	value        interface{}
}

func LoadEnv(filename string) error {

	if len(filename) == 0 {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				os.Setenv(key, value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func EnvKey(key string) *Getenv {
	g := &Getenv{key: key}
	return g
}

func (g *Getenv) Key() string { return g.key }

func (g *Getenv) Default(d interface{}) *Getenv {
	g.defaultValue = d
	return g
}

func (g *Getenv) Interface() (interface{}, bool) {
	value := os.Getenv(g.key)
	g.value = g.defaultValue
	if value != "" {
		g.value = value
		return g.value, true
	}
	return g.value, false
}

func (g *Getenv) String() string {
	if g.defaultValue == nil {
		g.defaultValue = ""
	}
	v, _ := g.Interface()
	return v.(string)
}

func (g *Getenv) Int() int {
	if g.defaultValue == nil {
		g.defaultValue = 0
	}
	v, d := g.Interface()
	if !d {
		return v.(int)
	}
	value, error := strconv.ParseInt(v.(string), 10, 64)
	if error != nil {
		return g.defaultValue.(int)
	}
	return int(value)
}

func (g *Getenv) Float64() float64 {
	if g.defaultValue == nil {
		g.defaultValue = 0
	}
	v, d := g.Interface()
	if !d {
		return v.(float64)
	}
	value, error := strconv.ParseFloat(v.(string), 64)
	if error != nil {
		return g.defaultValue.(float64)
	}
	return value
}

func (g *Getenv) Float() float32 {
	if g.defaultValue == nil {
		g.defaultValue = 0
	}
	v, d := g.Interface()
	if !d {
		return float32(v.(float64))
	}
	value, error := strconv.ParseFloat(v.(string), 32)
	if error != nil {
		return g.defaultValue.(float32)
	}
	return float32(value)
}

func (g *Getenv) Boolean() bool {
	if g.defaultValue == nil {
		g.defaultValue = false
	}
	v, d := g.Interface()
	if !d {
		return v.(bool)
	}
	value, error := strconv.ParseBool(v.(string))
	if error != nil {
		return g.defaultValue.(bool)
	}
	return value
}
