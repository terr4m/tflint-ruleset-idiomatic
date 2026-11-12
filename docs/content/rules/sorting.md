---
title: idiomatic_sorting
description: Ensures that block and complex type attributes are sorted in a consistent alphanumeric order.
---

## Rule `{{< title >}}`

{{< description >}}

- Block or complex type attributes should be sorted alphanumerically within their group.
- A group is defined as a set of attributes that are on the same line or on adjacent lines. An empty line or a comment line starts a new group. This rule only compares neighbours inside a group. Where empty lines are required between attributes, for example around multi-line values, is the concern of [`idiomatic_grouping`](./grouping/).
- Sorting is by byte value, so uppercase letters sort before lowercase letters (`Name` comes before `age`).
- Object keys are compared as strings, the same way Terraform stores them, so the key `443` comes before `80`.
- Attributes that have a different fixed position in [`idiomatic_ordering`](./ordering/) (for example `source` and `version` in a `module`, or `description` and `type` in a `variable`) are not compared with each other, because their order is set by that rule. Attributes that are not listed there are sorted as normal.
- Attributes inside nested blocks, and object values, anywhere in a file are checked, including inside `variable`, `output` and `locals` blocks. Nested blocks themselves are not reordered.

### Severity

This rule reports issues at the `warning` level.

### Configuration

There are no configuration options for this rule. If you wish to disable this rule, you can do so in your `.tflint.hcl` file.

### Example

```hcl
resource "aws_instance" "web" {
  ami           = "ami-0123456789abcdef0"
  instance_type = "t3.micro"
  subnet_id     = var.subnet_id

  tags = {
    Environment = "dev"
    Name        = "web"
  }
}
```
