package main

import "log"

func main() {
	// x64 := flag.Bool("x64", false, "64-bit Dolphin mode")
	// flag.Parse()

	// if *x64 {
	// 	log.Println("Running MeleeTags with 64-bit Dolphin compatibility.")
	// } else {
	// 	log.Println("Running MeleeTags with 32-bit Dolphin compatibility.")
	// }
	log.Println("Running MeleeTags with 64-bit Dolphin compatibility.")
	melee, err := NewMeleeTags(true)
	if err != nil {
		log.Fatal(err)
	}
	defer melee.Close()

	melee.Run()
}
