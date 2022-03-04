package snippet

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"
)

type Var struct {
	Key   string
	Value string
	Type  string
}

type Snippet struct {
	Name     string
	Type     string
	Text     string
	Vars     map[string]*Var
	FxEscape func(string) string
}
type Snippets struct {
	snippet  map[string]*Snippet
	comment  string
	fxEscape func(string) string
}

func NewSnippetStorage() *Snippets {
	return &Snippets{
		snippet:  make(map[string]*Snippet),
		comment:  "--",
		fxEscape: template.HTMLEscapeString,
	}
}
func (s *Snippets) Comment(comment string) *Snippets {
	s.comment = comment
	return s
}

func (s *Snippets) Escape(fx func(string) string) *Snippets {
	s.fxEscape = fx
	return s
}

func (s *Snippets) LoadFile(name string) *Snippets {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	return s.Load(scanner)
}

func (s *Snippets) LoadString(code string) *Snippets {
	scanner := bufio.NewScanner(strings.NewReader(code))
	return s.Load(scanner)
}

func (s *Snippets) parce(lineTmp, field string) (string, bool) {
	selectorPrefix := s.comment + field + ":"
	if len(lineTmp) > 0 && strings.HasPrefix(lineTmp, selectorPrefix) {
		key := strings.TrimPrefix(lineTmp, selectorPrefix)
		return key, true
	}
	return "", false
}

func (s *Snippets) Load(scanner *bufio.Scanner) *Snippets {
	key := ""
	for scanner.Scan() {
		line := scanner.Text()
		lineTmp := strings.ReplaceAll(line, " ", "")

		val, flag := s.parce(lineTmp, "name")
		if flag {
			key = val
			s.snippet[key] = &Snippet{Name: key, Type: "", Text: "", Vars: make(map[string]*Var), FxEscape: s.fxEscape}
			continue
		}

		if key == "" {
			continue
		}

		val, flag = s.parce(lineTmp, "var")
		if flag {
			split := strings.Split(val, "=")
			varKey := strings.ReplaceAll(split[0], " ", "")
			i := strings.Index(line, "=")
			value := ""
			if i > 0 {
				value = line[i+1:]
			}
			nameSplit := strings.Split(varKey, ":")
			varName := nameSplit[0]
			varType := "string"
			if len(nameSplit) > 1 {
				varType = nameSplit[1]
			}
			if varType == "float" {
				varType = "float64"
			}
			s.snippet[key].Vars[varName] = &Var{Key: varName, Value: value, Type: varType}
			continue
		}

		s.snippet[key].Text += line + "\n"
	}
	return s
}

func (s *Snippets) GetSnippet(name string) Snippet {
	return *s.snippet[name]
}

func (q *Snippet) Escape(f func(string) string) *Snippet {
	q.FxEscape = f
	return q
}

func (q *Snippet) Param(key string, value interface{}) *Snippet {
	str := fmt.Sprint(value)
	typeParamName := reflect.TypeOf(value).Name()
	if q.Vars[key].Type != typeParamName &&
		(q.Vars[key].Type == "number" && !(strings.HasPrefix(typeParamName, "int") ||
			strings.HasPrefix(typeParamName, "uint") ||
			strings.HasPrefix(typeParamName, "float"))) {
		panic(fmt.Sprintf("Type of param %s is %s and it must  be not %s", key, typeParamName, q.Vars[key].Type))
	}
	q.Vars[key].Value = str
	return q

}

func (q *Snippet) Get() string {
	t := q.Text
	flagEscape := q.FxEscape != nil
	for k, v := range q.Vars {
		literal := "${{" + k + "}}"
		if flagEscape {
			t = strings.ReplaceAll(t, literal, q.FxEscape(v.Value))
			continue
		}
		t = strings.ReplaceAll(t, literal, v.Value)
	}
	return t
}
