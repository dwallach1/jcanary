package rules

import "jcanary/engine/rules/operators"

type Rule struct {
	Name  string               `json:"name"`
	Steps []operators.Operator `json:"steps"`
}

// func NewRule(rawConfig map[string]interface{})
