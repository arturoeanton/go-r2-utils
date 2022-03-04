package snippet_test

import (
	"strings"
	"testing"
	"text/template"

	"github.com/arturoeanton/go-r2-utils/snippet"
)

func TestSnippet(t *testing.T) {

	code := `
--name:test1
Select * from test

--name: test2 
--var:id
--var:name
--var:age:number=14
--var:text:string=Hola como estas 1 = 1
Select *
from test
where id = "${{id}}"
and name = "${{name}}"
and age = ${{age}}
and text = "${{text}}"
`

	qb := snippet.NewSnippetStorage().Escape(template.HTMLEscapeString).Comment("--").LoadString(code)

	qt := qb.GetSnippet("test2")

	result := `Select * from test where id = "hola" and name = "Elias&#34;--" and age = 4.1 and text = "Hola como estas 1 = 1"`
	out := qt.
		Param("id", "hola").
		Param("name", "Elias\"--").
		Param("age", 4.1).
		Get()
	out1 := strings.Trim(out, " \n\t\r")
	result1 := strings.Trim(result, " \n\t\r")
	out1 = strings.Replace(out1, "\n", " ", -1)
	out1 = strings.Replace(out1, "\t", " ", -1)

	result1 = strings.Replace(result1, "\n", " ", -1)
	result1 = strings.Replace(result1, "\t", " ", -1)

	if out1 != result1 {
		t.Errorf("Error\n|%s|\n|%s|\n", result1, out1)
	}

}
