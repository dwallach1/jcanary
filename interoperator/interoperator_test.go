package interoperator_test

import (
	"testing"

	"jcanary/interoperator"

	"github.com/stretchr/testify/assert"
)

func TestBuildString(t *testing.T) {

	cases := []struct {
		Template string
		Bag      interoperator.VariableBag
		Expected string
	}{
		{
			"${baseurl}/checkins",
			interoperator.VariableBag{
				"baseurl": "http://localhost:8080/blend/v1",
			},
			"http://localhost:8080/blend/v1/checkins",
		},
		{
			"nothing/checkins",
			interoperator.VariableBag{
				"baseurl": "http://localhost:8080/blend/v1",
			},
			"nothing/checkins",
		},
		{
			"${baseurl}/${someOtherVal}/checkins",
			interoperator.VariableBag{
				"baseurl":      "http://localhost:8080/blend/v1",
				"someOtherVal": "some other thing",
			},
			"http://localhost:8080/blend/v1/some other thing/checkins",
		},
	}
	for _, c := range cases {
		res := interoperator.BuildString(
			c.Template,
			c.Bag,
		)
		assert.Equal(t, c.Expected, res)
	}
}
