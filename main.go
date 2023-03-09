package main

import (
	"fmt"
	"os"
)

var RULES_CONFIG = getEnvVar("RULES_CONFIG", "./rules.json")

func main() {
	fmt.Println("running verify.me...")

	// parse rules file
	// os.load RULES_CONFIG
}

func getEnvVar(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
