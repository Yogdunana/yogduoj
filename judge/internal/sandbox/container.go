package sandbox

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sys/unix"
)

// Container represents an isolated execution environment.
type Container struct {
	ID        string
	PID       int
	busy      int32 // atomic: 0 = free, 1 = busy
	createdAt time.Time
	rootPath  string
	pipeR     *os.File
	pipeW     *os.File
}

// IsBusy returns whether the container is currently in use.
func (c *Container) IsBusy() bool {
	return atomic.LoadInt32(&c.busy) == 1
}

// markBusy atomically sets the container as busy.
func (c *Container) markBusy() {
	atomic.StoreInt32(&c.busy, 1)
}

// markFree atomically sets the container as free.
func (c *Container) markFree() {
	atomic.StoreInt32(&c.busy, 0)
}

// ContainerPool manages a pool of pre-created containers.
type ContainerPool struct {
	containers []*Container
	freeCh     chan *Container
	rootDir    string
	size       int
	mu         sync.Mutex
	closed     bool
}

// NewContainerPool creates and initializes a pool of N containers.
func NewContainerPool(rootDir string, size int) (*ContainerPool, error) {
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create sandbox root %s: %w", rootDir, err)
	}

	pool := &ContainerPool{
		containers: make([]*Container, 0, size),
		freeCh:     make(chan *Container, size),
		rootDir:    rootDir,
		size:       size,
	}

	for i := 0; i < size; i++ {
		c, err := pool.createContainer(i)
		if err != nil {
			// Cleanup already created containers on failure.
			pool.Close()
			return nil, fmt.Errorf("failed to create container %d: %w", i, err)
		}
		pool.containers = append(pool.containers, c)
		pool.freeCh <- c
	}

	log.Printf("Container pool initialized with %d containers at %s", size, rootDir)
	return pool, nil
}

// createContainer creates a single container with isolated namespaces.
func (p *ContainerPool) createContainer(index int) (*Container, error) {
	id := fmt.Sprintf("judge-%d", index)
	rootPath := filepath.Join(p.rootDir, id)

	// Create container root filesystem.
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create container root: %w", err)
	}

	// Create standard directories inside the container.
	for _, dir := range []string{"bin", "lib", "lib64", "tmp", "proc", "dev", "etc", "usr"} {
		if err := os.MkdirAll(filepath.Join(rootPath, dir), 0755); err != nil {
			return nil, fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	// Create pipes for parent-child communication.
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe: %w", err)
	}

	// Create the child process in new namespaces.
	// CLONE_NEWPID  - new process ID namespace
	// CLONE_NEWNS   - new mount namespace
	// CLONE_NEWNET  - new network namespace
	// CLONE_NEWUTS  - new UTS (hostname) namespace
	// CLONE_NEWIPC  - new IPC namespace
	cloneFlags := uintptr(unix.CLONE_NEWPID | unix.CLONE_NEWNS | unix.CLONE_NEWNET | unix.CLONE_NEWUTS | unix.CLONE_NEWIPC | unix.SIGCHLD)

	// We use a simple approach: fork a child that waits for commands via pipe.
	// The child process runs in isolated namespaces.
	pid, _, errno := unix.Syscall6(
		unix.SYS_CLONE,
		cloneFlags,
		0, // child stack (0 = copy parent stack)
		0, // parent tid
		0, // tls
		0, // child tid
		0,
	)
	if errno != 0 {
		pipeR.Close()
		pipeW.Close()
		return nil, fmt.Errorf("clone failed: %v", errno)
	}

	if pid == 0 {
		// Child process - enter the sandbox.
		pipeR.Close()
		childSandbox(pipeW, rootPath)
		// Should never reach here.
		os.Exit(0)
	}

	// Parent process.
	pipeW.Close()

	container := &Container{
		ID:        id,
		PID:       int(pid),
		busy:      0,
		createdAt: time.Now(),
		rootPath:  rootPath,
		pipeR:     pipeR,
		pipeW:     pipeW,
	}

	// Create cgroup for this container.
	if err := SetCgroupLimits(id, 1, 256, 100); err != nil {
		log.Printf("Warning: failed to create cgroup for %s: %v", id, err)
	}

	return container, nil
}

// childSandbox is the entry point for child processes running in isolated namespaces.
// It waits for commands from the parent via the pipe and executes them.
func childSandbox(pipeW *os.File, rootPath string) {
	// Mount proc filesystem.
	if err := unix.Mount("proc", filepath.Join(rootPath, "proc"), "proc", 0, ""); err != nil {
		// Non-fatal for basic operation.
	}

	// Set hostname.
	if err := unix.Sethostname([]byte("sandbox")); err != nil {
		// Non-fatal.
	}

	// Signal ready to parent.
	pipeW.Write([]byte("ready"))
	pipeW.Close()

	// Wait to be killed. The parent will send signals to control execution.
	// In a production system, this would be a command loop reading from a pipe/socket.
	select {}
}

// Get returns a free container from the pool, blocking if necessary.
func (p *ContainerPool) Get() (*Container, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, fmt.Errorf("container pool is closed")
	}
	p.mu.Unlock()

	select {
	case c := <-p.freeCh:
		c.markBusy()
		return c, nil
	default:
		return nil, fmt.Errorf("no free containers available")
	}
}

// Put returns a container to the pool after use, resetting it for reuse.
func (p *ContainerPool) Put(c *Container) {
	if err := p.resetContainer(c); err != nil {
		log.Printf("Warning: failed to reset container %s: %v", c.ID, err)
	}
	c.markFree()
	p.freeCh <- c
}

// resetContainer cleans up the container filesystem for reuse.
func (p *ContainerPool) resetContainer(c *Container) error {
	// Clean the tmp directory.
	tmpDir := filepath.Join(c.rootPath, "tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("failed to clean tmp: %w", err)
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate tmp: %w", err)
	}

	// Reset cgroup limits.
	if err := ResetCgroup(c.ID); err != nil {
		log.Printf("Warning: failed to reset cgroup for %s: %v", c.ID, err)
	}

	return nil
}

// Size returns the total number of containers in the pool.
func (p *ContainerPool) Size() int {
	return p.size
}

// FreeCount returns the number of currently available containers.
func (p *ContainerPool) FreeCount() int {
	return len(p.freeCh)
}

// BusyCount returns the number of currently busy containers.
func (p *ContainerPool) BusyCount() int {
	return p.size - len(p.freeCh)
}

// Close shuts down all containers in the pool and cleans up resources.
func (p *ContainerPool) Close() {
	p.mu.Lock()
	p.closed = true
	p.mu.Unlock()

	for _, c := range p.containers {
		// Send SIGKILL to the child process.
		if c.PID > 0 {
			unix.Kill(c.PID, unix.SIGKILL)
			// Wait for the child to avoid zombies.
			unix.Wait4(c.PID, nil, 0, nil)
		}
		if c.pipeR != nil {
			c.pipeR.Close()
		}
		// Cleanup cgroup.
		CleanupCgroup(c.ID)
		// Cleanup filesystem.
		os.RemoveAll(c.rootPath)
	}

	log.Printf("Container pool shut down, %d containers cleaned up", len(p.containers))
}
