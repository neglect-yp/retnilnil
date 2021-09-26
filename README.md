# retnilnil

[![Test](https://github.com/neglect-yp/retnilnil/actions/workflows/test.yml/badge.svg)](https://github.com/neglect-yp/retnilnil/actions/workflows/test.yml)

Retnilnil is a static analysis tool for Golang that detects `return nil, nil` in functions with `(*T, error)` as the return type.

```go
func f() (*T, error) {
	return nil, nil // retnilnil detects this
}
```

Retnilnil ignores a code which has an ignore comment.

```go
func f() (*T, error) {
	//lint:ignore retnilnil reason
	return nil, nil // retnilnil doesn't detect this
}
```

Retnilnil also ignores a function that has a comment includes `nil, nil`.

```go
// f always returns `nil, nil`
func f() (*T, error) {
    return nil, nil // retnilnil doesn't detect this
}
```
