package operators

import (
	"fmt"
	"jcanary/interpreter"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type SchematizeOperator struct {
	Type       OperatorType           `json:"type"`
	StepRef    int                    `json:"stepRef"`
	Path       string                 `json:"path"`
	RootSchema map[string]interface{} `json:"rootSchema"`
}

func (o *SchematizeOperator) Operate(varBag interpreter.VariableBag, pipeline *[]*Result) *Result {
	var result Result
	resultStringToTest := (*pipeline)[o.StepRef].Container.Path(o.Path).String()

	sl := gojsonschema.NewSchemaLoader()
	sl.Draft = gojsonschema.Draft7
	sl.AutoDetect = false

	rootLoader := gojsonschema.NewGoLoader(o.RootSchema)
	schema, err := sl.Compile(rootLoader)
	if err != nil {
		Print("unable to compile schema: %v", err)
		result.Err = err
		return &result
	}

	document := gojsonschema.NewStringLoader(resultStringToTest)
	res, err := schema.Validate(document)
	if err != nil {
		result.Err = fmt.Errorf("failed to perform schema validation: %w", err)
		return &result
	}
	if res.Valid() {
		Print("The document is valid according to input schema\n")
	} else {
		errStrs := []string{}
		for _, desc := range res.Errors() {
			errStrs = append(errStrs, fmt.Sprintf("- %s", desc))
		}
		e := fmt.Errorf("object does not conform to schema: %v", strings.Join(errStrs, ", "))
		Print(e.Error())
		result.Err = e
	}
	return &result
}
