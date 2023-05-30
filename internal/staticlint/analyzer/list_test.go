package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/stylecheck"
)

func TestList(t *testing.T) {
	l := NewList()
	l.AddAnalyzers(shift.Analyzer, structtag.Analyzer)
	assert.Equal(t, []*analysis.Analyzer{structtag.Analyzer}, l.GetChecks("printf", "structtag"))

	l.AddAnalyzers(printf.Analyzer)
	assert.Equal(t, []*analysis.Analyzer{printf.Analyzer, structtag.Analyzer}, l.GetChecks("printf", "structtag"))

	st1000, st1000Orig := getStaticCheckAnalyzer("ST1000", stylecheck.Analyzers)
	require.NotNil(t, st1000)
	assert.Empty(t, l.GetChecks("ST1000"))
	l.AddStaticCheckAnalyzers(st1000Orig)
	assert.Equal(t, []*analysis.Analyzer{st1000}, l.GetChecks("ST1000"))
}

func getStaticCheckAnalyzer(name string, analyzers []*lint.Analyzer) (*analysis.Analyzer, *lint.Analyzer) {
	for _, a := range analyzers {
		if a.Analyzer.Name == name {
			return a.Analyzer, a
		}
	}

	return nil, nil
}
