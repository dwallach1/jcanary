package engine

import (
	"verify.me/engine/rules"
	"verify.me/interoperator"
)

type Config struct {
	Rules []rules.Rule `json:"rules"`
	Vars  map[string]interoperator.Variable
}
