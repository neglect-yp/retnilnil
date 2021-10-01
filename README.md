# retnilnil

[![Go Reference](https://pkg.go.dev/badge/github.com/neglect-yp/retnilnil.svg)](https://pkg.go.dev/github.com/neglect-yp/retnilnil)
[![Test](https://github.com/neglect-yp/retnilnil/actions/workflows/test.yml/badge.svg)](https://github.com/neglect-yp/retnilnil/actions/workflows/test.yml)

`retnilnil` is a static analysis tool for Golang that detects `return nil, nil` in functions with `(*T, error)` as the return type.

```go
func f() (*T, error) {
	return nil, nil // retnilnil detects this
}
```

`retnilnil` ignores a code which has an ignore comment.

```go
func f() (*T, error) {
	//lint:ignore retnilnil reason
	return nil, nil // retnilnil doesn't detect this
}
```

`retnilnil` also ignores a function that has a comment includes `nil, nil`.

```go
// f always returns `nil, nil`
func f() (*T, error) {
	return nil, nil // retnilnil doesn't detect this
}
```

## Install

You can get `retnilnil` by go install command (Go 1.16 and higher).

```
$ go install github.com/neglect-yp/retnilnil/cmd/retnilnil@v0.1.0
```

## How to use

```
$ retnilnil ./...
```
