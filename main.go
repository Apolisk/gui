package main

import (
	"fyne.io/fyne/v2/app"
	"main/gui"
)

func main() {
	myApp := app.New()
	gui.Ui(myApp)

}
