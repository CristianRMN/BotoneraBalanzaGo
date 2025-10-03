package ui

import (
	"botonera-balanza/scripts/model"
	"botonera-balanza/scripts/service"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

/*
Metodo principal del modal de ruedita de
*/
func showModalSettings(win fyne.Window, myApp fyne.App, saliBotonera bool, ventanas *[]fyne.Window) {
	configSettings := CrearUiModalSettings(win)
	model.SetEntreSettings(true)

	modal := dialog.NewCustomWithoutButtons("Ajustes", configSettings.Container, win)

	configSettings.ButtonAceptar.OnTapped = func() {
		usuario := service.GetValueNotIntDesdeConf(service.GetRutaDependsModo(), "usuario=")
		password := service.GetValueNotIntDesdeConf(service.GetRutaDependsModo(), "password=")

		if configSettings.TextUsuario.Text == usuario && configSettings.TextPassword.Text == password {
			model.SetEntreSettings(false)
			modal.Hide()
			win.Resize(fyne.NewSize(210, 40))
			saliBotonera = false
			service.CerrarVentanas(*ventanas)
			LoginUI(myApp)
		} else {
			dialog.ShowError(fmt.Errorf("usuario o contraseña incorrectos"), win)
		}
	}

	configSettings.ButtonCerrar.OnTapped = func() {
		model.SetEntreSettings(false)
		modal.Hide()
		win.Resize(fyne.NewSize(210, 40))
	}

	modal.Show()
}
