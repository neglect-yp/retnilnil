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

func fNaked() (t *T, err error) {
	return
}

func fIf() (*T, error) {
	if true {
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fElse() (*T, error) {
	if true {
		return &T{}, nil
	} else {
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fFor() (*T, error) {
	for {
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fRange() (*T, error) {
	for range []string{""} {
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fSwitch() (*T, error) {
	switch true {
	case true:
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fDefault() (*T, error) {
	switch true {
	default:
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fTypeSwitch() (*T, error) {
	var v interface{} = 1
	switch v.(type) {
	case int:
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fBlock() (*T, error) {
	{
		return nil, nil // want "return nil, nil"
	}

	return &T{}, nil
}

func fIgnore() (*T, error) {
	//lint:ignore retnilnil reason
	return nil, nil
}

func fIgnoreWithoutReason() (*T, error) {
	//lint:ignore retnilnil
	return nil, nil // want "return nil, nil"
}

// fDocumented always returns `nil, nil`
func fDocumented() (*T, error) {
	return nil, nil
}

func fFuncLit() {
	f := func() (*T, error) {
		return nil, nil // want "return nil, nil"
	}
	f()
}

func fFuncLit2() {
	func() (*T, error) {
		return nil, nil // want "return nil, nil"
	}()
}

var fFuncLit3 = func() (*T, error) {
	return nil, nil // want "return nil, nil"
}
