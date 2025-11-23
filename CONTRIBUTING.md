# Contributing to Reverse Proxy

First off, thank you for considering contributing to this project! 🎉

## Code of Conduct

This project and everyone participating in it is expected to follow our Code of Conduct. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples to demonstrate the steps**
- **Describe the behavior you observed and what you expected**
- **Include logs and error messages**
- **Specify your environment** (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Use a clear and descriptive title**
- **Provide a detailed description of the suggested enhancement**
- **Explain why this enhancement would be useful**
- **List any similar features in other projects**

### Pull Requests

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Run tests (`go test ./...`)
6. Run linter (`golangci-lint run`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for container testing)
- golangci-lint (for linting)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/bunnydevv/reverse-proxy.git
cd reverse-proxy

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Build the project
go build -o reverse-proxy .

# Run linter
golangci-lint run
```

### Project Structure

```
.
├── main.go                 # Application entry point
├── config/
│   ├── config.go          # Configuration management
│   └── config_test.go     # Configuration tests
├── proxy/
│   ├── proxy.go           # Main proxy logic
│   ├── proxy_test.go      # Proxy tests
│   ├── load_balancer.go   # Load balancing algorithms
│   ├── load_balancer_test.go
│   ├── health_check.go    # Health check implementation
│   └── health_check_test.go
└── internal/
    ├── logger/            # Structured logging
    └── metrics/           # Metrics collection
```

## Coding Standards

### Go Style Guide

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Write clear, idiomatic Go code
- Add comments for exported functions and types

### Testing

- Write unit tests for all new code
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

Example test structure:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests liberally

### Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported symbols
- Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/)

## Release Process

1. Update version in appropriate files
2. Update CHANGELOG.md
3. Create a git tag (`git tag -a v1.0.0 -m "Release v1.0.0"`)
4. Push the tag (`git push origin v1.0.0`)
5. GitHub Actions will automatically build and publish the release

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

Thank you for contributing! 🚀