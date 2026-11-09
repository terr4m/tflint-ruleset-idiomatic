package rules

import (
	"cmp"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// keyRule describes where a well-known attribute or nested block belongs.
type keyRule struct {
	// order is the position: <0 must come first, >0 must come last, 0 is "other".
	order int
	// group controls spacing (see the grouping rule): adjacent items in the same
	// group may share lines, adjacent items in different groups must be separated
	// by an empty line. An empty group defaults to the key's own name.
	group string
}

// blockLayout describes the expected layout of a top-level block type.
// Keys that are not listed are "other" attributes: order 0 and a shared group.
type blockLayout struct {
	keys map[string]keyRule
}

// otherGroup is the group shared by every key that is not listed in a layout.
const otherGroup = ""

// rule returns the rule for name and whether the key is listed in the layout.
func (l blockLayout) rule(name string) (keyRule, bool) {
	kr, ok := l.keys[name]
	if !ok {
		return keyRule{group: otherGroup}, false
	}
	if kr.group == "" {
		kr.group = name
	}
	return kr, true
}

// layoutEntry is an attribute or nested block of a body together with its layout rule.
type layoutEntry struct {
	name        string
	order       int
	group       string
	known       bool
	itemRange   hcl.Range
	targetRange hcl.Range
}

// multiLine reports whether the entry spans more than one line.
func (e layoutEntry) multiLine() bool {
	return e.itemRange.End.Line > e.itemRange.Start.Line
}

// layoutEntries returns the attributes and nested blocks of body in source order.
func layoutEntries(body *hclsyntax.Body, layout blockLayout) []layoutEntry {
	entries := make([]layoutEntry, 0, len(body.Attributes)+len(body.Blocks))

	add := func(name string, itemRange, targetRange hcl.Range) {
		kr, known := layout.rule(name)
		entries = append(entries, layoutEntry{
			name:        name,
			order:       kr.order,
			group:       kr.group,
			known:       known,
			itemRange:   itemRange,
			targetRange: targetRange,
		})
	}

	for _, attr := range body.Attributes {
		add(attr.Name, attr.Range(), attr.NameRange)
	}
	for _, b := range body.Blocks {
		add(b.Type, hcl.RangeBetween(b.TypeRange, b.CloseBraceRange), b.TypeRange)
	}

	slices.SortFunc(entries, func(a, b layoutEntry) int {
		if c := cmp.Compare(a.itemRange.Start.Line, b.itemRange.Start.Line); c != 0 {
			return c
		}
		return cmp.Compare(a.itemRange.Start.Column, b.itemRange.Start.Column)
	})

	return entries
}

// orderedGroup returns rules for keys that share one group and must appear in the
// given order, before any other attribute.
func orderedGroup(group string, keys ...string) map[string]keyRule {
	rules := make(map[string]keyRule, len(keys))
	for i, k := range keys {
		rules[k] = keyRule{order: i - len(keys), group: group}
	}
	return rules
}

var resourceKeys = map[string]keyRule{
	"count":       {order: -2},
	"for_each":    {order: -2},
	"provider":    {order: -1},
	"connection":  {order: 1},
	"provisioner": {order: 2},
	"lifecycle":   {order: 3},
	"depends_on":  {order: 4},
}

// blockLayouts is read by both the ordering and sorting rules.
var blockLayouts = map[string]blockLayout{
	"module": {keys: map[string]keyRule{
		"source":                     {order: -5, group: "header"},
		"version":                    {order: -4, group: "header"},
		"ignore_nested_deprecations": {order: -3, group: "header"},
		"count":                      {order: -2},
		"for_each":                   {order: -2},
		"providers":                  {order: -1},
		"depends_on":                 {order: 1},
	}},
	"resource": {keys: resourceKeys},
	"data":     {keys: resourceKeys},
	"variable": {keys: orderedGroup("variable",
		"description", "type", "nullable", "sensitive", "const",
		"ephemeral", "default", "deprecated", "validation",
	)},
	"output": {keys: orderedGroup("output",
		"description", "type", "sensitive", "ephemeral",
		"deprecated", "value", "precondition", "depends_on",
	)},
}
