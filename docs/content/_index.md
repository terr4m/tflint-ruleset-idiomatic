---
title: Overview
type: docs
---

## TFLint Ruleset Idiomatic

Welcome to the rulebook for [TFLint Ruleset Idiomatic](https://github.com/terr4m/tflint-ruleset-idiomatic).

### Usage

To use the TFLint Ruleset Idiomatic plugin, add the following configuration to your `.tflint.hcl` file.

```hcl
plugin "idiomatic" {
  source    = "github.com/terr4m/tflint-ruleset-idiomatic"
  signature = "attestation"
  version   = "0.0.1"
  enabled   = true
}
```

### Configuration

There is currently no configuration required for this ruleset. All rules are enabled by default with the default settings. If you wish to disable an individual rule, you can do so in your `.tflint.hcl` file.

```hcl
rule "idiomatic_sorting" {
  enabled = false
}
```

### Rules

See [Rules](./rules) for more details.
