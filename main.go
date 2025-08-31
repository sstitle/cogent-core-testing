package main

import (
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
)

func main() {
	// Create the main application body
	b := core.NewBody()
	b.SetTitle("Cogent Core Example")

	// Create a button that changes its text when clicked
	core.NewButton(b).SetText("Send").SetIcon(icons.Send).OnClick(func(e events.Event) {
		core.MessageSnackbar(b, "Message sent")
	})

	// Run the application
	b.RunMainWindow()
}
