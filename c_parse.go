package main

import (
	"fmt"
	"log"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
)

func main() {
	cCode := []byte(`
int add(int a, int b) { return a + b; }
void greet() { printf("Hello"); }
static char* getName() { return "Gemini"; }
`)

	parser := tree_sitter.NewParser()
	defer parser.Close()

	lang := tree_sitter.NewLanguage(tree_sitter_c.Language())
	_ = parser.SetLanguage(lang)

	tree := parser.Parse(cCode, nil)
	defer tree.Close()

	// 1. Write a query to find function_definition nodes and capture the declarator (the name/signature)
	// 'declarator' usually contains the function identifier.
	queryString := `
		(function_definition
			declarator: (function_declarator
				declarator: (identifier) @func_name))
	`

	// 2. Compile the query
	query, err := tree_sitter.NewQuery(lang, queryString)
	if err != nil {
		log.Fatalf("Error creating query: %v", err)
	}
	defer query.Close()

	// 3. Execute the query on the root node
	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	captures := cursor.Captures(query, tree.RootNode(), cCode)
	captureNames := query.CaptureNames()

	fmt.Println("Functions found via Query:")

	// 4. Iterate over the matches
	for {
		match, _ := captures.Next()
		if match == nil || len(match.Captures) == 0 {
			break
		}

		for _, capture := range match.Captures {
			// Look up the name by indexing into the slice using capture.Index
			captureName := captureNames[capture.Index]

			if captureName == "func_name" {
				node := capture.Node
				funcName := string(cCode[node.StartByte():node.EndByte()])
				fmt.Printf("- %s (Line: %d)\n", funcName, node.StartPosition().Row+1)
			}
		}
	}
}
