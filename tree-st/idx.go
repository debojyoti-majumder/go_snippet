package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// 1. Initialize the shared parser wrapper
	cParser, err := NewCParser()
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	defer cParser.Close()

	fileNames := []string{
		"/Users/debojyotim/Documents/module/f1.c",
		"/Users/debojyotim/Documents/module/f2.c",
	}

	for _, filename := range fileNames {
		fileBuffer, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal("Unable to read from the file")
		}

		// 3. Extract metadata
		funcs, err := cParser.ExtractFunctions(fileBuffer)
		if err != nil {
			log.Fatalf("Error parsing: %v", err)
		}

		// 4. Output results
		fmt.Printf("Extracted Function Boundaries in file %s\n", filename)
		fmt.Println("Function count:", len(funcs))

		for _, f := range funcs {
			fmt.Printf("     Function '%s()':\n", f.Name)
			fmt.Printf("        Starts at -> Line: %d, Column: %d\n", f.StartLine, f.StartCol)
			fmt.Printf("        Ends at   -> Line: %d, Column: %d\n\n", f.EndLine, f.EndCol)
		}
	}
}
