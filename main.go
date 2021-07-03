package main

func main() {
	view := NewView()
	controller := Controller{view, false}

	for {
		// process user input
		controller.processInput()

		// update

		// convert grid into setContent calls
		view.render()
	}
}
