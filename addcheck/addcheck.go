package addcheck

import (
	"go/ast"
	"regexp"
	"strconv"

	"github.com/pppestto/test_assignment_linter/addcheck/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "checks log messages for slog; sensitive = concat with variable; -config; -fix",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var configPath string

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", "", "path to loglint JSON config (optional)")
}

func run(pass *analysis.Pass) (interface{}, error) {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	var patterns []*regexp.Regexp
	for _, p := range cfg.SensitivePatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, re)
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		if !IsSlogMessageCall(pass, call) {
			return
		}

		if cfg.ruleEnabled("sensitive") && len(call.Args) > 0 {
			arg0 := call.Args[0]
			if _, _, ok := LogLiteral(call); !ok {
				reportSensitiveConcat(pass, arg0, cfg.SensitiveKeywords)
			}
			if lit, msg, ok := LogLiteral(call); ok {
				if v := rules.CheckSensitivePatterns(msg, patterns); v != nil {
					pass.Reportf(lit.Pos(), "[%s] %s", v.Rule, v.Message)
				}
			}
		}

		lit, msg, ok := LogLiteral(call)
		if !ok {
			return
		}

		report := func(rule, message string) {
			pass.Reportf(lit.Pos(), "[%s] %s", rule, message)
		}

		reportWithFix := func(rule, message string, fixed string) {
			if fixed == msg {
				report(rule, message)
				return
			}
			newText := []byte(strconv.Quote(fixed))
			pass.Report(analysis.Diagnostic{
				Pos:            lit.Pos(),
				Message:        "[" + rule + "] " + message,
				SuggestedFixes: []analysis.SuggestedFix{{Message: "apply fix", TextEdits: []analysis.TextEdit{{Pos: lit.Pos(), End: lit.End(), NewText: newText}}}},
			})
		}

		if cfg.ruleEnabled("lowercase") {
			if v := rules.CheckLowerCase(msg); v != nil {
				reportWithFix(v.Rule, v.Message, rules.FixLowercaseFirst(msg))
			}
		}
		if cfg.ruleEnabled("english") {
			if v := rules.CheckEnglish(msg); v != nil {
				report(v.Rule, v.Message)
			}
		}
		if cfg.ruleEnabled("emoji") {
			if v := rules.CheckEmoji(msg); v != nil {
				reportWithFix(v.Rule, v.Message, rules.FixEmojiTZ(msg))
			}
		}
	})
	return nil, nil
}
