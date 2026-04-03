package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CppCompiler compiles C and C++ source files using g++.
type CppCompiler struct{}

// Language returns "cpp".
func (c *CppCompiler) Language() string {
	return "cpp"
}

// Compile compiles a C/C++ source file using g++ with -O2 optimization.
// The output binary is placed in outputDir with the name "solution".
func (c *CppCompiler) Compile(sourcePath string, outputDir string, timeLimit time.Duration, memoryLimitKB int64) CompileResult {
	start := time.Now()

	outputPath := filepath.Join(outputDir, "solution")

	cmd := exec.Command("g++", "-O2", "-o", outputPath, sourcePath, "-lm")
	cmd.Env = append(os.Environ(),
		"PATH=/usr/bin:/bin",
	)

	// Set process limits via the command context.
	var output, stderr []byte
	var err error

	if timeLimit > 0 {
		timer := time.NewTimer(timeLimit)
		done := make(chan error, 1)

		go func() {
			output, err = cmd.CombinedOutput()
			done <- err
		}()

		select {
		case <-done:
			timer.Stop()
		case <-timer.C:
			// Kill the compiler process.
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			<-done
			return CompileResult{
				Success:  false,
				Error:    fmt.Sprintf("compilation timed out after %v", timeLimit),
				TimeUsed: time.Since(start),
			}
		}
	} else {
		output, err = cmd.CombinedOutput()
	}

	_ = stderr // not used separately

	elapsed := time.Since(start)

	if err != nil {
		return CompileResult{
			Success:  false,
			Error:    fmt.Sprintf("compilation error: %s\n%s", err.Error(), string(output)),
			TimeUsed: elapsed,
		}
	}

	// Verify the output file exists.
	if _, statErr := os.Stat(outputPath); statErr != nil {
		return CompileResult{
			Success:  false,
			Error:    fmt.Sprintf("compilation succeeded but output file not found: %v", statErr),
			TimeUsed: elapsed,
		}
	}

	return CompileResult{
		Success:    true,
		OutputPath: outputPath,
		TimeUsed:   elapsed,
	}
}
