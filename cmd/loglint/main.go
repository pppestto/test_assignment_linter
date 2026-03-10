package main

import (
	"github.com/pppestto/test_assignment_linter/addcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(addcheck.Analyzer)
}
