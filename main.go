package main

import (
	"fmt"
	"log"
	"os"

	"jcanary/engine"
	"jcanary/engine/rules/operators"

	"github.com/Jeffail/gabs"
	"github.com/fatih/color"
)

var RULES_CONFIG = getEnvVar("RULES_CONFIG", "./rules.json")

func main() {
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("running jcanary...")

	successPrinter := color.New(color.FgGreen)
	errPrinter := color.New(color.FgRed)

	// parse rules file
	rawConfig, err := gabs.ParseJSONFile(RULES_CONFIG)
	if err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}
	conf, err := engine.New(rawConfig)
	if err != nil {
		log.Fatalf("unable to create engine instance: %v", err)
	}
	pipeline := []*operators.Result{}

	results := map[string]bool{}

	for r, rule := range conf.Rules {
		c.Printf("\tprocessing rule #%v\n", r)
		for s, step := range rule.Steps {
			c.Printf("\t\tprocessing step #%v\n", s)
			res := step.Operate(conf.Vars, &pipeline)
			if res.HasError() {
				errPrinter.Printf("\t\t\tstep #%v failed: %v\n", s, res.Err)
			}
			pipeline = append(pipeline, res)
		}

		ruleHasFailure := false
		for _, res := range pipeline {
			if res.HasError() {
				ruleHasFailure = true
			}
		}
		results[rule.Name] = ruleHasFailure

		// reset pipeline for new rule
		pipeline = []*operators.Result{}
	}

	fmt.Printf("\n\n----------------------------\n")
	failureDetected := false
	for ruleName, hasFailure := range results {
		verb := successPrinter.Sprintf("passed")
		if hasFailure {
			verb = errPrinter.Sprintf("failed")
			failureDetected = true
		}
		fmt.Printf("rule \"%v\": %v\n", ruleName, verb)
	}
	fmt.Printf("----------------------------\n")
	fmt.Println("finished running jcanary")
	if failureDetected {
		os.Exit(1)
	}
}

func getEnvVar(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
