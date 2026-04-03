package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// JavaCompiler compiles Java source files using javac.
type JavaCompiler struct{}

// Language returns "java".
func (c *JavaCompiler) Language() string {
	return "java"
}

// Compile compiles a Java source file using javac.
// The compiled .class files are placed in outputDir.
func (c *JavaCompiler) Compile(sourcePath string, outputDir string, timeLimit time.Duration, memoryLimitKB int64) CompileResult {
	start := time.Now()

	// Ensure output directory exists.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return CompileResult{
			Success:  false,
			Error:    fmt.Sprintf("failed to create output directory: %v", err),
			TimeUsed: time.Since(start),
		}
	}

	cmd := exec.Command("javac", "-d", outputDir, sourcePath)
	cmd.Env = append(os.Environ(),
		"PATH=/usr/bin:/bin",
		"JAVA_HOME=/usr/lib/jvm/default-jvm",
	)

	var output []byte
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

	elapsed := time.Since(start)

	if err != nil {
		return CompileResult{
			Success:  false,
			Error:    fmt.Sprintf("compilation error: %s\n%s", err.Error(), string(output)),
			TimeUsed: elapsed,
		}
	}

	// The output path for Java is the directory containing .class files.
	// The runner will need to know the main class name.
	return CompileResult{
		Success:    true,
		OutputPath: outputDir,
		TimeUsed:   elapsed,
	}
}
