/*
staticlint создает multichecker и выполняет статические проверки кода заданными анализаторами.

# Анализаторы по умолчанию

По умолчанию для создания multichecker используются:

1) все анализаторы из golang.org/x/tools/go/analysis/passes,

2) анализаторы staticcheck SA, ST, S, U,

3) анализаторы github.com/ivanpodgorny/urlshortener/pkg/analyzer/osexit,
github.com/kisielk/errcheck/errcheck, github.com/timakin/bodyclose/passes/bodyclose.

# Файл конфигурации

Список анализаторов для создания multichecker по умолчанию можно переопределить,
указав массив имен анализаторов в файле staticlint.json в директории исполняемого
файла. Доступны любые анализаторы staticcheck, golang.org/x/tools/go/analysis/passes,
github.com/ivanpodgorny/urlshortener/pkg/analyzer, github.com/kisielk/errcheck/errcheck,
github.com/timakin/bodyclose/passes/bodyclose.

# Пример файла конфигурации

	{
	    "analyzers": ["printf", "shadow", "SA1019", "ST1003", "osexit"]
	}

# Пример использования

Запуск всех анализаторов

	staticlint ./...

Запуск определенных анализаторов

	staticlint -printf -SA1019 ./...
*/
package main

import (
	"log"

	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"honnef.co/go/tools/unused"

	"github.com/ivanpodgorny/urlshortener/internal/staticlint/analyzer"
	"github.com/ivanpodgorny/urlshortener/internal/staticlint/config"
)

func main() {
	cfg, err := config.NewBuilder().LoadFile().Build()
	if err != nil {
		log.Fatal(err)
	}

	checks := analyzer.NewList().
		AddAnalyzers(analyzer.Analyzers...).
		AddStaticCheckAnalyzers(staticcheck.Analyzers...).
		AddStaticCheckAnalyzers(simple.Analyzers...).
		AddStaticCheckAnalyzers(quickfix.Analyzers...).
		AddStaticCheckAnalyzers(stylecheck.Analyzers...).
		AddStaticCheckAnalyzers(unused.Analyzer).
		GetChecks(cfg.AnalyzersNames()...)
	multichecker.Main(
		checks...,
	)
}
