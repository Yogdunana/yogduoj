package compiler

import "time"

// CompileResult holds the result of a compilation attempt.
type CompileResult struct {
	Success   bool
	OutputPath string
	Error     string
	TimeUsed  time.Duration
}

// Compiler defines the interface for language-specific compilers.
type Compiler interface {
	// Compile compiles the source file and returns the result.
	// sourcePath is the path to the source file.
	// outputDir is the directory where compiled output should be placed.
	// timeLimit is the maximum allowed compilation time.
	// memoryLimitKB is the maximum allowed memory during compilation.
	Compile(sourcePath string, outputDir string, timeLimit time.Duration, memoryLimitKB int64) CompileResult

	// Language returns the language identifier (e.g., "cpp", "java", "python").
	Language() string
}

// GetCompiler returns the appropriate compiler for the given language.
func GetCompiler(language string) Compiler {
	switch language {
	case "cpp", "c++", "cc", "cxx", "c":
		return &CppCompiler{}
	case "java":
		return &JavaCompiler{}
	case "python", "python3", "py":
		return &PythonCompiler{}
	default:
		return nil
	}
}
