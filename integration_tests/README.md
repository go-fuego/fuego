# Integration Tests for Fuego

This directory contains integration tests for the Fuego framework using [Hurl](https://hurl.dev/).

## Prerequisites

- Go 1.21 or later
- Hurl (install via `winget install Orange-OpenSource.Hurl` on Windows or `brew install hurl` on macOS)

## Running Tests

1. Start the Fuego server in one terminal:

   ```bash
   cd ../examples/basic
   go run .
   ```

2. Run the tests in another terminal using:

   Windows:

   ```powershell
   .\run_tests.ps1
   ```

   Unix/Linux/Git Bash:

   ```bash
   # First time only
   chmod +x run_tests.sh

   # Run tests
   ./run_tests.sh
   ```

## Test Structure

- `basic.hurl`: Tests basic functionality including:
  - Successful POST request with JSON response
  - Validation errors
  - Header validation
  - Standard GET endpoint

## Adding New Tests

1. Create a new `.hurl` file in this directory
2. Add your test cases following the Hurl syntax
3. Update the test runners if needed to include your new test file

## Troubleshooting

- If tests fail with connection errors, ensure no other process is using port 8088
- If using the shell scripts, ensure they have the correct line endings for your OS
- For Windows users, Git Bash is recommended when using the bash script
