package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func Test_allRules_valid_fixtures(t *testing.T) {
	t.Parallel()

	paths, err := filepath.Glob("testdata/valid/*.tf")
	if err != nil || len(paths) == 0 {
		t.Fatalf("no fixtures found: %v", err)
	}

	for _, path := range paths {
		t.Run(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), func(t *testing.T) {
			t.Parallel()

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			helper.AssertIssues(t, helper.Issues{}, runRules(t, string(content), All()...))
		})
	}
}

func runRules(t *testing.T, content string, rules ...tflint.Rule) helper.Issues {
	t.Helper()

	runner := helper.TestRunner(t, map[string]string{"test.tf": content})
	for _, rule := range rules {
		if err := rule.Check(runner); err != nil {
			t.Fatalf("%s: unexpected check error: %s", rule.Name(), err)
		}
	}

	return runner.Issues
}
