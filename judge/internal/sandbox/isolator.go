package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const cgroupBase = "/sys/fs/cgroup"

// SetCgroupLimits creates a cgroup for the given containerID and sets resource limits.
// cpuLimit is in CPU cores (e.g., 1.0 = 1 core, 0.5 = half a core).
// memoryLimitMB is the maximum memory in megabytes.
// maxPids is the maximum number of processes allowed.
func SetCgroupLimits(containerID string, cpuLimit float64, memoryLimitMB int, maxPids int) error {
	cgroupPath := filepath.Join(cgroupBase, containerID)

	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup directory %s: %w", cgroupPath, err)
	}

	// Set CPU quota.
	// cpu.cfs_quota_us = cpuLimit * cpu.cfs_period_us (default period is 100000us)
	quotaUs := int(cpuLimit * 100000)
	if err := os.WriteFile(filepath.Join(cgroupPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(quotaUs)), 0644); err != nil {
		return fmt.Errorf("failed to set cpu quota: %w", err)
	}
	if err := os.WriteFile(filepath.Join(cgroupPath, "cpu.cfs_period_us"), []byte("100000"), 0644); err != nil {
		return fmt.Errorf("failed to set cpu period: %w", err)
	}

	// Set memory limit.
	memoryBytes := int64(memoryLimitMB) * 1024 * 1024
	if err := os.WriteFile(filepath.Join(cgroupPath, "memory.max"), []byte(strconv.FormatInt(memoryBytes, 10)), 0644); err != nil {
		// Fallback for cgroup v1.
		if err := os.WriteFile(filepath.Join(cgroupPath, "memory.limit_in_bytes"), []byte(strconv.FormatInt(memoryBytes, 10)), 0644); err != nil {
			return fmt.Errorf("failed to set memory limit: %w", err)
		}
	}

	// Set PID limit.
	if err := os.WriteFile(filepath.Join(cgroupPath, "pids.max"), []byte(strconv.Itoa(maxPids)), 0644); err != nil {
		return fmt.Errorf("failed to set pid limit: %w", err)
	}

	return nil
}

// AddPidToCgroup adds a process to the cgroup for the given containerID.
func AddPidToCgroup(containerID string, pid int) error {
	cgroupPath := filepath.Join(cgroupBase, containerID)

	// Try cgroup v2 first, then fall back to cgroup v1.
	cgroupProcs := filepath.Join(cgroupPath, "cgroup.procs")
	if _, err := os.Stat(cgroupProcs); os.IsNotExist(err) {
		cgroupProcs = filepath.Join(cgroupPath, "tasks")
	}

	if err := os.WriteFile(cgroupProcs, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to add pid %d to cgroup %s: %w", pid, containerID, err)
	}

	return nil
}

// ResetCgroup resets the cgroup limits for a container to default values.
func ResetCgroup(containerID string) error {
	cgroupPath := filepath.Join(cgroupBase, containerID)

	// Reset CPU quota to unlimited.
	if err := os.WriteFile(filepath.Join(cgroupPath, "cpu.cfs_quota_us"), []byte("-1"), 0644); err != nil {
		// Non-fatal.
	}

	// Reset memory to max.
	if err := os.WriteFile(filepath.Join(cgroupPath, "memory.max"), []byte("max"), 0644); err != nil {
		if err := os.WriteFile(filepath.Join(cgroupPath, "memory.limit_in_bytes"), []byte(strconv.FormatInt(int64(1<<30), 10)), 0644); err != nil {
			// Non-fatal.
		}
	}

	// Reset PID limit.
	if err := os.WriteFile(filepath.Join(cgroupPath, "pids.max"), []byte("max"), 0644); err != nil {
		// Non-fatal.
	}

	return nil
}

// CleanupCgroup removes the cgroup directory for a container.
func CleanupCgroup(containerID string) error {
	cgroupPath := filepath.Join(cgroupBase, containerID)

	// Kill all processes in the cgroup first.
	procsPath := filepath.Join(cgroupPath, "cgroup.procs")
	if _, err := os.Stat(procsPath); os.IsNotExist(err) {
		procsPath = filepath.Join(cgroupPath, "tasks")
	}

	if data, err := os.ReadFile(procsPath); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			pid, err := strconv.Atoi(line)
			if err != nil {
				continue
			}
			// Send SIGKILL to all processes in the cgroup.
			procPath := filepath.Join("/proc", strconv.Itoa(pid))
			if _, err := os.Stat(procPath); err == nil {
				os.WriteFile(filepath.Join(cgroupPath, "cgroup.kill"), []byte("1"), 0644)
			}
		}
	}

	if err := os.RemoveAll(cgroupPath); err != nil {
		return fmt.Errorf("failed to remove cgroup %s: %w", cgroupPath, err)
	}

	return nil
}

// GetMemoryUsageKb reads the current memory usage of a cgroup in kilobytes.
func GetMemoryUsageKb(containerID string) (int64, error) {
	cgroupPath := filepath.Join(cgroupBase, containerID)

	// Try cgroup v2 first.
	data, err := os.ReadFile(filepath.Join(cgroupPath, "memory.current"))
	if err != nil {
		// Fallback to cgroup v1.
		data, err = os.ReadFile(filepath.Join(cgroupPath, "memory.usage_in_bytes"))
		if err != nil {
			return 0, fmt.Errorf("failed to read memory usage: %w", err)
		}
	}

	usage := strings.TrimSpace(string(data))
	bytes, err := strconv.ParseInt(usage, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse memory usage: %w", err)
	}

	return bytes / 1024, nil
}

// GetCpuTimeUs reads the total CPU time used by a cgroup in microseconds.
func GetCpuTimeUs(containerID string) (int64, error) {
	cgroupPath := filepath.Join(cgroupBase, containerID)

	// Try cgroup v2 first.
	data, err := os.ReadFile(filepath.Join(cgroupPath, "cpu.stat"))
	if err != nil {
		return 0, fmt.Errorf("failed to read cpu stat: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "usage_usec ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				val, err := strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					return 0, fmt.Errorf("failed to parse cpu usage: %w", err)
				}
				return val, nil
			}
		}
	}

	return 0, fmt.Errorf("usage_usec not found in cpu.stat")
}
