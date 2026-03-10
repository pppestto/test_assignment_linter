package addcheck

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

const (
	pkgLogSlog = "log/slog"
	pkgUberZap = "go.uber.org/zap"
)

var slogMsgFuncs = map[string]bool{
	"Info": true, "Warn": true, "Error": true, "Debug": true, "Log": true,
	"InfoContext": true, "WarnContext": true, "ErrorContext": true, "DebugContext": true, "LogContext": true,
}

func LogLiteral(call *ast.CallExpr) (lit *ast.BasicLit, msg string, ok bool) {
	if len(call.Args) == 0 {
		return nil, "", false
	}
	l, isLit := call.Args[0].(*ast.BasicLit)
	if !isLit || l.Kind != token.STRING {
		return nil, "", false
	}
	s, err := strconv.Unquote(l.Value)
	if err != nil {
		return nil, "", false
	}
	return l, s, true
}

func IsSlogMessageCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	if len(call.Args) == 0 {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if !slogMsgFuncs[sel.Sel.Name] {
		return false
	}
	if selObj, exists := pass.TypesInfo.Selections[sel]; exists {
		if fn, isFunc := selObj.Obj().(*types.Func); isFunc {
			if pkg := fn.Pkg(); pkg != nil && pkg.Path() == pkgLogSlog {
				return true
			}
		}
	}
	if obj := pass.TypesInfo.Uses[sel.Sel]; obj != nil {
		if fn, isFunc := obj.(*types.Func); isFunc {
			if pkg := fn.Pkg(); pkg != nil && pkg.Path() == pkgLogSlog {
				return true
			}
		}
	}
	// go.uber.org/zap: logger.Info("msg", ...) — тот же набор имён методов
	if selObj, exists := pass.TypesInfo.Selections[sel]; exists {
		if fn, isFunc := selObj.Obj().(*types.Func); isFunc {
			if pkg := fn.Pkg(); pkg != nil && pkg.Path() == pkgUberZap {
				if sel.Sel.Name == "Info" || sel.Sel.Name == "Error" || sel.Sel.Name == "Debug" || sel.Sel.Name == "Warn" || sel.Sel.Name == "DPanic" {
					return true
				}
			}
		}
	}
	return false
}

func TryExtractLogMessage(pass *analysis.Pass, call *ast.CallExpr) (msg string, pos token.Pos, ok bool) {
	if !IsSlogMessageCall(pass, call) {
		return "", 0, false
	}
	l, s, ok := LogLiteral(call)
	if !ok {
		return "", 0, false
	}
	return s, l.Pos(), true
}
