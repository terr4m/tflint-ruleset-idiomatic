package rules

import (
	"fmt"
	"slices"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// orderingRule enforces a consistent order of well-known attributes and blocks in
// top-level blocks. Spacing between them is handled by the grouping rule.
type orderingRule struct {
	tflint.DefaultRule
}

func NewOrderingRule() tflint.Rule {
	return &orderingRule{}
}

func (r *orderingRule) Name() string {
	return "idiomatic_ordering"
}

func (r *orderingRule) Enabled() bool {
	return true
}

func (r *orderingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *orderingRule) Link() string {
	return "https://terr4m.github.io/tflint-ruleset-idiomatic/rules/ordering/"
}

func (r *orderingRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}
		for _, block := range body.Blocks {
			if err := r.checkBlock(runner, block); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *orderingRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	if block == nil || block.Body == nil {
		return nil
	}

	layout, ok := blockLayouts[block.Type]
	if !ok {
		return nil
	}

	return r.checkOrder(runner, layoutEntries(block.Body, layout))
}

// checkOrder reports well-known items that are on the wrong side of other items.
func (r *orderingRule) checkOrder(runner tflint.Runner, entries []layoutEntry) error {
	for i, e := range entries {
		var msg string

		switch {
		case e.order < 0:
			if idx := slices.IndexFunc(entries[:i], func(o layoutEntry) bool { return o.order > e.order }); idx >= 0 {
				msg = fmt.Sprintf("`%s` should appear before `%s`.", e.name, entries[idx].name)
			}
		case e.order > 0:
			if idx := slices.IndexFunc(entries[i+1:], func(o layoutEntry) bool { return o.order < e.order }); idx >= 0 {
				msg = fmt.Sprintf("`%s` should appear after `%s`.", e.name, entries[i+1+idx].name)
			}
		}

		if msg != "" {
			if err := runner.EmitIssue(r, msg, e.targetRange); err != nil {
				return err
			}
		}
	}

	return nil
}
