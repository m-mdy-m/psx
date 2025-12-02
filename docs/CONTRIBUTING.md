# Contributing to PSX

hey thanks for wanting to help out!

## how to contribute

### reporting bugs

just open an issue and tell me:
- what you did
- what happened
- what you expected to happen
- your OS and version

### suggesting features

open an issue with the `enhancement` label. describe what you want and why it'd be useful.

### code contributions

1. fork the repo
2. make a branch (`git checkout -b fix-something`)
3. do your thing
4. write tests if you can (not mandatory but nice)
5. make sure it builds (`make build`)
6. commit with a decent message
7. push and open a PR

## dev setup

you need:
- go 1.23 or newer
- make (optional but easier)
- git obviously

clone and build:
```bash
git clone https://github.com/m-mdy-m/psx
cd psx
make build
```

run tests:
```bash
make test
```

## code style

i use `gofmt` so just run that. don't worry too much about it.
```bash
gofmt -w .
```

also `golangci-lint` if you have it:
```bash
golangci-lint run
```
## adding new rules

1. go to `internal/rules/`
2. add your rule in the right category folder
3. implement the `Rule` interface:
```go
type Rule interface {
    ID() string
    Check(project *Project) *Result
    Fix(project *Project) error
    Severity() Severity
}
```
4. register it in `registry.go`
5. add tests in `*_test.go`

## commit messages

something like:
```
fix: typo in readme
feat: add python project detection
refactor: clean up rule engine
```

## pull requests

- one feature/fix per PR
- update docs if needed
- add yourself to contributors if you want

## questions?

just ask in issues or email me: bitsgenix@gmail.com

## code of conduct

basically don't be a jerk. see CODE_OF_CONDUCT.md

---

thanks! 🙏