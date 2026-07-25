# Contributing to HLFA

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/bug.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Commit: `git commit -am 'Add your feature'`
6. Push: `git push origin feature/your-feature`
7. Open a Pull Request

## Development Setup

```bash
# Install dependencies
go mod download

# Setup development environment
make docker-up

# Run tests
make test
```

## Code Style

- Go: Follow [Effective Go](https://golang.org/doc/effective_go)
- Python: PEP 8
- TypeScript: ESLint configuration

## Testing Requirements

- All new features must include tests
- Minimum 80% code coverage
- Integration tests for API changes

## Commit Messages

Use conventional commits:
```
feat: add new scanner
fix: resolve race condition in worker pool
docs: update API documentation
test: add integration tests
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new functionality
3. Ensure all tests pass: `make test`
4. Update CHANGELOG.md
5. Request review from maintainers

## Reporting Bugs

Include:
- Description of the bug
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment information
- Logs/error messages

## Feature Requests

Include:
- Description of the feature
- Use case/motivation
- Proposed implementation approach
- Potential risks or considerations

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
