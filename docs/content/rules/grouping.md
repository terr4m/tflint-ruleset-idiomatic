---
title: idiomatic_grouping
description: Ensures that block and complex type attributes are grouped in a consistent manner.
---

## Rule `{{< title >}}`

{{< description >}}

- Entries are attributes, nested blocks and object keys. Entries that are next to each other without an empty line between them are grouped.
- An entry that spans multiple lines must not be grouped with another entry.
- A well-known entry of a top-level block, as listed in [`idiomatic_ordering`](./ordering/), must not be grouped with an entry that belongs to a different group.
- A comment line between two entries counts as an empty line, and entries on the same line are not compared.
- Entries anywhere in a file are checked, including in nested blocks and object values.

### Severity

This rule reports issues at the `warning` level.

### Configuration

There are no configuration options for this rule. If you wish to disable this rule, you can do so in your `.tflint.hcl` file.

### Example

```hcl
resource "aws_instance" "web" {
  ami           = "ami-0123456789abcdef0"
  instance_type = "t3.micro"

  tags = {
    Name = "web"
  }
}
```

### Groups

Each well-known entry is its own group, with these exceptions:

- `source`, `version` and `ignore_nested_deprecations` in a `module` share a group.
- The well-known entries of a `variable` share a group, as do those of an `output`.
- Entries that are not well-known share a group.

The position of entries is the concern of [`idiomatic_ordering`](./ordering/) and their alphanumeric order is the concern of [`idiomatic_sorting`](./sorting/).
