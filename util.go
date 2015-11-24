package main

import (
	"time"
)

func RunForever(fn func(), interval int) {
	var currentTime, accumTime, deltaTime int64
	currentTime = time.Now().UnixNano()
	accumTime = 0
	deltaTime = 1000000000 / int64(interval)

	for {
		newTime := time.Now().UnixNano()
		frameTime := newTime - currentTime
		currentTime = newTime
		accumTime += frameTime

		for accumTime >= deltaTime {
			fn()
			accumTime -= deltaTime
		}
	}
}
