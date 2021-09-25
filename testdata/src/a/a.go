package a

import "errors"

type T struct {
	I int
}

type IF interface{}

func f() (*T, error) {
	return &T{I: 1}, nil
}

func f2() (*T, error) {
	return nil, errors.New("error")
}

func f3() (*T, error) {
	return nil, nil // want "return nil, nil"
}

func f4() (T, error) {
	return T{}, nil
}

func f5() (IF, error) {
	return &T{}, nil
}

func f6() (IF, error) {
	return &T{}, errors.New("error")
}

func f7() (IF, error) {
	return nil, nil // want "return nil, nil"
}

func f8() ([]T, error) {
	return nil, nil
}

func f9() (map[string]T, error) {
	return nil, nil
}
