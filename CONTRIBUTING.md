# Contributing

Thanks for helping improve this project — we welcome contributions. Please follow these guidelines.

1. Fork the repo and create a feature branch from `main`.
2. Keep changes small and focused. One logical change per PR.
3. Run formatting and linting before opening a PR:

```
gofmt -w .
go vet ./...
go test ./... -v
```

4. Commit messages should be clear and reference the issue when applicable.
5. Add tests for new features and ensure existing tests pass.

PR Checklist

- [ ] Code is formatted (`gofmt`).
- [ ] Tests added/updated and passing.
- [ ] No secrets or credentials committed (`credentials.json` must never be committed).

If you're unsure about any change, open an issue to discuss it first.
