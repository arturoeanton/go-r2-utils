package commons

import (
	"io/ioutil"
	"os"
)

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func FileToString(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	return string(content), err
}

func FileToString_(name string) string {
	content, _ := ioutil.ReadFile(name)
	return string(content)
}

func StringToFile(filenanme string, content string) error {
	d1 := []byte(content)
	err := ioutil.WriteFile(filenanme, d1, 0644)
	return err
}
func RemoveRepeat(s []string) []string {
	m := make(map[string]bool)
	for _, v := range s {
		m[v] = true
	}
	set := make([]string, 0)
	for k := range m {
		set = append(set, k)
	}
	return set
}
