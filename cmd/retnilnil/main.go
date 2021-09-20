package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/neglect-yp/retnilnil"
)

func main() {
	singlechecker.Main(retnilnil.Analyzer)
}
