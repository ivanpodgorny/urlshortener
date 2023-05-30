package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/analysis/lint"
)

// List реализует методы для управления списком анализаторов
// и получения списка проверок.
type List struct {
	analyzers map[string]*analysis.Analyzer
}

// NewList возвращает указатель на новый экземпляр List.
func NewList() *List {
	return &List{
		analyzers: map[string]*analysis.Analyzer{},
	}
}

// GetChecks возвращает те из добавленных анализаторов, которые встречаются в names.
func (l *List) GetChecks(names ...string) []*analysis.Analyzer {
	checks := make([]*analysis.Analyzer, 0, len(names))
	for _, n := range names {
		if a, ok := l.analyzers[n]; ok {
			checks = append(checks, a)
		}
	}

	return checks
}

// AddAnalyzers добавляет в список стандартные анализаторы.
func (l *List) AddAnalyzers(analyzers ...*analysis.Analyzer) *List {
	for _, a := range analyzers {
		l.analyzers[a.Name] = a
	}

	return l
}

// AddStaticCheckAnalyzers добавляет в список анализаторы staticcheck.
//
// https://staticcheck.io/docs/checks/
func (l *List) AddStaticCheckAnalyzers(analyzers ...*lint.Analyzer) *List {
	for _, a := range analyzers {
		l.analyzers[a.Analyzer.Name] = a.Analyzer
	}

	return l
}
