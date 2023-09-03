// Package exitanalyzer определяет Analyzer который
// проверяет, что внутри основной функции нет вызова os.Exit
package exitanalyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

const Doc = `check that there are no calls to os.Exit inside the main function`

var Analyzer = &analysis.Analyzer{
	Name: "exitanalyzer",
	Doc:  Doc,
	Run:  run,
}

// run выполняет анализ исходного кода Go, чтобы найти прямые вызовы функции os.Exit внутри функции main.
// Если такие вызовы обнаруживаются, функция создает диагностическое сообщение.
func run(pass *analysis.Pass) (interface{}, error) {
	// Перебираем все файлы в проекте.
	for _, file := range pass.Files {
		if o := file.Scope.Objects["examples"]; o != nil && file.Name.String() == "main" {
			continue
		}
		// Используем AST для обхода узлов в дереве синтаксического анализа.
		ast.Inspect(file, func(n ast.Node) bool {
			// Проверяем ялвеяется ли узел вызвовом функции
			if callExpr, ok := n.(*ast.CallExpr); ok {
				// Порверяем является ли функция селектором
				if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					// Проверяем ялвяется ли селекторы вызовом функции os.Exit
					if ident, ok := selExpr.X.(*ast.Ident); ok &&
						ident.Name == "os" &&
						selExpr.Sel.Name == "Exit" {
						// Проверяем, находится ли вызов функции os.Exit внутри функции main.
						if isWithinMain(pass, callExpr) {
							// Создаем диагностическое сообщение, если условие выполняется.
							pass.Reportf(callExpr.Pos(), "direct call to os.Exit found within main function")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

// isWithinMain проверяет, находится ли вызов функции os.Exit внутри функции main.
func isWithinMain(pass *analysis.Pass, callExpr *ast.CallExpr) bool {
	// Получаем путь AST от вызова функции os.Exit к верхнему уровню файла.
	path, _ := astutil.PathEnclosingInterval(pass.Files[0], callExpr.Pos(), callExpr.End())
	for _, n := range path {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Проверяем, что это функция main без параметров получателя и анализируется пакет "main".
			if funcDecl.Name.Name == "main" && funcDecl.Recv == nil {
				return true
			}
			break
		}
	}
	return false
}
