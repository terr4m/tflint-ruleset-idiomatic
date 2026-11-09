package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type sortable interface {
	sortName() string
	itemRange() hcl.Range
	targetRange() hcl.Range
	kind() string
}

type attrItem struct{ *hclsyntax.Attribute }

func (a attrItem) sortName() string       { return a.Name }
func (a attrItem) itemRange() hcl.Range   { return a.Range() }
func (a attrItem) targetRange() hcl.Range { return a.NameRange }
func (a attrItem) kind() string           { return "Attribute" }

type objItem struct {
	key  string
	item hclsyntax.ObjectConsItem
}

func (o objItem) sortName() string       { return o.key }
func (o objItem) targetRange() hcl.Range { return o.item.KeyExpr.Range() }
func (o objItem) kind() string           { return "Key" }
func (o objItem) itemRange() hcl.Range {
	return hcl.RangeBetween(o.item.KeyExpr.Range(), o.item.ValueExpr.Range())
}

// sortingRule enforces alphabetical sorting of attributes and complex values (objects/maps),
// using empty lines to delimit independent logical groups.
type sortingRule struct {
	tflint.DefaultRule
}

func NewSortingRule() tflint.Rule {
	return &sortingRule{}
}

func (r *sortingRule) Name() string {
	return "idiomatic_sorting"
}

func (r *sortingRule) Enabled() bool {
	return true
}

func (r *sortingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *sortingRule) Link() string {
	return "https://terr4m.github.io/tflint-ruleset-idiomatic/rules/sorting/"
}

func (r *sortingRule) Check(runner tflint.Runner) error {
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

func (r *sortingRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block, layout blockLayout) error {
	if block == nil || block.Body == nil {
		return nil
	}

	// 1. Check block-level attributes
	if len(block.Body.Attributes) >= 2 {
		// The ordering rule dictates the relative order of keys with different positions.
		exempt := func(a, b sortable) bool {
			ka, _ := layout.rule(a.sortName())
			kb, _ := layout.rule(b.sortName())
			return ka.order != 0 && kb.order != 0 && ka.order != kb.order
		}

		items := make([]sortable, 0, len(block.Body.Attributes))
		for _, attr := range block.Body.Attributes {
			items = append(items, attrItem{attr})
		}
		if err := r.checkOrder(runner, items, exempt); err != nil {
			return err
		}
	}

	// 2. Check expressions inside attributes for complex values
	for _, attr := range block.Body.Attributes {
		if attr.Expr != nil {
			if err := r.checkExpressions(runner, attr.Expr); err != nil {
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

// nodeWalker implements [hclsyntax.Walker].
type nodeWalker struct {
	enter func(hclsyntax.Node) hcl.Diagnostics
}

func (w *nodeWalker) Enter(node hclsyntax.Node) hcl.Diagnostics {
	if w.enter != nil {
		return w.enter(node)
	}
	return nil
}

func (w *nodeWalker) Exit(_ hclsyntax.Node) hcl.Diagnostics {
	return nil
}

// objectExprs returns every object expression within expr (including expr itself)
// that has at least two items.
func objectExprs(expr hclsyntax.Expression) []*hclsyntax.ObjectConsExpr {
	// The walker can't return a Go error, so collect the objects and let the caller check them.
	var objects []*hclsyntax.ObjectConsExpr

	walker := &nodeWalker{
		enter: func(node hclsyntax.Node) hcl.Diagnostics {
			if obj, ok := node.(*hclsyntax.ObjectConsExpr); ok && len(obj.Items) >= 2 {
				objects = append(objects, obj)
			}
			return nil
		},
	}

	// The walker never returns diagnostics, so the result carries no information.
	_ = hclsyntax.Walk(expr, walker)

	return objects
}

func (r *sortingRule) checkExpressions(runner tflint.Runner, expr hclsyntax.Expression) error {
	for _, obj := range objectExprs(expr) {
		items := make([]sortable, 0, len(obj.Items))
		for _, item := range obj.Items {
			if key := keyString(item.KeyExpr); key != "" {
				items = append(items, objItem{
					key:  key,
					item: item,
				})
			}
		}

		if len(items) >= 2 {
			if err := r.checkOrder(runner, items, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkOrder reports items that are out of alphabetical order within a logical group.
// Pairs for which exempt returns true are not compared.
func (r *sortingRule) checkOrder(runner tflint.Runner, items []sortable, exempt func(a, b sortable) bool) error {
	sort.Slice(items, func(i, j int) bool {
		r1, r2 := items[i].itemRange(), items[j].itemRange()
		if r1.Start.Line != r2.Start.Line {
			return r1.Start.Line < r2.Start.Line
		}
		return r1.Start.Column < r2.Start.Column
	})

	for i := 0; i < len(items)-1; i++ {
		curr, next := items[i], items[i+1]

		if next.itemRange().Start.Line > curr.itemRange().End.Line+1 {
			continue
		}
		if exempt != nil && exempt(curr, next) {
			continue
		}

		if strings.Compare(curr.sortName(), next.sortName()) > 0 {
			msg := fmt.Sprintf("%s `%s` should appear after `%s` in this logical group (attributes must be sorted alphabetically).", curr.kind(), curr.sortName(), next.sortName())
			if err := runner.EmitIssue(r, msg, curr.targetRange()); err != nil {
				return err
			}
		}
	}

	return nil
}

// keyString returns the key as Terraform sees it: a statically known value
// converted to a string (so `80` becomes "80" and `true` becomes "true").
// It returns an empty string for keys that can't be evaluated statically.
func keyString(expr hclsyntax.Expression) string {
	val, diags := expr.Value(nil)
	if diags.HasErrors() || !val.IsWhollyKnown() || val.IsNull() {
		return ""
	}

	val, err := convert.Convert(val, cty.String)
	if err != nil || !val.IsKnown() || val.IsNull() {
		return ""
	}

	return val.AsString()
}
