---
title: idiomatic_ordering
description: Ensures that specific well-known block attributes and nested blocks are ordered in a consistent manner.
---

## Rule `{{< title >}}`

{{< description >}}

- Well known attributes and nested blocks should follow a specific order within their block.
- Only the position of an entry relative to the other entries is checked. Empty lines between entries are not part of this rule, they are covered by [`idiomatic_grouping`](./grouping/).
- Only top-level `variable`, `output`, `resource`, `data` and `module` blocks are checked. Other blocks, and blocks nested inside them, are ignored.
- Files written in JSON syntax (`.tf.json`) are not checked.

The `...` entries below stand for any other attributes or blocks, which can appear in any order here. Their order is the concern of [`idiomatic_sorting`](./sorting/).

### Variables

1. `description`
2. `type`
3. `nullable`
4. `sensitive`
5. `const`
6. `ephemeral`
7. `default`
8. `deprecated`
9. `validation`

### Outputs

1. `description`
2. `type`
3. `sensitive`
4. `ephemeral`
5. `deprecated`
6. `value`
7. `precondition`
8. `depends_on`

### Resources/Data Sources

1. `count` or `for_each`
2. `provider`
3. `...` (other attributes)
4. `connection`
5. `provisioner`
6. `lifecycle`
7. `depends_on`

### Modules

1. `source`
2. `version`
3. `ignore_nested_deprecations`
4. `count` or `for_each`
5. `providers`
6. `...` (other attributes)
7. `depends_on`

### Severity

This rule reports issues at the `warning` level.

### Configuration

There are no configuration options for this rule. If you wish to disable this rule, you can do so in your `.tflint.hcl` file.
