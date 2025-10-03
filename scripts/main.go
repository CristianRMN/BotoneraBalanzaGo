package main

import (
	"botonera-balanza/scripts/model"
	"botonera-balanza/scripts/service"
	"botonera-balanza/scripts/ui"
	"botonera-balanza/scripts/utils"

	"fyne.io/fyne/v2/app"
)

// metodo que ejecuta la aplicacion
func main() {
	model.InitTurnos()
	myApp := app.New()
	if utils.InitApp(myApp) {
		ui.LoginUI(myApp)
	} else {
		service.AuxLevantarnavegador()
		ui.ShowBotonera(myApp)
	}

	myApp.Run()
}
