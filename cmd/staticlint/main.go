// Package main запускает multichecker, который выполняет набор анализаторов
// статического кода для проверки качества кода в проекте
// для запуска анализатора необходимо выполнить команду
// make analyze или go run cmd/staticlint/main.go -source ./...
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

// ignoreChecks содержит анализаторы, которые следует игнорировать во
// время выполнения multichecker.
var ignoreChecks = map[string]bool{
	"ST1000": true,
}

type checks struct {
	Analyzers []*analysis.Analyzer
}

func main() {
	multichecker.Main(getAnalyzers()...)
}

// getAnalyzers возвращает набор анализаторов для проверки.
func getAnalyzers() []*analysis.Analyzer {
	c := checks{}

	c.addAnalyzers(staticcheck.Analyzers)
	c.addAnalyzers(stylecheck.Analyzers)

	c.setStandardAnalyzers()

	return c.Analyzers
}

// addAnalyzers добавляет анализаторы из заданного списка в набор анализаторов,
// если они не были помечены как игнорируемые.
func (c *checks) addAnalyzers(analyzers []*lint.Analyzer) {
	for _, v := range analyzers {
		if !ignoreChecks[v.Analyzer.Name] {
			c.Analyzers = append(c.Analyzers, v.Analyzer)
		}
	}
}

// setStandardAnalyzers добавляет стандартные анализаторы, такие как анализаторы
// из стандартной библиотеки Go и другие общепринятые анализаторы, в набор.
//
// exitanalyzer.Analyzer проверяет, что внутри основной функции нет вызова os.Exit
//
// printf.Analyzer: Проверяет аргументы и форматирование строк в функциях пакета fmt. Обнаруживает потенциальные ошибки форматирования.
//
// shadow.Analyzer: Выявляет случаи, когда переменные с одинаковыми именами объявлены во вложенных областях видимости.
//
// structtag.Analyzer: Проверяет использование тегов в структурах, выявляет неиспользуемые и неправильно отформатированные теги.
//
// asmdecl.Analyzer: Проверяет использование инструкций сборки в коде, обнаруживает ошибки в инструкциях сборки.
//
// assign.Analyzer: Выявляет случаи, когда переменные объявляются, но не используются.
//
// atomic.Analyzer: Проверяет операции с атомарными переменными, обнаруживает ошибки и проблемы многопоточности.
//
// bools.Analyzer: Исследует операции с логическими значениями и выявляет неправильное использование логических операторов.
//
// atomicalign.Analyzer: Проверяет выравнивание данных в структурах, используемых с атомарными переменными.
//
// composite.Analyzer: Проверяет использование составных литералов в коде и находит лишние литералы.
//
// buildtag.Analyzer: Проверяет использование тегов сборки в коде, находит неиспользуемые теги сборки.т улучшить читаемость и эффективность кода.
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
