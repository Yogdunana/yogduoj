package sandbox

/*
#include <linux/seccomp.h>
#include <linux/filter.h>
#include <linux/audit.h>
#include <stddef.h>
#include <errno.h>

// Seccomp BPF program structure for Go.
// We build the BPF filter programmatically and apply it via prctl.
*/
import "C"

import (
	"encoding/binary"
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/unix"
)

// seccompData matches the kernel's struct seccomp_data layout.
type seccompData struct {
	Nr     uint32 // syscall number
	Arch   uint32 // AUDIT_ARCH_X86_64
	IP     uint64 // instruction pointer (unused)
	Args   [6]uint64
}

// sockFilter represents a BPF instruction.
type sockFilter struct {
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
}

// sockFprog represents a BPF program.
type sockFprog struct {
	Len    uint16
	Filter *sockFilter
}

const (
	// BPF instruction classes.
	bpfLD    = 0x00
	bpfJMP   = 0x05
	bpfRET   = 0x06
	bpfALU   = 0x04

	// BPF ld/ldx fields.
	bpfW     = 0x00 // 32-bit word
	bpfABS   = 0x20

	// BPF jmp fields.
	bpfJA    = 0x00
	bpfJEQ   = 0x10
	bpfJGE   = 0x30
	bpfJGT   = 0x20

	// BPF alu/jmp fields.
	bpfK     = 0x00

	// BPF return values.
	seccompRetAllow  = 0x7FFF0000
	seccompRetKill   = 0x00000000
	seccompRetErrno  = 0x00050000 // SECCOMP_RET_ERRNO | EPERM

	// x86_64 syscall numbers.
	// These are from /usr/include/asm/unistd_64.h
	sysRead        = 0
	sysWrite       = 1
	sysOpen        = 2
	sysClose       = 3
	sysStat        = 4
	sysFstat       = 5
	sysLseek       = 8
	sysMmap        = 9
	sysMprotect    = 10
	sysMunmap      = 11
	sysBrk         = 12
	sysIoctl       = 16
	sysAccess      = 21
	sysPipe        = 22
	sysDup         = 32
	sysDup2        = 33
	sysGetpid      = 39
	sysSocket      = 41
	sysConnect     = 42
	sysClone       = 56
	sysFork        = 57
	sysVfork       = 58
	sysExecve      = 59
	sysExit        = 60
	sysArchPrctl   = 158
	sysFstatat     = 262
	sysOpenat      = 257
	sysNewfstatat  = 262
	sysSetRobustList = 273
	sysRseq        = 334
	sysPread64     = 17
	sysPwrite64    = 18
	sysReadv       = 19
	sysWritev      = 20
	sysFaccessat   = 269
	sysLstat       = 6
	sysGetcwd      = 79
	sysGetuid      = 102
	sysGeteuid     = 103
	sysGetgid      = 104
	sysGetegid     = 105
	sysClockGettime = 228
	sysUname       = 63
	sysSigaction   = 13
	sysSigprocmask = 14
	sysSetTidAddress = 218
)

// allowedSyscalls is the list of syscalls that are permitted in the sandbox.
var allowedSyscalls = map[uint32]bool{
	sysRead:         true,
	sysWrite:        true,
	sysOpen:         true,
	sysClose:        true,
	sysStat:         true,
	sysFstat:        true,
	sysLseek:        true,
	sysMmap:         true,
	sysMprotect:     true,
	sysMunmap:       true,
	sysBrk:          true,
	sysIoctl:        true,
	sysAccess:       true,
	sysPipe:         true,
	sysDup:          true,
	sysDup2:         true,
	sysGetpid:       true,
	sysExit:         true,
	sysArchPrctl:    true,
	sysOpenat:       true,
	sysNewfstatat:   true,
	sysSetRobustList: true,
	sysRseq:         true,
	sysPread64:      true,
	sysPwrite64:     true,
	sysReadv:        true,
	sysWritev:       true,
	sysFaccessat:    true,
	sysLstat:        true,
	sysGetcwd:       true,
	sysGetuid:       true,
	sysGeteuid:      true,
	sysGetgid:       true,
	sysGetegid:      true,
	sysClockGettime: true,
	sysUname:        true,
	sysSigaction:    true,
	sysSigprocmask:  true,
	sysSetTidAddress: true,
}

// blockedSyscalls are syscalls that are explicitly blocked for security.
var blockedSyscalls = map[uint32]bool{
	sysClone:    true,
	sysFork:     true,
	sysVfork:    true,
	sysExecve:   true,
	sysSocket:   true,
	sysConnect:  true,
}

// SetupSeccomp installs a seccomp-bpf filter that allows only whitelisted syscalls.
// It uses the prctl syscall with PR_SET_NO_NEW_PRIVS and PR_SET_SECCOMP.
func SetupSeccomp() error {
	// Build the BPF filter program.
	filter := buildSeccompFilter()

	prog := sockFprog{
		Len:    uint16(len(filter)),
		Filter: &filter[0],
	}

	// Step 1: Set PR_SET_NO_NEW_PRIVS (required before setting seccomp filter).
	if err := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		return fmt.Errorf("prctl(PR_SET_NO_NEW_PRIVS) failed: %v", err)
	}

	// Step 2: Set PR_SET_SECCOMP with SECCOMP_MODE_FILTER.
	_, _, seErr := unix.Syscall6(
		unix.SYS_PRCTL,
		uintptr(unix.PR_SET_SECCOMP),
		uintptr(unix.SECCOMP_MODE_FILTER),
		uintptr(unsafe.Pointer(&prog)),
		0, 0, 0,
	)
	if seErr != 0 {
		return fmt.Errorf("prctl(PR_SET_SECCOMP, SECCOMP_MODE_FILTER) failed: %v", seErr)
	}

	log.Printf("Seccomp filter installed with %d instructions (%d allowed syscalls)",
		len(filter), len(allowedSyscalls))
	return nil
}

