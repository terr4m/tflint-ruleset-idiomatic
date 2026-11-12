package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_groupingRule(t *testing.T) {
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

  precondition {
    condition     = true
    error_message = "..."
  }

  depends_on = [a]
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_single_line_attributes_without_a_layout",
			content: `
locals {
  alpha = "..."
  beta  = "..."
  gamma = "..."
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_comment_line_separates_items",
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
			name: "valid_object_items_separated",
			content: `
locals {
  tags = {
    a = 1
    b = 2

    c = {
      d = 1
    }
  }
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_object_items_on_one_line",
			content: `
locals {
  alpha = { x = 1, y = { z = 1 } }
}
`,
			want: helper.Issues{},
		},
		{
			name: "valid_object_items_sharing_a_line_with_a_multi_line_item",
			content: `
locals {
  alpha = { x = 1, y = {
    z = 1
  } }
}
`,
			want: helper.Issues{},
		},
		{
			name: "invalid_variable_validation_not_separated",
			content: `
variable "test" {
  description = "..."
  type        = string
  validation {
    condition     = true
    error_message = "..."
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`validation` spans multiple lines and should be separated from `type` by an empty line.",
					Range:   rng(5, 3, 13),
				},
			},
		},
		{
			name: "invalid_multi_line_attribute_after_single_line_attribute",
			content: `
locals {
  alpha = "..."
  beta = {
    a = 1
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`beta` spans multiple lines and should be separated from `alpha` by an empty line.",
					Range:   rng(4, 3, 7),
				},
			},
		},
		{
			name: "invalid_single_line_attribute_after_multi_line_attribute",
			content: `
locals {
  alpha = {
    a = 1
  }
  beta = "..."
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`alpha` spans multiple lines and should be separated from `beta` by an empty line.",
					Range:   rng(3, 3, 8),
				},
			},
		},
		{
			name: "invalid_adjacent_multi_line_blocks",
			content: `
resource "dummy_resource" "example" {
  provisioner "local-exec" {
    command = "echo"
  }
  provisioner "local-exec" {
    command = "echo"
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`provisioner` spans multiple lines and should be separated from `provisioner` by an empty line.",
					Range:   rng(6, 3, 14),
				},
			},
		},
		{
			name: "invalid_nested_block_multi_line_child",
			content: `
resource "dummy_resource" "example" {
  lifecycle {
    ignore_changes = [a]
    precondition {
      condition     = true
      error_message = "..."
    }
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`precondition` spans multiple lines and should be separated from `ignore_changes` by an empty line.",
					Range:   rng(5, 5, 17),
				},
			},
		},
		{
			name: "invalid_object_multi_line_item",
			content: `
locals {
  tags = {
    a = 1
    b = {
      c = 1
    }
    d = 2
  }
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`b` spans multiple lines and should be separated from `a` by an empty line.",
					Range:   rng(5, 5, 6),
				},
				{
					Rule:    NewGroupingRule(),
					Message: "`b` spans multiple lines and should be separated from `d` by an empty line.",
					Range:   rng(5, 5, 6),
				},
			},
		},
		{
			name: "invalid_resource_well_known_grouped_with_other",
			content: `
resource "dummy_resource" "example" {
  count = 1
  alpha = "..."
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`count` should be separated from `alpha` by an empty line (they are not part of the same group).",
					Range:   rng(3, 3, 8),
				},
			},
		},
		{
			name: "invalid_resource_well_known_grouped_with_well_known",
			content: `
resource "dummy_resource" "example" {
  provider = dummy.alt
  count    = 1
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`count` should be separated from `provider` by an empty line (they are not part of the same group).",
					Range:   rng(4, 3, 8),
				},
			},
		},
		{
			name: "invalid_module_header_grouped_with_other",
			content: `
module "example" {
  source = "./example"
  alpha  = "..."
}
`,
			want: helper.Issues{
				{
					Rule:    NewGroupingRule(),
					Message: "`source` should be separated from `alpha` by an empty line (they are not part of the same group).",
					Range:   rng(3, 3, 9),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, map[string]string{
				"test.tf": tt.content,
			})

			if err := NewGroupingRule().Check(runner); err != nil {
				t.Fatalf("unexpected check error: %s", err)
			}

			helper.AssertIssues(t, tt.want, runner.Issues)
		})
	}
}
