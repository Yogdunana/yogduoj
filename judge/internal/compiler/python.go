package compiler

import (
	"fmt"
	"os/exec"
	"time"
)

// PythonCompiler handles Python submissions. Python is interpreted, so no
// compilation step is needed. The source file is validated and returned as-is.
type PythonCompiler struct{}

// Language returns "python".
func (c *PythonCompiler) Language() string {
	return "python"
}

// Compile validates that Python3 is available and returns the source path.
func (c *PythonCompiler) Compile(sourcePath string, outputDir string, timeLimit time.Duration, memoryLimitKB int64) CompileResult {
	start := time.Now()

	// Check that python3 is available.
	cmd := exec.Command("python3", "--version")
	if output, err := cmd.CombinedOutput(); err != nil {
		return CompileResult{
			Success:  false,
			Error:    fmt.Sprintf("python3 not available: %s\n%s", err.Error(), string(output)),
			TimeUsed: time.Since(start),
		}
	}

	return CompileResult{
		Success:    true,
		OutputPath: sourcePath,
		TimeUsed:   time.Since(start),
	}
}
