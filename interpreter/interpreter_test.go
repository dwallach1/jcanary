package interpreter_test

import (
	"jcanary/interpreter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildString(t *testing.T) {

	cases := []struct {
		Template string
		Bag      interpreter.VariableBag
		Expected string
	}{
		{
			"${baseurl}/checkins",
			interpreter.VariableBag{
				"baseurl": "http://localhost:8080/blend/v1",
			},
			"http://localhost:8080/blend/v1/checkins",
		},
		{
			"nothing/checkins",
			interpreter.VariableBag{
				"baseurl": "http://localhost:8080/blend/v1",
			},
			"nothing/checkins",
		},
		{
			"${baseurl}/${someOtherVal}/checkins",
			interpreter.VariableBag{
				"baseurl":      "http://localhost:8080/blend/v1",
				"someOtherVal": "some other thing",
			},
			"http://localhost:8080/blend/v1/some other thing/checkins",
		},
	}
	for _, c := range cases {
		res := interpreter.BuildString(
			c.Template,
			c.Bag,
		)
		assert.Equal(t, c.Expected, res)
	}
}
