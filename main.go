package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	view := NewView()
	controller := Controller{view, false}

	refreshRate, err := time.ParseDuration(fmt.Sprintf("%vms", 1000.0/30))
	if err != nil {
		log.Fatalln("could not parse refreshRate")
	}

	for {
		// process user input
		controller.processInput()

		// update
		view.update()

		// convert grid into setContent calls
		view.render()

		time.Sleep(refreshRate)
	}
}
