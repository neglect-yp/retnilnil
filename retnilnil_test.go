package retnilnil_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/neglect-yp/retnilnil"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, retnilnil.Analyzer, "a")
}
