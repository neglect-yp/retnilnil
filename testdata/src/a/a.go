package a

import "errors"

type T struct {
	I int
}

func f() (*T, error) {
	return &T{I: 1}, nil
}

func f2() (*T, error) {
	return nil, nil // want "return nil, nil"
}

func f3() (*T, error) {
	return nil, errors.New("error")
}

func f4() (T, error) {
	return T{}, nil
}
