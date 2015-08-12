// Package prettyprint contains tools for formatting text.
package prettyprint

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

// Colors contains a map of the standard ANSI color codes.
//
// There are four variants:
// 	- Bare color names (Red, Black) color the characters.
// 	- Bold color names add bolding to the characters.
// 	- Under color names add underlining to the characters.
// 	- Hi color names add highlighting (background colors).
//
// These can be used within `text/template` to provide colors. The convenience
// function `Colorize()` provides this feature.
var Colors = map[string]string{
	"Default":     "\033[0m",
	"Black":       "\033[0;30m",
	"Red":         "\033[0;31m",
	"Green":       "\033[0;32m",
	"Yellow":      "\033[0;33m",
	"Blue":        "\033[0;34m",
	"Purple":      "\033[0;35m",
	"Cyan":        "\033[0;36m",
	"White":       "\033[0;37m",
	"BoldBlack":   "\033[1;30m",
	"BoldRed":     "\033[1;31m",
	"BoldGreen":   "\033[1;32m",
	"BoldYellow":  "\033[1;33m",
	"BoldBlue":    "\033[1;34m",
	"BoldPurple":  "\033[1;35m",
	"BoldCyan":    "\033[1;36m",
	"BoldWhite":   "\033[1;37m",
	"UnderBlack":  "\033[4;30m",
	"UnderRed":    "\033[4;31m",
	"UnderGreen":  "\033[4;32m",
	"UnderYellow": "\033[4;33m",
	"UnderBlue":   "\033[4;34m",
	"UnderPurple": "\033[4;35m",
	"UnderCyan":   "\033[4;36m",
	"UnderWhite":  "\033[4;37m",
	"HiBlack":     "\033[30m",
	"HiRed":       "\033[31m",
	"HiGreen":     "\033[32m",
	"HiYellow":    "\033[33m",
	"HiBlue":      "\033[34m",
	"HiPurple":    "\033[35m",
	"HiCyan":      "\033[36m",
	"HiWhite":     "\033[37m",
	"Deis1":       "\033[31m● \033[34m▴ \033[32m■\033[0m",
	"Deis2":       "\033[32m■ \033[31m● \033[34m▴\033[0m",
	"Deis3":       "\033[34m▴ \033[32m■ \033[31m●\033[0m",
	"Deis":        "\033[31m● \033[34m▴ \033[32m■\n\033[32m■ \033[31m● \033[34m▴\n\033[34m▴ \033[32m■ \033[31m●\n",
}

// DeisIfy returns a pretty-printed deis logo along with the corresponding message
func DeisIfy(msg string) string {
	var t = struct {
		Msg string
		C   map[string]string
	}{
		Msg: msg,
		C:   Colors,
	}
	tpl := "{{.C.Deis1}}\n{{.C.Deis2}} {{.Msg}}\n{{.C.Deis3}}\n"
	var buf bytes.Buffer
	template.Must(template.New("deis").Parse(tpl)).Execute(&buf, t)
	return buf.String()
}

// Logo returns a colorized Deis logo with no space for text.
func Logo() string {
	return Colorize("{{.Deis}}")
}

// NoColor strips colors from the template.
//
// NoColor provides support for non-color ANSI terminals. It can be used
// as an alternative to Colorize when it is detected that the terminal does
// not support colors.
func NoColor(msg string) string {
	empties := make(map[string]string, len(Colors))
	for k := range Colors {
		empties[k] = ""
	}
	return colorize(msg, empties)
}

// Colorize makes it easy to add colors to ANSI terminal output.
//
// This takes any of the colors defined in the Colors map. Colors are rendered
// through the `text/template` system, so you may use pipes and functions as
// well.
//
// Example:
//	Colorize("{{.Red}}ERROR:{{.Default}} Something happened.")
func Colorize(msg string) string {
	return colorize(msg, Colors)
}

// ColorizeVars provides template rendering with color support.
//
// The template is given a datum with two objects: `.V` and `.C`. `.V` contains
// the `vars` passed into the function. `.C` contains the color map.
//
// Assuming `vars` contains a member named `Msg`, a template can be constructed
// like this:
//	{{.C.Red}}Message:{{.C.Default}} .V.Msg
func ColorizeVars(msg string, vars interface{}) string {
	var t = struct {
		V interface{}
		C map[string]string
	}{
		V: vars,
		C: Colors,
	}
	return colorize(msg, t)
}

func colorize(msg string, vars interface{}) string {
	tpl, err := template.New(msg).Parse(msg)
	// If the template's not valid, we just ignore and return.
	if err != nil {
		return msg
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, vars); err != nil {
		return msg
	}

	return buf.String()
}

// Overwrite sends a line that will be replaced by a subsequent overwrite.
//
// Example:
// 	Overwrite("foo")
// 	Overwrite("bar")
//
// The above will print "foo" and then immediately replace it with "var".
//
// (Interpretation of \r is left to the shell.)
func Overwrite(msg string) string {
	lm := len(msg)
	if lm >= 80 {
		return msg + "\r"
	}
	pad := 80 - len(msg)
	return msg + strings.Repeat(" ", pad) + "\r"

}

// Overwritef formats a string and then returns an overwrite line.
//
// See `Overwrite` for details.
func Overwritef(msg string, args ...interface{}) string {
	return Overwrite(fmt.Sprintf(msg, args...))
}

// PrettyTabs formats a map with with alligned keys and values.
//
// Example:
// test := map[string]string {
//    "test": "testing",
//    "foo": "bar",
//  }
//
// Prettytabs(test, 5)
//
// This will return a formatted string.
// The previous example would return:
// foo      bar
// test     testing
func PrettyTabs(msg map[string]string, spaces int) string {
	// find the longest key so we know how much padding to use
	max := 0
	for key := range msg {
		if len(key) > max {
			max = len(key)
		}
	}
	max += spaces

	// sort the map keys so we can print them alphabetically
	var keys []string
	for k := range msg {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var output string
	for _, k := range keys {
		output += fmt.Sprintf("%s%s%s\n", k, strings.Repeat(" ", max-len(k)), msg[k])
	}
	return output
}
