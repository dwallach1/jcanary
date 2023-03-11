package main

import (
	"fmt"
	"log"
	"os"

	"jcanary/engine"
	"jcanary/engine/rules/operators"

	"github.com/Jeffail/gabs"
)

var RULES_CONFIG = getEnvVar("RULES_CONFIG", "./rules.json")

func main() {
	fmt.Println("running jcanary...")

	// parse rules file
	// os.load RULES_CONFIG
	rawConfig, err := gabs.ParseJSONFile(RULES_CONFIG)
	if err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}
	conf, err := engine.New(rawConfig)
	if err != nil {
		log.Fatalf("unable to create engine instance: %v", err)
	}
	println(conf)
	pipeline := []*operators.Result{}
	for r, rule := range conf.Rules {
		fmt.Printf("processing rule #%v\n", r)
		for s, step := range rule.Steps {
			fmt.Printf("processing step #%v\n", s)
			res := step.Operate(conf.Vars, &pipeline)
			fmt.Printf("res -> %v\n", res)
			pipeline = append(pipeline, res)
		}
	}
	fmt.Println("finished running jcanary ")
}

func getEnvVar(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
