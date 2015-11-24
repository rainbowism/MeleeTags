package main

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	// TH32CSSnapProcess https://msdn.microsoft.com/en-us/library/windows/desktop/ms682489(v=vs.85).aspx
	TH32CSSnapProcess = 0x00000002
	// TH32CSSnapModule https://msdn.microsoft.com/en-us/library/windows/desktop/ms682489(v=vs.85).aspx
	TH32CSSnapModule = 0x00000008

	//ProcessVMRead https://msdn.microsoft.com/en-ca/library/windows/desktop/ms684880(v=vs.85).aspx
	ProcessVMRead = 0x0010

	// MaxModuleName32 ??
	MaxModuleName32 = 255
)

type (
	// HANDLE https://msdn.microsoft.com/en-us/library/windows/desktop/aa383751(v=vs.85).aspx
	HANDLE uintptr
	// HMODULE https://msdn.microsoft.com/en-us/library/windows/desktop/aa383751(v=vs.85).aspx
	HMODULE HANDLE

	// PROCESSENTRY32 https://msdn.microsoft.com/en-us/library/windows/desktop/ms684839(v=vs.85).aspx
	PROCESSENTRY32 struct {
		Size            uint32
		Usage           uint32
		ProcessID       uint32
		DefaultHeapID   int
		ModuleID        uint32
		Threads         uint32
		ParentProcessID uint32
		PriClassBase    int32
		Flags           uint32
		ExeFile         [syscall.MAX_PATH]uint16
	}

	// MODULEENTRY32 https://msdn.microsoft.com/en-us/library/windows/desktop/ms684225(v=vs.85).aspx
	MODULEENTRY32 struct {
		Size        uint32
		ModuleID    uint32
		ProcessID   uint32
		GlobalUsage uint32
		ProcUsage   uint32
		BaseAddr    *uint8
		BaseSize    uint32
		HModule     HMODULE
		Module      [MaxModuleName32 + 1]uint16
		ExePath     [syscall.MAX_PATH]uint16
	}
)

var (
	kernel32 = syscall.MustLoadDLL("kernel32.dll")

	procCloseHandle              = kernel32.MustFindProc("CloseHandle")
	procCreateToolhelp32Snapshot = kernel32.MustFindProc("CreateToolhelp32Snapshot")
	procProcess32First           = kernel32.MustFindProc("Process32FirstW")
	procProcess32Next            = kernel32.MustFindProc("Process32NextW")
	procModule32First            = kernel32.MustFindProc("Module32FirstW")
	procModule32Next             = kernel32.MustFindProc("Module32NextW")
	procOpenProcess              = kernel32.MustFindProc("OpenProcess")
	procReadProcessMemory        = kernel32.MustFindProc("ReadProcessMemory")
)

// CloseHandle https://msdn.microsoft.com/en-us/library/windows/desktop/ms724211(v=vs.85).aspx
func CloseHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(
		uintptr(object),
	)
	return ret != 0
}

// CreateToolhelp32Snapshot https://msdn.microsoft.com/en-us/library/windows/desktop/ms682489(v=vs.85).aspx
func CreateToolhelp32Snapshot(flags, processID uint32) HANDLE {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processID),
	)

	if ret <= 0 {
		return HANDLE(0)
	}

	return HANDLE(ret)
}

// Process32First https://msdn.microsoft.com/en-us/library/windows/desktop/ms684834(v=vs.85).aspx
func Process32First(snapshot HANDLE, pe *PROCESSENTRY32) bool {
	ret, _, _ := procProcess32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(pe)),
	)
	return ret != 0
}

// Process32Next https://msdn.microsoft.com/en-us/library/windows/desktop/ms684836(v=vs.85).aspx
func Process32Next(snapshot HANDLE, pe *PROCESSENTRY32) bool {
	ret, _, _ := procProcess32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(pe)),
	)
	return ret != 0
}

// Module32First https://msdn.microsoft.com/en-us/library/windows/desktop/ms684218(v=vs.85).aspx
func Module32First(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)),
	)

	return ret != 0
}

// Module32Next https://msdn.microsoft.com/en-us/library/windows/desktop/ms684221(v=vs.85).aspx
func Module32Next(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)),
	)

	return ret != 0
}

// OpenProcess https://msdn.microsoft.com/en-us/library/windows/desktop/ms684320(v=vs.85).aspx
func OpenProcess(desiredAccess uint32, inheritHandle bool, processID uint32) HANDLE {
	inherit := 0
	if inheritHandle {
		inherit = 1
	}

	ret, _, _ := procOpenProcess.Call(
		uintptr(desiredAccess),
		uintptr(inherit),
		uintptr(processID),
	)
	return HANDLE(ret)
}

// ReadProcessMemory https://msdn.microsoft.com/en-ca/library/windows/desktop/ms680553(v=vs.85).aspx
func ReadProcessMemory(process HANDLE, address uint64, buf []byte, size int) bool {
	ret, _, _ := procReadProcessMemory.Call(
		uintptr(process),
		uintptr(address),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(size),
		uintptr(0),
	)
	return ret != 0
}

// ReadProcessMemoryN https://msdn.microsoft.com/en-ca/library/windows/desktop/ms680553(v=vs.85).aspx
func ReadProcessMemoryN(process HANDLE, address uint64, buf []byte, size int, numberOfBytesRead *uint64) bool {
	ret, _, _ := procReadProcessMemory.Call(
		uintptr(process),
		uintptr(address),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(numberOfBytesRead)),
	)
	return ret != 0
}

// GetProgram returns a HANDLE to the first matching process in processes
func GetProgram(processes ...string) (HANDLE, uint64) {
	snapshot := CreateToolhelp32Snapshot(TH32CSSnapProcess, 0)
	if snapshot == 0 {
		panic(errors.New("snapshot could not be created"))
	}

	var pe PROCESSENTRY32
	var handle HANDLE
	pe.Size = uint32(unsafe.Sizeof(pe))
	if !Process32First(snapshot, &pe) {
		panic(errors.New("Process32First"))
	}
findProcess:
	for Process32Next(snapshot, &pe) {
		for _, v := range processes {
			if syscall.UTF16ToString(pe.ExeFile[:]) == v {
				handle = OpenProcess(ProcessVMRead, false, pe.ProcessID)
				break findProcess
			}
		}
	}
	CloseHandle(snapshot)

	var me MODULEENTRY32
	var baseAddress uint64

	snapshot = CreateToolhelp32Snapshot(TH32CSSnapModule, pe.ProcessID)
	me.Size = uint32(unsafe.Sizeof(me))
	if Module32First(snapshot, &me) {
		module := syscall.UTF16ToString(me.Module[:])
		for _, v := range processes {
			if module == v {
				baseAddress = uint64(me.HModule)
			}
		}
	}
	CloseHandle(snapshot)

	return handle, baseAddress
}
