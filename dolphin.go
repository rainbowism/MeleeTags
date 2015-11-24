package main

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
)

const noOffset = math.MaxUint64

// Dolphin is for reading from Dolphin's memory
type Dolphin struct {
	process     HANDLE
	baseAddress uint64
	buf         []byte
	bytesRead   uint64
}

// Read returns true if it read memory
func (d *Dolphin) Read(address uint64) bool {
	if ReadProcessMemoryN(d.process, address+d.baseAddress, d.buf, len(d.buf), &d.bytesRead); int(d.bytesRead) != len(d.buf) {
		return false
	}
	return true
}

// ReadBuf returns true if it read memory, as well as the buffer
func (d *Dolphin) ReadBuf(address uint64) (bool, []byte) {
	if ReadProcessMemoryN(d.process, address+d.baseAddress, d.buf, len(d.buf), &d.bytesRead); d.bytesRead == 0 {
		return false, nil
	}
	return true, d.buf
}

// ReadOffset returns true if it read memory
func (d *Dolphin) ReadOffset(address, offset uint64) bool {
	if ReadProcessMemoryN(d.process, address+d.baseAddress, d.buf, len(d.buf), &d.bytesRead); int(d.bytesRead) != len(d.buf) {
		return false
	}
	if offset != noOffset {
		if ReadProcessMemoryN(d.process, binary.LittleEndian.Uint64(d.buf)+offset, d.buf, len(d.buf), &d.bytesRead); int(d.bytesRead) != len(d.buf) {
			return false
		}
	}
	return true
}

// ReadBufOffset returns true if it read memory, as well as the buffer
func (d *Dolphin) ReadBufOffset(address, offset uint64) (bool, []byte) {
	if ReadProcessMemoryN(d.process, address+d.baseAddress, d.buf, len(d.buf), &d.bytesRead); d.bytesRead == 0 {
		return false, nil
	}
	if offset != noOffset {
		if ReadProcessMemoryN(d.process, binary.LittleEndian.Uint64(d.buf)+offset, d.buf, len(d.buf), &d.bytesRead); d.bytesRead == 0 {
			return false, nil
		}
	}
	return true, d.buf
}

// Close cleans up Dolphin
func (d *Dolphin) Close() {
	CloseHandle(d.process)
}

// NewDolphin returns a new Dolphin instance
func NewDolphin() (*Dolphin, error) {
	var dolphin HANDLE
	var baseAddress uint64
	if dolphin, baseAddress = GetProgram("Dolphin.exe"); dolphin == 0 {
		return nil, errors.New("Dolphin isn't running")
	}
	log.Printf("%x\n", baseAddress)
	return &Dolphin{
		process:     dolphin,
		baseAddress: baseAddress,
		buf:         make([]byte, 4),
	}, nil
}
