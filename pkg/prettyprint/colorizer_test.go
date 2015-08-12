package prettyprint

import (
	"fmt"
	"strings"
	"testing"
)

func TestColorize(t *testing.T) {
	out := Colorize("{{.Red}}Hello {{.Default}}World{{.UnderGreen}}!{{.Default}}")
	expected := "\033[0;31mHello \033[0mWorld\033[4;32m!\033[0m"
	if out != expected {
		t.Errorf("Expected '%s', got '%s'", expected, out)
	}
}

func TestColorizeVars(t *testing.T) {
	vars := map[string]string{"Who": "World"}
	tpl := "{{.C.Red}}Hello {{.C.Default}}{{.V.Who}}{{.C.UnderGreen}}!{{.C.Default}}"
	out := ColorizeVars(tpl, vars)
	expected := "\033[0;31mHello \033[0mWorld\033[4;32m!\033[0m"
	if out != expected {
		t.Errorf("Expected '%s', got '%s'", expected, out)
	}
}

func TestNoColor(t *testing.T) {
	tpl := "{{.Red}}{{.Yellow}}{{.Green}}coffee all the things!{{.Default}}"
	expected := "coffee all the things!"
	out := NoColor(tpl)
	if out != expected {
		t.Errorf("Expected `%s`, got `%s`", expected, out)
	}
}

func ExampleColorize() {
	out := Colorize("{{.Red}}Hello {{.Default}}World{{.UnderGreen}}!{{.Default}}")
	fmt.Println(out)
}

func ExampleColorizeVars() {
	vars := map[string]string{"Who": "World"}
	tpl := "{{.C.Red}}Hello {{.C.Default}}{{.V.Who}}!"
	out := ColorizeVars(tpl, vars)
	fmt.Println(out)
}

func TestDeisIfy(t *testing.T) {
	d := DeisIfy("Test")
	if strings.Contains(d, "Deis1") {
		t.Errorf("Failed to compile template")
	}
	if !strings.Contains(d, "Test") {
		t.Errorf("Failed to render template")
	}
}

func TestLogo(t *testing.T) {
	l := Logo()
	if l != Colors["Deis"] {
		t.Errorf("Expected \n%s\n, Got\n%s\n", Colors["Deis"], Logo())
	}
}

func TestPrettyTabs(t *testing.T) {
	test := map[string]string{
		"test": "testing",
		"foo":  "bar",
	}

	expected := `foo  bar
test testing
`

	output := PrettyTabs(test, 1)

	if output != expected {
		t.Errorf("Expected '%s', Got '%s'", expected, output)
	}
}
