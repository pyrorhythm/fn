# Contributing

## Getting Started

```bash
git clone https://github.com/pyrorhythm/fn
cd fn
go test ./...
```

## Guidelines

- Keep it simple. No over-engineering.
- Write tests. Aim for >90% coverage.
- Run `go fmt` and `go vet` before committing.
- One feature per PR.

## Pull Requests

1. Fork the repo
2. Create a branch (`git checkout -b feature/thing`)
3. Make changes
4. Run tests (`go test -cover ./...`)
5. Push and open PR

## Code Style

- No unnecessary abstractions
- Prefer explicit over clever
- Comments only where logic isn't obvious
- Follow existing patterns in the codebase

## Reporting Issues

Open an issue with:
- What you expected
- What happened
- Minimal reproduction if possible