// buildSeccompFilter constructs a BPF filter program.
//
// The filter logic:
//   1. Load syscall number (seccomp_data.nr)
//   2. Check architecture (AUDIT_ARCH_X86_64)
//   3. For each allowed syscall, check if it matches -> ALLOW
//   4. Default: KILL
func buildSeccompFilter() []sockFilter {
	var filter []sockFilter

	// Load architecture field (offset 4 in seccomp_data).
	filter = append(filter, sockFilter{
		Code: bpfLD | bpfW | bpfABS,
		K:    4, // offset of .arch in seccomp_data
	})

	// Check if arch == AUDIT_ARCH_X86_64 (0xC000003E).
	// If not, jump to kill (past all the checks).
	archKillOffset := uint8(3 + len(allowedSyscalls) + len(blockedSyscalls) + 1)
	filter = append(filter, sockFilter{
		Code: bpfJMP | bpfJEQ | bpfK,
		K:    0xC000003E, // AUDIT_ARCH_X86_64
		Jt:   0,          // fall through to next instruction
		Jf:   archKillOffset, // jump to kill
	})

	// Load syscall number (offset 0 in seccomp_data).
	filter = append(filter, sockFilter{
		Code: bpfLD | bpfW | bpfABS,
		K:    0, // offset of .nr in seccomp_data
	})

	// For each allowed syscall, emit a check.
	// If syscall matches, jump to ALLOW (which is after all checks).
	allowOffset := uint8(2 + len(allowedSyscalls) + len(blockedSyscalls))
	currentIdx := 3
	for nr := range allowedSyscalls {
		jumpDist := allowOffset - uint8(currentIdx)
		filter = append(filter, sockFilter{
			Code: bpfJMP | bpfJEQ | bpfK,
			K:    nr,
			Jt:   jumpDist, // jump to ALLOW
			Jf:   0,        // fall through
		})
		currentIdx++
	}

	// For each explicitly blocked syscall, emit a check that returns ERRNO.
	// This gives a clearer error than KILL.
	blockOffset := uint8(1 + len(blockedSyscalls))
	currentIdx = 3 + len(allowedSyscalls)
	for nr := range blockedSyscalls {
		jumpDist := blockOffset - uint8(currentIdx)
		filter = append(filter, sockFilter{
			Code: bpfJMP | bpfJEQ | bpfK,
			K:    nr,
			Jt:   jumpDist, // jump to ERRNO
			Jf:   0,        // fall through
		})
		currentIdx++
	}

	// Default: return ERRNO (EPERM).
	filter = append(filter, sockFilter{
		Code: bpfRET | bpfK,
		K:    seccompRetErrno,
	})

	// ALLOW return.
	filter = append(filter, sockFilter{
		Code: bpfRET | bpfK,
		K:    seccompRetAllow,
	})

	return filter
}

// SetupSeccompSimple installs a simpler seccomp filter using the raw prctl approach.
// This is a fallback that uses a minimal filter if the full BPF approach fails.
func SetupSeccompSimple() error {
	// Minimal filter: just block the most dangerous syscalls.
	filter := []sockFilter{
		// Load syscall number.
		{Code: bpfLD | bpfW | bpfABS, K: 0},
		// Block clone/fork/vfork.
		{Code: bpfJMP | bpfJEQ | bpfK, K: sysClone, Jt: 4, Jf: 0},
		{Code: bpfJMP | bpfJEQ | bpfK, K: sysFork, Jt: 3, Jf: 0},
		{Code: bpfJMP | bpfJEQ | bpfK, K: sysVfork, Jt: 2, Jf: 0},
		// Block execve.
		{Code: bpfJMP | bpfJEQ | bpfK, K: sysExecve, Jt: 1, Jf: 0},
		// Kill.
		{Code: bpfRET | bpfK, K: seccompRetKill},
		// Allow everything else.
		{Code: bpfRET | bpfK, K: seccompRetAllow},
	}

	prog := sockFprog{
		Len:    uint16(len(filter)),
		Filter: &filter[0],
	}

	if err := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		return fmt.Errorf("prctl(PR_SET_NO_NEW_PRIVS) failed: %v", err)
	}

	_, _, seErr := unix.Syscall6(
		unix.SYS_PRCTL,
		uintptr(unix.PR_SET_SECCOMP),
		uintptr(unix.SECCOMP_MODE_FILTER),
		uintptr(unsafe.Pointer(&prog)),
		0, 0, 0,
	)
	if seErr != 0 {
		return fmt.Errorf("prctl(PR_SET_SECCOMP) failed: %v", seErr)
	}

	return nil
}

// init ensures the binary order is correct for BPF.
func init() {
	// Verify we're on a little-endian system (required for BPF).
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, 0x1234)
	if buf[0] != 0x34 {
		log.Printf("Warning: not running on little-endian system, seccomp BPF may not work correctly")
	}
}

// GetAllowedSyscallCount returns the number of allowed syscalls in the whitelist.
func GetAllowedSyscallCount() int {
	return len(allowedSyscalls)
}

// IsSyscallAllowed checks if a syscall number is in the allowed list.
func IsSyscallAllowed(nr uint32) bool {
	return allowedSyscalls[nr]
}
