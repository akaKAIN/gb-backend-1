package main

import (
	"sync"

	"github.com/akaKAIN/gb-backend-1/lesson4/example"
)

func main() {
	var wg = new(sync.WaitGroup)

	wg.Add(1)
	example.StartSimpleServe(wg)
	wg.Wait()
}
