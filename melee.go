package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"time"
)

var CharMap = map[uint16]string{
	0x200A: "A",
	0x200B: "B",
	0x200C: "C",
	0x200D: "D",
	0x200E: "E",
	0x200F: "F",
	0x2010: "G",
	0x2011: "H",
	0x2012: "I",
	0x2013: "J",
	0x2014: "K",
	0x2015: "L",
	0x2016: "M",
	0x2017: "N",
	0x2018: "O",
	0x2019: "P",
	0x201A: "Q",
	0x201B: "R",
	0x201C: "S",
	0x201D: "T",
	0x201E: "U",
	0x201F: "V",
	0x2020: "W",
	0x2021: "X",
	0x2022: "Y",
	0x2023: "Z",
	0x2000: "0",
	0x2001: "1",
	0x2002: "2",
	0x2003: "3",
	0x2004: "4",
	0x2005: "5",
	0x2006: "6",
	0x2007: "7",
	0x2008: "8",
	0x2009: "9",
	0x20E3: " ",
	0x20E7: ".",
	0x20EB: "?",
	0x20EC: "!",
	0x20FB: "+",
	0x20FC: "-",
	0x20FE: "=",
	0x2102: "$",
	0x2103: "%",
	0x2105: "&",
	0x2107: "@",
}

type MeleeTags struct {
	Dolphin *Dolphin
	x64     bool // 64-bit Dolphin?
}

func NewMeleeTags(x64 bool) (*MeleeTags, error) {
	var dolphin *Dolphin
	var err error
	if dolphin, err = NewDolphin(); err != nil {
		return nil, err
	}
	return &MeleeTags{
		Dolphin: dolphin,
		x64:     x64,
	}, nil
}

func (melee *MeleeTags) Close() {
	melee.Dolphin.Close()
}

func (melee *MeleeTags) Update() (ok bool, err error) {
	// Player 1
	player1 := melee.ReadPlayer(0x80C45D37, 0x80C45D5F+(0x4B0*0))
	err = ioutil.WriteFile("player1.txt", []byte(player1), 0644)
	if err != nil {
		return false, err
	}
	player2 := melee.ReadPlayer(0x80C45EC3, 0x80C45D5F+(0x4B0*1))
	err = ioutil.WriteFile("player2.txt", []byte(player2), 0644)
	if err != nil {
		return false, err
	}
	player3 := melee.ReadPlayer(0x80C46373, 0x80C45D5F+(0x4B0*2))
	err = ioutil.WriteFile("player3.txt", []byte(player3), 0644)
	if err != nil {
		return false, err
	}
	player4 := melee.ReadPlayer(0x80C46823, 0x80C45D5F+(0x4B0*3))
	err = ioutil.WriteFile("player4.txt", []byte(player4), 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (melee *MeleeTags) ReadPlayer(charCount, base uint64) string {
	if melee.x64 {
		charCount += 0x100000000
		base += 0x100000000
	}
	melee.Dolphin.Read(charCount)
	charNum := melee.Dolphin.buf[0] & 0xF
	var nametag string
	switch charNum {
	case 0x3:
		melee.Dolphin.Read(base + 0x0)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		break
	case 0x6:
		melee.Dolphin.Read(base + 0x0)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		melee.Dolphin.Read(base + 0x3)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		break
	case 0x9:
		melee.Dolphin.Read(base + 0x0)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		melee.Dolphin.Read(base + 0x3)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		melee.Dolphin.Read(base + 0x6)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		break
	case 0xC:
		melee.Dolphin.Read(base + 0x0)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		melee.Dolphin.Read(base + 0x3)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		melee.Dolphin.Read(base + 0x6)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		melee.Dolphin.Read(base + 0x9)
		nametag += CharMap[binary.BigEndian.Uint16(melee.Dolphin.buf[:2])]
		break
	}
	return nametag
}

func (melee *MeleeTags) Run() {
	var ok bool
	var err error
	ticker := time.NewTicker(time.Second / time.Duration(8))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if ok, err = melee.Update(); !ok {
				fmt.Println(err)
				return
			}
		}
	}
}
