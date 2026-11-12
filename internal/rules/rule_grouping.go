package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// groupingRule enforces consistent spacing between attributes, nested blocks and
// complex values (objects/maps). Items that span multiple lines are kept apart from
// their neighbours, and well-known items of top-level blocks are kept apart from
// items that are not part of the same group.
type groupingRule struct {
	tflint.DefaultRule
}

func NewGroupingRule() tflint.Rule {
	return &groupingRule{}
}

func (r *groupingRule) Name() string {
	return "idiomatic_grouping"
}

func (r *groupingRule) Enabled() bool {
	return true
}

func (r *groupingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *groupingRule) Link() string {
	return "https://terr4m.github.io/tflint-ruleset-idiomatic/rules/grouping/"
}

func (r *groupingRule) Check(runner tflint.Runner) error {
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
			// Layouts only apply to top-level blocks.
			if err := r.checkBlock(runner, block, blockLayouts[block.Type]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *groupingRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block, layout blockLayout) error {
	if block == nil || block.Body == nil {
		return nil
	}

	// 1. Check the attributes and nested blocks of the block
	if err := r.checkEntries(runner, layoutEntries(block.Body, layout)); err != nil {
		return err
	}

	// 2. Check expressions inside attributes for complex values
	for _, attr := range block.Body.Attributes {
		if attr.Expr == nil {
			continue
		}
		for _, obj := range objectExprs(attr.Expr) {
			if err := r.checkEntries(runner, objectEntries(obj)); err != nil {
				return err
			}
		}
	}

	// 3. Recurse into child blocks
	for _, child := range block.Body.Blocks {
		if err := r.checkBlock(runner, child, blockLayout{}); err != nil {
			return err
		}
	}

	return nil
}

// checkEntries reports adjacent entries that need an empty line between them but
// don't have one. An entry needs separating from its neighbours when it spans
// multiple lines, or when the two entries belong to different groups.
func (r *groupingRule) checkEntries(runner tflint.Runner, entries []layoutEntry) error {
	for i := 0; i < len(entries)-1; i++ {
		curr, next := entries[i], entries[i+1]

		// Entries that are on the same line (only possible in objects) can't be
		// separated, and an empty line or a comment line already separates the rest.
		if next.itemRange.Start.Line != curr.itemRange.End.Line+1 {
			continue
		}

		var target, other layoutEntry
		var msg string

		switch {
		case curr.multiLine() || next.multiLine():
			// Report the multi-line entry; if both are, report the later one.
			target, other = next, curr
			if !next.multiLine() {
				target, other = curr, next
			}
			msg = fmt.Sprintf("`%s` spans multiple lines and should be separated from `%s` by an empty line.", target.name, other.name)
		case curr.group != next.group:
			// Report the well-known entry; if both are, report the later one.
			target, other = next, curr
			if !next.known {
				target, other = curr, next
			}
			msg = fmt.Sprintf("`%s` should be separated from `%s` by an empty line (they are not part of the same group).", target.name, other.name)
		default:
			continue
		}

		if err := runner.EmitIssue(r, msg, target.targetRange); err != nil {
			return err
		}
	}

	return nil
}

// objectEntries returns the items of obj in source order. Items have no layout.
func objectEntries(obj *hclsyntax.ObjectConsExpr) []layoutEntry {
	entries := make([]layoutEntry, 0, len(obj.Items))

	for _, item := range obj.Items {
		name := keyString(item.KeyExpr)
		if name == "" {
			name = "(dynamic key)"
		}

		entries = append(entries, layoutEntry{
			name:        name,
			group:       otherGroup,
			itemRange:   hcl.RangeBetween(item.KeyExpr.Range(), item.ValueExpr.Range()),
			targetRange: item.KeyExpr.Range(),
		})
	}

	return entries
}
