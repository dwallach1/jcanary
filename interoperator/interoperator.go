package interoperator

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)\${[a-zA-Z0-9]*}`)

type VariableBag map[string]string

type VariableType string

const (
	Constant VariableType = "constant"
)

type Variable struct {
	Type VariableType
}

type MissingVariableErr struct {
	Var string
}

func (v *MissingVariableErr) Error() string {
	return "reference to undefined variable " + v.Var
}

// BuildString
//
//	template -> ${basurl}/checkins , bag -> { "baseurl": "localhost:8080" }
//	should return "localhost:8080/checkins"
func BuildString(template string, bag VariableBag) string {
	builder := template
	for _, match := range re.FindAllString(template, -1) {
		key := match[2 : len(match)-1]
		val, ok := bag[key]
		if !ok {
			return ""
		}
		builder = strings.Replace(builder, match, val, 1)
	}
	return builder
}
