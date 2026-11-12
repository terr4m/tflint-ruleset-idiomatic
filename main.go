package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/terr4m/tflint-ruleset-idiomatic/internal/rules"
)

var version = "dev"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "idiomatic",
			Version: version,
			Rules:   rules.All(),
		},
	})
}
