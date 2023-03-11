package engine

import (
	"errors"
	"fmt"

	"jcanary/engine/rules"
	"jcanary/engine/rules/operators"
	"jcanary/interpreter"

	"github.com/Jeffail/gabs"
)

type Config struct {
	Rules []rules.Rule `json:"rules"`
	Vars  interpreter.VariableBag
}

func New(rawconfig *gabs.Container) (*Config, error) {
	configMap, err := rawconfig.ChildrenMap()
	if err != nil {
		return nil, fmt.Errorf("unable to load rules array from config: %w", err)
	}
	// @TODO: this can cause nil pointer exceptions.. fix later
	ruleContainer := configMap["rules"].Data()
	var myrules []rules.Rule
	fmt.Println(ruleContainer)
	rawRules, ok := ruleContainer.([]interface{})
	if !ok {
		return nil, errors.New("invalid rules object in config")
	}
	for i, rawRuleContainer := range rawRules {
		rawRule := rawRuleContainer.(map[string]interface{})
		rulename := rawRule["name"].(string)
		rawsteps := rawRule["steps"].([]interface{})
		var mysteps []operators.Operator
		if !ok {
			return nil, fmt.Errorf("malformed steps array in config at idx %v", i)
		}
		for j, rawstep := range rawsteps {
			step := rawstep.(map[string]interface{})
			operatorname := step["action"].(string)
			op, err := operators.New(
				operators.OperatorType(operatorname),
				step,
			)
			if err != nil {
				return nil, fmt.Errorf("invalid step instance at idx %v: %w", j, err)
			}
			mysteps = append(mysteps, op)
		}
		myrules = append(myrules, rules.Rule{
			Name:  rulename,
			Steps: mysteps,
		})

	}

	varsmap := configMap["vars"].Data().(map[string]interface{})
	return &Config{
		Rules: myrules,
		Vars:  parseVars(varsmap),
	}, nil
}

func parseVars(vars map[string]interface{}) interpreter.VariableBag {
	bag := interpreter.VariableBag{}
	for k, v := range vars {

		varconfig := v.(map[string]interface{})
		varType := varconfig["type"]
		if varType == "constant" {
			bag[k] = varconfig["value"].(string)
		}
	}
	return bag
}
