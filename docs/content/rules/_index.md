---
title: Rules
weight: 1
---

## Rules

The following rules are provided by this ruleset and only apply to `*.tf` files. Files written in JSON syntax (`*.tf.json`) are not checked. Individual rules can be disabled in your `.tflint.hcl` file if desired.

{{< children >}}

### Modules

All of the rules only check the root module, which is the directory that TFLint is run against. TFLint doesn't support these rules evaluating child modules, so the files of modules that are called with a `module` block are skipped, whatever the value of `--call-module-type`. To check every module in a repository, run TFLint with `--recursive` so that each directory is checked as a root module.

### How the rules fit together

The rules each cover one aspect of layout and are designed to agree with each other, so a block that is laid out correctly passes all of them.

- [`idiomatic_ordering`](./ordering/) sets the position of well-known attributes and nested blocks.
- [`idiomatic_grouping`](./grouping/) sets where empty lines are required between them.
- [`idiomatic_sorting`](./sorting/) sets the alphanumeric order of attributes within a group.

For example, in a `variable` the well-known attributes follow their fixed order and the `validation` blocks come last (ordering), and the `validation` blocks are separated from the single line attributes above them by an empty line (grouping):

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

Attributes that have no fixed position are sorted alphanumerically within their group (sorting). In this `resource` the `count` comes first, the multi-line `tags` and `lifecycle` entries are separated from their neighbours, and `ami`, `instance_type` and `subnet_id` are in alphanumeric order:

```hcl
resource "aws_instance" "web" {
  count = 2

  ami           = "ami-0123456789abcdef0"
  instance_type = "t3.micro"
  subnet_id     = var.subnet_id

  tags = {
    Name = "web"
  }

  lifecycle {
    create_before_destroy = true
  }
}
```
