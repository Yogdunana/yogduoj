package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Yogdunana/yogduoj/judge/internal/sandbox"
)

// RunResult holds the result of executing a program in the sandbox.
type RunResult struct {
	ExitCode     int
	TimeUsedMs   int64
	MemoryUsedKb int64
	Stdout       string
	Stderr       string
	TimedOut     bool
	OOMKilled    bool
}

// Runner executes compiled programs inside sandboxed containers.
type Runner struct {
	sandboxRoot string
}

// NewRunner creates a new Runner with the given sandbox root directory.
func NewRunner(sandboxRoot string) *Runner {
	return &Runner{sandboxRoot: sandboxRoot}
}

// Run executes a program in a sandboxed environment.
// containerID is used for cgroup tracking.
// execPath is the path to the executable (or for Python, the script).
// language is the programming language ("cpp", "java", "python").
// input is the stdin data to feed to the program.
// timeLimitMs is the maximum execution time in milliseconds.
// memoryLimitKB is the maximum memory in kilobytes.
func (r *Runner) Run(containerID string, execPath string, language string, input string, timeLimitMs int64, memoryLimitKB int64) RunResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitMs)*time.Millisecond+time.Second)
	defer cancel()

	var cmd *exec.Cmd

	switch language {
	case "python", "python3":
		cmd = exec.CommandContext(ctx, "python3", execPath)
	case "java":
		// For Java, execPath is the directory containing .class files.
		// We need to find the main class (assume Main.class).
		mainClass := "Main"
		cmd = exec.CommandContext(ctx, "java", "-Xmx"+fmt.Sprintf("%dK", memoryLimitKB), "-cp", execPath, mainClass)
	default:
		// C/C++ and other compiled languages.
		cmd = exec.CommandContext(ctx, execPath)
	}

	// Set up stdin.
	cmd.Stdin = bytes.NewReader([]byte(input))

	// Capture stdout and stderr.
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Set environment for the sandbox.
	cmd.Env = []string{
		"PATH=/usr/bin:/bin",
		"HOME=/tmp",
		"LANG=C.UTF-8",
		"TERM=dumb",
	}

	// Restrict the working directory.
	workDir := filepath.Join(r.sandboxRoot, containerID, "tmp")
	if err := os.MkdirAll(workDir, 0755); err == nil {
		cmd.Dir = workDir
	}

	// Start the process.
	startTime := time.Now()
	err := cmd.Run()
	elapsed := time.Since(startTime)

	result := RunResult{
		Stdout: stdoutBuf.String(),
		Stderr: stderrBuf.String(),
	}

	// Check if it timed out.
	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.ExitCode = -1
		result.TimeUsedMs = timeLimitMs
		return result
	}

	// Check exit code.
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
			result.Stderr += fmt.Sprintf("\nExecution error: %v", err)
		}
	} else {
		result.ExitCode = 0
	}

	// Try to read memory usage from cgroup.
	memKB, memErr := sandbox.GetMemoryUsageKb(containerID)
	if memErr == nil {
		result.MemoryUsedKb = memKB
	}

	// Check for OOM (exit code 137 = SIGKILL, typically from OOM killer).
	if result.ExitCode == 137 || result.ExitCode == -9 {
		result.OOMKilled = true
	}

	// Calculate time used.
	result.TimeUsedMs = elapsed.Milliseconds()

	// Clamp to the limit.
	if result.TimeUsedMs > timeLimitMs {
		result.TimeUsedMs = timeLimitMs
	}

	return result
}

// RunWithCgroup executes a program with explicit cgroup resource limits.
// This version adds the process to a cgroup for accurate resource tracking.
func (r *Runner) RunWithCgroup(containerID string, execPath string, language string, input string, timeLimitMs int64, memoryLimitKB int64, cpuLimit float64, maxPids int) RunResult {
	// Set up cgroup limits for this execution.
	memoryLimitMB := memoryLimitKB / 1024
	if memoryLimitMB < 1 {
		memoryLimitMB = 1
	}
	if cpuLimit <= 0 {
		cpuLimit = 1.0
	}
	if maxPids <= 0 {
		maxPids = 50
	}

	if err := sandbox.SetCgroupLimits(containerID, cpuLimit, int(memoryLimitMB), maxPids); err != nil {
		// Log warning but continue - cgroup limits are best-effort in this environment.
		fmt.Printf("Warning: failed to set cgroup limits for %s: %v\n", containerID, err)
	}

	// Execute the program.
	result := r.Run(containerID, execPath, language, input, timeLimitMs, memoryLimitKB)

	// Try to get more accurate resource usage from cgroup.
	if memKB, err := sandbox.GetMemoryUsageKb(containerID); err == nil && memKB > result.MemoryUsedKb {
		result.MemoryUsedKb = memKB
	}

	if cpuUs, err := sandbox.GetCpuTimeUs(containerID); err == nil {
		cpuMs := cpuUs / 1000
		if cpuMs > result.TimeUsedMs {
			result.TimeUsedMs = cpuMs
		}
	}

	return result
}
