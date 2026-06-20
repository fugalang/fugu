package builder

import (
	"os"
	"runtime"
)

type TypeOs uint

const (
	_ TypeOs = iota
	Darwin
	GnuLinux

	Windows10p
	Windows7
	WindowsXP

	Freebsd
	Openbsd
	Netbsd

	KernelMod
)

func CurrentOs() TypeOs {
	switch runtime.GOOS {
	case "darwin":
		return Darwin
	case "windows":
		return Windows10p
	case "linux":
		return GnuLinux
	case "freebsd":
		return Freebsd
	case "openbsd":
		return Openbsd
	case "netbsd":
		return Netbsd
	default:
		return KernelMod // хз почему
	}
}

func GetTargetOsEnv() TypeOs {
	tgos := os.Getenv("TARGET-OS")
	switch tgos {
	case "DARWIN", "MACOS":
		return Darwin
	case "WINDOWS", "WINDOWS10", "WINDOWS11":
		return Windows10p
	case "GNU-LINUX", "LINUX":
		return GnuLinux
	case "WINDOWS7":
		return Windows7
	case "WINDOWS-XP":
		return WindowsXP
	case "FREEBSD":
		return Freebsd
	case "OPENBSD":
		return Openbsd
	case "NETBSD":
		return Netbsd
	case "KERNEL-MOD", "NO-STD", "LOW-MOD":
		return KernelMod
	default:
		return CurrentOs()
	}
}
