package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"metric-alert/cmd/staticlint/exitanalyzer"
)

var ignoreChecks = map[string]bool{
	"ST1000": true,
}

type checks struct {
	Analyzers []*analysis.Analyzer
}

func main() {
	multichecker.Main(getAnalyzers()...)
}

func getAnalyzers() []*analysis.Analyzer {
	c := checks{}

	c.addAnalyzers(staticcheck.Analyzers)
	c.addAnalyzers(stylecheck.Analyzers)

	c.setStandardAnalyzers()

	return c.Analyzers
}

func (c *checks) addAnalyzers(analyzers []*lint.Analyzer) {
	for _, v := range analyzers {
		if !ignoreChecks[v.Analyzer.Name] {
			c.Analyzers = append(c.Analyzers, v.Analyzer)
		}
	}
}

func (c *checks) setStandardAnalyzers() {
	c.Analyzers = append(c.Analyzers,
		exitanalyzer.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		atomicalign.Analyzer,
		composite.Analyzer,
		buildtag.Analyzer,
	)
}
