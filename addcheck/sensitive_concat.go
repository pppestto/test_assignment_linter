package addcheck

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/pppestto/test_assignment_linter/addcheck/rules"
	"golang.org/x/tools/go/analysis"
)

func stringLitUnquoted(e ast.Expr) (s string, ok bool) {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return s, true
}

func reportSensitiveConcat(pass *analysis.Pass, arg ast.Expr, keywords []string) {
	var visit func(ast.Expr)
	visit = func(e ast.Expr) {
		be, ok := e.(*ast.BinaryExpr)
		if !ok || be.Op != token.ADD {
			return
		}
		visit(be.X)
		visit(be.Y)

		xStr, xLit := stringLitUnquoted(be.X)
		yStr, yLit := stringLitUnquoted(be.Y)
		if xLit && !yLit && rules.ContainsSensitiveSubstring(xStr, keywords) {
			pass.Reportf(be.Y.Pos(), "[sensitive] do not concatenate sensitive prefix with variable; use static message like \"token validated\"")
		}
		if yLit && !xLit && rules.ContainsSensitiveSubstring(yStr, keywords) {
			pass.Reportf(be.X.Pos(), "[sensitive] do not concatenate variable with sensitive suffix; use static message")
		}
	}
	visit(arg)
}
