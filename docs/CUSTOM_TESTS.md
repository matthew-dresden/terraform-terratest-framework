# Writing Custom Tests

This guide explains how to write custom tests for your Terraform modules using the Terraform Terratest Framework.

## Overview

While the framework provides common assertions for basic testing, real-world Terraform modules often require custom tests specific to their functionality. The framework supports running custom test functions on all examples in parallel.

## Writing Custom Test Functions

A custom test function is any function that matches the `CustomTestFunc` type:

```go
type CustomTestFunc func(t *testing.T, ctx testctx.TestContext)
```

Inside this function, you can:
- Access the Terraform outputs
- Make assertions about resources
- Call provider APIs or other external services to verify resources
- Check for specific conditions
- Use any of the built-in assertions

Example of a custom test function:

```go
func verifyRenderedFile(t *testing.T, ctx testctx.TestContext) {
    // Get the file path from Terraform outputs
    path := terraform.Output(t, ctx.Terraform, "output_file_path")

    // Verify the file the module rendered actually exists on disk
    info, err := os.Stat(path)
    require.NoError(t, err, "rendered file should exist")
    assert.False(t, info.IsDir(), "rendered path should be a file, not a directory")
}
```

## Running Custom Tests on All Examples

There are two main approaches to running custom tests on all examples:

### Approach 1: Run Examples First, Then Custom Tests

```go
func TestModule(t *testing.T) {
    // Run all examples with default configs
    results := testctx.RunAllExamples(t, "../..", nil)
    
    // Define custom test functions
    verifyFiles := func(t *testing.T, ctx testctx.TestContext) {
        // Custom rendered-file verification logic
    }
    
    verifyPermissions := func(t *testing.T, ctx testctx.TestContext) {
        // Custom permissions verification logic
    }
    
    // Run custom tests on all examples
    testctx.RunCustomTests(t, results, verifyFiles, verifyPermissions)
}
```

### Approach 2: Run Examples and Custom Tests in One Go

```go
func TestModule(t *testing.T) {
    // Define custom test functions
    verifyFiles := func(t *testing.T, ctx testctx.TestContext) {
        // Custom rendered-file verification logic
    }
    
    verifyPermissions := func(t *testing.T, ctx testctx.TestContext) {
        // Custom permissions verification logic
    }
    
    // Run all examples and then run custom tests on each
    testctx.RunAllExamplesWithTests(t, "../..", nil, verifyFiles, verifyPermissions)
}
```

## Example-Specific Tests

You can run different tests for different examples:

```go
func TestModule(t *testing.T) {
    // Run all examples
    results := testctx.RunAllExamples(t, "../..", nil)
    
    // Define a custom test that behaves differently based on the example
    customTest := func(t *testing.T, ctx testctx.TestContext) {
        switch ctx.Config.Name {
        case "basic":
            // Test logic for basic example
            assertions.AssertOutputEquals(t, ctx, "mode", "minimal")
        
        case "advanced":
            // Test logic for advanced example
            assertions.AssertOutputEquals(t, ctx, "mode", "full")
            
            // Additional advanced-specific tests
            name := terraform.Output(t, ctx.Terraform, "resource_name")
            assert.Contains(t, name, "advanced")
        }
    }
    
    // Run the custom test on all examples
    testctx.RunCustomTests(t, results, customTest)
}
```

## Best Practices

1. **Organize Tests by Resource Type**: Group tests by the resource type they're testing (storage, compute, identity, etc.)

2. **Use Helper Functions**: Create reusable helper functions for common verification tasks

3. **Handle Example-Specific Logic**: Use conditionals or switch statements to handle differences between examples

4. **Clean Up Resources**: The framework automatically cleans up Terraform resources, but if your tests create additional resources, clean them up

5. **Use Descriptive Test Names**: When defining subtests, use descriptive names that indicate what's being tested

6. **Error Handling**: Include proper error handling and descriptive assertion messages

## Example: Complete Custom Test

```go
package functional

import (
    "os"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/gruntwork-io/terratest/modules/terraform"
    
    "github.com/matthew-dresden/terraform-terratest-framework/internal/assertions"
    "github.com/matthew-dresden/terraform-terratest-framework/pkg/testctx"
)

func TestFileModule(t *testing.T) {
    // Define custom test functions
    verifyFileExists := func(t *testing.T, ctx testctx.TestContext) {
        path := terraform.Output(t, ctx.Terraform, "output_file_path")
        
        _, err := os.Stat(path)
        require.NoError(t, err, "rendered file should exist")
    }
    
    verifyFileContent := func(t *testing.T, ctx testctx.TestContext) {
        path := terraform.Output(t, ctx.Terraform, "output_file_path")
        
        content, err := os.ReadFile(path)
        require.NoError(t, err)
        assert.Contains(t, string(content), "expected marker")
    }
    
    // Run all examples and custom tests
    testctx.RunAllExamplesWithTests(t, "../..", nil, 
        verifyFileExists, 
        verifyFileContent,
    )
}
```

## Advanced: Test Fixtures

For complex tests, you might want to create test fixtures:

```go
type FileTestFixture struct {
    Path    string
    Content string
}

func setupFileFixture(t *testing.T, ctx testctx.TestContext) *FileTestFixture {
    path := terraform.Output(t, ctx.Terraform, "output_file_path")
    
    content, err := os.ReadFile(path)
    require.NoError(t, err)
    
    return &FileTestFixture{
        Path:    path,
        Content: string(content),
    }
}

func TestFileModuleWithFixture(t *testing.T) {
    results := testctx.RunAllExamples(t, "../..", nil)
    
    for name, ctx := range results {
        t.Run(fmt.Sprintf("FileTests_%s", name), func(t *testing.T) {
            fixture := setupFileFixture(t, ctx)
            
            t.Run("FileExists", func(t *testing.T) {
                // Use fixture to test file existence
            })
            
            t.Run("FileContent", func(t *testing.T) {
                // Use fixture to test file content
            })
        })
    }
}
```

This approach allows you to organize complex tests with shared setup logic.