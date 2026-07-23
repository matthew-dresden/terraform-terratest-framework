# Example-Specific Tests

This guide explains how to write tests for specific examples in your Terraform module.

## Overview

The Terraform Terratest Framework supports two testing approaches:

1. **Centralized Testing**: Running all examples from a single test file
2. **Example-Specific Testing**: Dedicated test files for each example

This document focuses on the second approach, which is useful when:
- Each example requires unique test logic
- Examples are complex and warrant dedicated test files
- You want to organize tests by example

## Directory Structure

For example-specific testing, your module should follow this structure:

```
terraform-module/
├── examples/
│   ├── example1/
│   │   ├── main.tf
│   │   └── ...
│   ├── example2/
│   │   ├── main.tf
│   │   └── ...
│   └── ...
├── tests/
│   ├── common/
│   │   └── helpers.go
│   ├── example1/
│   │   └── module_test.go
│   ├── example2/
│   │   └── module_test.go
│   └── ...
└── ...
```

## Writing Example-Specific Tests

In each example's test file (e.g., `tests/example1/module_test.go`):

```go
package example1

import (
	"testing"

	"github.com/matthew-dresden/terraform-terratest-framework/internal/assertions"
	"github.com/matthew-dresden/terraform-terratest-framework/pkg/testctx"
)

func TestExample(t *testing.T) {
	// Run just this specific example
	ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{
		Name: "example1",
		ExtraVars: map[string]interface{}{
			"region": "us-west-2",
		},
	})
	
	// Run example-specific assertions
	assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
	
	// Custom verification for this example
	verifyResources(t, ctx)
}

func verifyResources(t *testing.T, ctx testctx.TestContext) {
	// Example-specific verification logic
	// ...
}
```

## Sharing Common Test Logic

You can create shared test helpers in the `tests/common` directory:

```go
// tests/common/helpers.go
package common

import (
	"testing"
	"github.com/matthew-dresden/terraform-terratest-framework/pkg/testctx"
)

func VerifyS3Bucket(t *testing.T, ctx testctx.TestContext) {
	// Common S3 bucket verification logic
	// ...
}
```

Then import and use these helpers in your example-specific tests:

```go
package example1

import (
	"testing"

	"github.com/matthew-dresden/terraform-terratest-framework/pkg/testctx"
	"terraform-module/tests/common"
)

func TestExample(t *testing.T) {
	ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{})
	
	// Use common test helpers
	common.VerifyS3Bucket(t, ctx)
	
	// Example-specific tests
	// ...
}
```

## Best Practices

1. **Use Relative Paths**: Always use relative paths (`"../../"`) to reference the module root

2. **Organize by Resource Type**: Group test functions by the resources they verify

3. **Share Common Logic**: Put reusable test functions in the `tests/common` package

4. **Use Descriptive Test Names**: Name tests based on what they're verifying

5. **Test One Example Per File**: Keep each test file focused on a single example

## Example: Complete Test File

```go
package advanced

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/gruntwork-io/terratest/modules/terraform"
	
	"github.com/matthew-dresden/terraform-terratest-framework/internal/assertions"
	"github.com/matthew-dresden/terraform-terratest-framework/pkg/testctx"
	"terraform-module/tests/common"
)

func TestAdvanced(t *testing.T) {
	// Run the specific example
	ctx := testctx.RunSingleExample(t, "../../", "advanced", testctx.TestConfig{
		Name: "advanced",
		ExtraVars: map[string]interface{}{
			"name": "test-name",
			"mode": "full",
		},
	})
	
	// Verify outputs
	assertions.AssertOutputEquals(t, ctx, "name", "test-name")
	
	// Verify the rendered file
	verifyRenderedFile(t, ctx)
	
	// Use common test helpers
	common.VerifyTags(t, ctx)
}

func verifyRenderedFile(t *testing.T, ctx testctx.TestContext) {
	path := terraform.Output(t, ctx.Terraform, "output_file_path")
	
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test-name")
}
```