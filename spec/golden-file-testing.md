# Golden file testing

We use golden file testing to make sure that we don't regress on the outputs of the CLI.

## Overview

Golden file testing is a technique where we:

1. Generate expected outputs (golden files) from our CLI commands
2. Store these golden files in version control
3. Run tests that compare actual CLI output against these golden files
4. Fail tests when outputs differ from the expected golden files

This approach helps us detect unintended changes in CLI behavior or output formatting.

## Directory Structure

```
project/
├── cmd/                  # CLI command implementations
├── tests/
│   ├── golden/          # Directory containing golden files
│   │   ├── command1/    # Golden files for command1
│   │   │   ├── case1.golden
│   │   │   └── case2.golden
│   │   ├── command2/    # Golden files for command2
│   │   │   └── case1.golden
│   └── golden_test.go   # Test file that runs and validates against golden files
└── scripts/
    └── generate_golden.sh  # Script to generate/update golden files
```

## Generating Golden Files

To generate or update golden files, run:

```bash
./scripts/generate_golden.sh
```

This script:

1. Builds the latest version of the CLI
2. Executes a predefined set of example commands
3. Captures their outputs to the appropriate files in `tests/golden/`

Only run this script when you intentionally want to update the expected behavior of the CLI. Each update should be carefully reviewed in your commit.

## Test Implementation

In `golden_test.go`, we implement tests that:

1. Build the CLI binary
2. Execute the same set of commands used to generate golden files
3. Compare the output to the corresponding golden file
4. Fail if there are any differences

Example test implementation:

```go
func TestGoldenFiles(t *testing.T) {
    // Find all golden files
    goldenFiles, err := filepath.Glob("golden/**/*.golden")
    if err != nil {
        t.Fatalf("Failed to find golden files: %v", err)
    }

    for _, goldenFile := range goldenFiles {
        t.Run(goldenFile, func(t *testing.T) {
            // Extract command and args from golden file path
            command := extractCommandFromPath(goldenFile)

            // Run the command
            cmd := exec.Command("../bin/cli", command...)
            output, err := cmd.CombinedOutput()
            if err != nil {
                t.Fatalf("Command failed: %v", err)
            }

            // Read expected output
            expected, err := ioutil.ReadFile(goldenFile)
            if err != nil {
                t.Fatalf("Failed to read golden file: %v", err)
            }

            // Compare outputs
            if !bytes.Equal(output, expected) {
                t.Errorf("Output does not match golden file.\nExpected:\n%s\n\nGot:\n%s",
                    string(expected), string(output))
            }
        })
    }
}
```

## CI Integration

These tests should be run as part of your continuous integration pipeline to catch unintended changes to CLI output.

## Updating Golden Files

When you intentionally change CLI behavior or output format:

1. Run `./scripts/generate_golden.sh` to update golden files
2. Review the changes to ensure they match your intended modifications
3. Commit the updated golden files along with your code changes
4. Include in your PR description which golden files were updated and why

## Best Practices

1. Golden files should have clear, descriptive names
2. Each golden file should test a specific use case or command combination
3. Include both standard and error outputs in your tests
4. Consider testing with different flag combinations and edge cases
5. For commands with varying output (timestamps, random IDs, etc.), use test hooks or regex patterns to handle dynamic content

## Troubleshooting

If tests fail with differences in line endings (CRLF vs LF), ensure consistent line endings in your project configuration and golden file generation script.
