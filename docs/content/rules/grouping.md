---
title: idiomatic_grouping
description: Ensures that block and complex type attributes are grouped in a consistent manner.
---

## Rule `{{< title >}}`

{{< description >}}

- Attributes, nested blocks and object keys that span multiple lines must be separated by an empty line from the entries next to them. This applies to single line entries and to other multi-line entries.
- In top-level `resource`, `data`, `module`, `variable` and `output` blocks, well-known attributes and nested blocks must be separated by an empty line from any entry that is not part of the same group (see below).
- A comment line between two entries counts as an empty line.
- Entries that are on the same line, such as the keys in `{ a = 1, b = 2 }`, are not compared with each other.
- Nested blocks and object values anywhere in a file are checked, including inside `variable`, `output` and `locals` blocks. The groups below only apply to top-level blocks.
- Files written in JSON syntax (`.tf.json`) are not checked.

The position of entries is the concern of [`idiomatic_ordering`](./ordering/) and their alphanumeric order is the concern of [`idiomatic_sorting`](./sorting/).

### Groups

Entries in the same group can be written on adjacent lines, as long as none of them spans multiple lines. Entries that are not listed below are other attributes, which share a single group.

#### Variables and Outputs

All of the well-known attributes and blocks listed in [`idiomatic_ordering`](./ordering/) form one group, and any other attribute is separated from them by an empty line. For example, `description`, `type` and `default` can be written together, but a `validation` block spans multiple lines so it is separated from the attributes above it:

```hcl
variable "port" {
  description = "The port to listen on."
  type        = number
  default     = 8080

  validation {
    condition     = var.port > 0 && var.port < 65536
    error_message = "The port must be between 1 and 65535."
  }
}
```

#### Resources/Data Sources

Each of these is its own group, separated from the others by an empty line:

1. `count` or `for_each`
2. `provider`
3. `...` (other attributes)
4. `connection`
5. `provisioner`
6. `lifecycle`
7. `depends_on`

#### Modules

1. `source`, `version` and `ignore_nested_deprecations` (one group)
2. `count` or `for_each`
3. `providers`
4. `...` (other attributes)
5. `depends_on`

### Multi-line Entries

Multi-line entries are separated from their neighbours wherever they appear, including in blocks that have no groups such as `locals`:

```hcl
resource "aws_instance" "web" {
  ami           = "ami-0123456789abcdef0"
  instance_type = "t3.micro"

  tags = {
    Name = "web"
  }
}
```

### Severity

This rule reports issues at the `warning` level.

### Configuration

There are no configuration options for this rule. If you wish to disable this rule, you can do so in your `.tflint.hcl` file.
