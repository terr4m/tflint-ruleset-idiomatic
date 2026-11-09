package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_orderingRule(t *testing.T) {
	t.Parallel()

	rng := func(line, startCol, endCol int) hcl.Range {
		return hcl.Range{
			Filename: "test.tf",
			Start:    hcl.Pos{Line: line, Column: startCol},
			End:      hcl.Pos{Line: line, Column: endCol},
		}
	}

	for _, tt := range []struct {
		name    string
		content string
		want    helper.Issues
	}{
		{
			name: "valid_module",
			content: `
module "example" {
  source                     = "./example"
  version                    = "1.0.0"
  ignore_nested_deprecations = true

  for_each = var.items

  providers = {
    aws = aws.alt
  }

  alpha = "..."
  beta  = "..."

  depends_on = [a]
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_resource",
			content: `
resource "dummy_resource" "example" {
  count = 1

  provider = dummy.alt

  alpha = "..."
  beta  = "..."

  connection {
    type = "ssh"
  }

  provisioner "local-exec" {
    command = "echo"
  }

  provisioner "local-exec" {
    command = "echo"
  }

  lifecycle {
    ignore_changes = [alpha]
  }

  depends_on = [a]
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_variable",
			content: `
variable "test" {
  description = "..."
  type        = string
  default     = "a"
  validation {
    condition     = true
    error_message = "..."
  }
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_output",
			content: `
output "test" {
  description = "..."
  value       = "a"
  depends_on  = [a]
}
`,
			want: helper.Issues{},
		},
		{
			name: "unconfigured_block_keys_have_no_special_position",
			content: `
locals {
  count      = 1
  depends_on = "..."
  source     = "..."
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_comment_line_between_items",
			content: `
resource "dummy_resource" "example" {
  count = 1
  # alpha is documented here
  alpha = "..."
}
`,
			want: helper.Issues{},
		},
		{
			name: "invalid_module_source_after_other",
			content: `
module "example" {
  alpha = "..."

  source = "./example"
}
`,
			want: helper.Issues{
				{Rule: NewOrderingRule(), Message: "`source` should appear before `alpha`.", Range: rng(5, 3, 9)},
			},
		},
		{
			name: "invalid_module_source_after_count",
			content: `
module "example" {
  count = 1

  source = "./example"
}
`,
			want: helper.Issues{
				{Rule: NewOrderingRule(), Message: "`source` should appear before `count`.", Range: rng(5, 3, 9)},
			},
		},
		{
			name: "invalid_module_depends_on_before_other",
			content: `
module "example" {
  source = "./example"

  depends_on = [a]

  alpha = "..."
}
`,
			want: helper.Issues{
				{Rule: NewOrderingRule(), Message: "`depends_on` should appear after `alpha`.", Range: rng(5, 3, 13)},
			},
		},
		{
			name: "invalid_resource_lifecycle_after_depends_on",
			content: `
resource "dummy_resource" "example" {
  alpha = "..."

  depends_on = [a]

  lifecycle {
    ignore_changes = [alpha]
  }
}
`,
			want: helper.Issues{
				{Rule: NewOrderingRule(), Message: "`depends_on` should appear after `lifecycle`.", Range: rng(5, 3, 13)},
			},
		},
		{
			name: "invalid_data_count_after_other",
			content: `
data "dummy_data" "example" {
  alpha = "..."

  count = 1
}
`,
			want: helper.Issues{
				{Rule: NewOrderingRule(), Message: "`count` should appear before `alpha`.", Range: rng(5, 3, 8)},
			},
		},
		{
			name: "valid_spacing_is_left_to_grouping",
			content: `
resource "dummy_resource" "example" {
  count = 1
  alpha = "..."
}
`,
			want: helper.Issues{},
		},
		{
			name: "invalid_resource_well_known_out_of_order",
			content: `
resource "dummy_resource" "example" {
  provider = dummy.alt
  count    = 1
}
`,
			want: helper.Issues{
				{
					Rule:    NewOrderingRule(),
					Message: "`count` should appear before `provider`.",
					Range:   rng(4, 3, 8),
				},
			},
		},
		{
			name: "invalid_variable_type_before_description",
			content: `
variable "test" {
  type        = string
  description = "..."
}
`,
			want: helper.Issues{
				{Rule: NewOrderingRule(), Message: "`description` should appear before `type`.", Range: rng(4, 3, 14)},
			},
		},
		{
			name: "invalid_output_depends_on_before_value",
			content: `
output "test" {
  depends_on = [a]
  value      = "a"
}
`,
			want: helper.Issues{
				{Rule: NewOrderingRule(), Message: "`value` should appear before `depends_on`.", Range: rng(4, 3, 8)},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, map[string]string{
				"test.tf": tt.content,
			})

			if err := NewOrderingRule().Check(runner); err != nil {
				t.Fatalf("unexpected check error: %s", err)
			}

			helper.AssertIssues(t, tt.want, runner.Issues)
		})
	}
}
