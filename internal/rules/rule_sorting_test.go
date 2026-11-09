package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_sortingRule(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		content string
		want    helper.Issues
	}{
		{
			name: "valid_root_block_attributes",
			content: `
locals {
  beta = "..."
  gamma  = "..."

  alpha = "..."
  zeta = "..."
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_resource_block_attributes",
			content: `
resource "dummy_resource" "example" {
  beta = "..."
  gamma  = "..."

  alpha = "..."
  zeta = "..."
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_resource_block_nested_block_attributes",
			content: `
resource "dummy_resource" "example" {
  child {
    beta = "..."
    gamma  = "..."

    alpha = "..."
    zeta = "..."
  }
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_complex_single_line",
			content: `
locals {
  test = {alpha = "...", beta = "...", gamma = "...", zeta = "..."}
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_complex_multi_line",
			content: `
locals {
  test = {
    beta = "..."
    gamma  = "..."

    alpha = "..."
    zeta = "..."
  }
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_complex_mixed_line",
			content: `
locals {
  test = {beta = "...", gamma = "...",

    alpha = "..."
    zeta = "..."
  }
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_complex_numeric_keys_sorted_as_strings",
			content: `
locals {
  test = {443 = "...", 80 = "..."}
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_complex_bool_keys",
			content: `
locals {
  test = {false = "...", true = "..."}
}
`,
			want: helper.Issues{},
		},
		{
			name: "invalid_root_block_attributes",
			content: `
locals {
  beta  = "2"
  alpha = "1"
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Attribute `beta` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 7},
					},
				},
			},
		},
		{
			name: "invalid_resource_block_attributes",
			content: `
resource "dummy_resource" "example" {
  beta  = "2"
  alpha = "1"
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Attribute `beta` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 7},
					},
				},
			},
		},
		{
			name: "invalid_resource_block_nested_block_attributes",
			content: `
resource "dummy_resource" "example" {
  child {
    beta  = "2"
    alpha = "1"
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Attribute `beta` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 4, Column: 9},
					},
				},
			},
		},
		{
			name: "invalid_complex_single_line",
			content: `
locals {
  test = {beta = "...", alpha = "..."}
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Key `beta` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 15},
					},
				},
			},
		},
		{
			name: "invalid_complex_numeric_keys",
			content: `
locals {
  test = {80 = "...", 443 = "..."}
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Key `80` should appear after `443` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 13},
					},
				},
			},
		},
		{
			name: "invalid_complex_multi_line",
			content: `
locals {
  test = {
    beta = "..."
    alpha = "..."
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Key `beta` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 4, Column: 9},
					},
				},
			},
		},
		{
			name: "invalid_complex_mixed_line",
			content: `
locals {
  test = {gamma = "...", beta = "...",

    zeta = "..."
    alpha = "..."
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Key `gamma` should appear after `beta` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 16},
					},
				},
				{
					Rule:    NewSortingRule(),
					Message: "Key `zeta` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 5, Column: 5},
						End:      hcl.Pos{Line: 5, Column: 9},
					},
				},
			},
		},
		{
			name: "skip_variable_ranked_keys_are_left_to_ordering",
			content: `
variable "test" {
  description = "..."
  type        = string
  nullable    = false
  default     = "a"
}
`,
			want: helper.Issues{},
		},
		{
			name: "skip_module_header_keys_are_left_to_ordering",
			content: `
module "test" {
  source                     = "./test"
  version                    = "1.0.0"
  ignore_nested_deprecations = true
}
`,
			want: helper.Issues{},
		},
		{
			name: "skip_resource_ranked_keys_out_of_order_are_left_to_ordering",
			content: `
resource "dummy_resource" "example" {
  provider = dummy.alt
  count    = 1
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_resource_ranked_and_other_keys_separated",
			content: `
resource "dummy_resource" "example" {
  zeta = "..."

  depends_on = [a]
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_comment_line_starts_new_group",
			content: `
locals {
  beta = "..."
  # alpha is documented here
  alpha = "..."
}
`,
			want: helper.Issues{},
		},
		{
			name: "invalid_resource_ranked_key_next_to_other_key",
			content: `
resource "dummy_resource" "example" {
  count = 1
  alpha = "..."
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Attribute `count` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 8},
					},
				},
			},
		},
		{
			name: "invalid_variable_complex",
			content: `
variable "test" {
  type = object({
    zeta  = string
    alpha = string
  })
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Key `zeta` should appear after `alpha` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 4, Column: 9},
					},
				},
			},
		},
		{
			name: "invalid_variable_child_block",
			content: `
variable "test" {
  validation {
    error_message = "..."
    condition     = "..."
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Attribute `error_message` should appear after `condition` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 4, Column: 18},
					},
				},
			},
		},
		{
			name: "invalid_unconfigured_block_well_known_names_are_sorted",
			content: `
locals {
  depends_on = "..."
  count      = 1
}
`,
			want: helper.Issues{
				{
					Rule:    NewSortingRule(),
					Message: "Attribute `depends_on` should appear after `count` in this logical group (attributes must be sorted alphabetically).",
					Range: hcl.Range{
						Filename: "test.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 13},
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rule := NewSortingRule()

			runner := helper.TestRunner(t, map[string]string{
				"test.tf": tt.content,
			})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected check error: %s", err)
			}

			helper.AssertIssues(t, tt.want, runner.Issues)
		})
	}
}
