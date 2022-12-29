package main

import (
	"golang.org/x/tools/go/analysis"
	"goroutine-with-recover/analyzer"
)

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}

var AnalyzerPlugin analyzerPlugin
