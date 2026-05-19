package main

import (
	"fmt"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
)

// FunctionInfo stores the name and structural boundaries of a parsed function.
type FunctionInfo struct {
	Name      string
	StartLine uint32
	StartCol  uint32
	EndLine   uint32
	EndCol    uint32
}

// CParser holds the reusable tree-sitter instances to optimize multi-file parsing.
type CParser struct {
	parser       *tree_sitter.Parser
	query        *tree_sitter.Query
	captureNames []string
}

// NewCParser initializes the parser and compiles the query exactly once.
func NewCParser() (*CParser, error) {
	parser := tree_sitter.NewParser()

	lang := tree_sitter.NewLanguage(tree_sitter_c.Language())
	if err := parser.SetLanguage(lang); err != nil {
		parser.Close()
		return nil, fmt.Errorf("failed to set C language: %w", err)
	}

	// S-Expression query string targeting the body wrapper and internal identifier names
	queryString := `
		(function_definition) @func_body
		(function_definition
			declarator: (function_declarator
				declarator: (identifier) @func_name))
	`
	query, err := tree_sitter.NewQuery(lang, queryString)
	if err != nil {
		parser.Close()
		return nil, fmt.Errorf("invalid tree-sitter query: %w", err)
	}

	return &CParser{
		parser:       parser,
		query:        query,
		captureNames: query.CaptureNames(),
	}, nil
}

// Close frees the underlying C memory allocated by Tree-sitter.
func (cp *CParser) Close() {
	if cp.query != nil {
		cp.query.Close()
	}
	if cp.parser != nil {
		cp.parser.Close()
	}
}

// ExtractFunctions returns a detailed slice of metadata for all functions found.
func (cp *CParser) ExtractFunctions(cCode []byte) ([]FunctionInfo, error) {
	tree := cp.parser.Parse(cCode, nil)
	if tree == nil {
		return nil, fmt.Errorf("parser returned an empty syntax tree")
	}
	defer tree.Close()

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	captures := cursor.Captures(cp.query, tree.RootNode(), cCode)
	var functions []FunctionInfo

	// Temporary structure to align a body match with its subsequent name match
	var currentFunc FunctionInfo

	for {
		match, _ := captures.Next()
		if match == nil || len(match.Captures) == 0 {
			break
		}

		for _, capture := range match.Captures {
			captureType := cp.captureNames[capture.Index]
			node := capture.Node

			if captureType == "func_body" {
				// Typecast uint to uint32 to satisfy Go's strict assignment rules
				currentFunc.StartLine = uint32(node.StartPosition().Row + 1)
				currentFunc.StartCol = uint32(node.StartPosition().Column)
				currentFunc.EndLine = uint32(node.EndPosition().Row + 1)
				currentFunc.EndCol = uint32(node.EndPosition().Column)

			} else if captureType == "func_name" {
				currentFunc.Name = string(cCode[node.StartByte():node.EndByte()])

				// Commit to results slice once both body limits and name are grouped
				functions = append(functions, currentFunc)
				currentFunc = FunctionInfo{}
			}
		}
	}

	return functions, nil
}
