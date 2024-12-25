package main

import (
	"github.com/voice0726/todo-app-api/di"
)

func main() {
	app := di.PrepareApp()
	app.Run()
}
