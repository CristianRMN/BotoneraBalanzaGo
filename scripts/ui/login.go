package ui

import (
	"botonera-balanza/scripts/model"
	"botonera-balanza/scripts/service"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Función principal de UI
func LoginUI(myApp fyne.App) {
	myApp.Settings().SetTheme(theme.DarkTheme())

	//secciones := service.GetConsultiva(endPointConsultiva)
	btncerrarApp := widget.NewButton("Cerrar aplicación", func() {
		if model.NavegadorCmd != nil && model.NavegadorCmd.Process != nil {
			_ = model.NavegadorCmd.Process.Kill()
		}
		os.Exit(0)
	})

	btnAceptar := widget.NewButton("Aceptar", nil)

	configSave := ConstruirVentanaLogin(myApp, btncerrarApp, btnAceptar)

	arrayTextosBotonera := []*widget.Entry{configSave.ColorRBotonera, configSave.ColorGBotonera,
		configSave.ColorBBotonera, configSave.PosicionYBotonera, configSave.PosicionXBotonera}
	arrayTextosVisor := []*widget.Entry{configSave.PosicionYVisor, configSave.PosicionXVisor}

	service.DisableEntrysByCheck(configSave.Visor, arrayTextosVisor)
	service.DisableEntrysByCheck(configSave.Botonera, arrayTextosBotonera)

	btnAceptar.OnTapped = func() {
		if service.EmptyText(configSave.Usuario.Text) ||
			service.EmptyText(configSave.Password.Text) || service.EmptyText(configSave.Ip.Text) || service.EmptyText(configSave.IpBalanza.Text) {
			dialog.ShowInformation("Campos vacíos", "Por favor, rellena todos los campos.", configSave.LoginWindow)
			return
		}

		if !configSave.Botonera.Checked && !configSave.Visor.Checked {
			dialog.ShowInformation("Configuración incorrecta", "Tienes que tener seleccionado turnos o publicidad", configSave.LoginWindow)
			return
		}

		if !service.CheckLoginIpValidSecciones(configSave.Ip.Text, configSave.IpBalanza.Text) {
			dialog.ShowInformation("Error", "Este terminal no tiene secciones, crealas primero", configSave.LoginWindow)
			return
		}

		if configSave.Botonera.Checked {
			if service.EmptyText(configSave.ColorRBotonera.Text) || service.EmptyText(configSave.ColorGBotonera.Text) ||
				service.EmptyText(configSave.ColorBBotonera.Text) ||
				service.EmptyText(configSave.PosicionYBotonera.Text) ||
				service.EmptyText(configSave.PosicionXBotonera.Text) {
				dialog.ShowInformation("Campos vacíos", "La botonera está activada y tienes que rellenar los campos correspondientes", configSave.LoginWindow)
				return
			}
			if !service.IsValidRGBLoginEntry(configSave.ColorRBotonera) ||
				!service.IsValidRGBLoginEntry(configSave.ColorGBotonera) ||
				!service.IsValidRGBLoginEntry(configSave.ColorBBotonera) {
				dialog.ShowInformation("Colores incorrectos", "Los colores deben de ir del 0 al 255", configSave.LoginWindow)
				return
			}
			if !service.CheckIfPisitionsLoginIsNumber(configSave.PosicionYBotonera) ||
				!service.CheckIfPisitionsLoginIsNumber(configSave.PosicionXBotonera) {
				dialog.ShowInformation("Posiciones incorrectas", "La posición debe de ser un número", configSave.LoginWindow)
				return
			}
			if !service.CheckRestrictionPositionY(configSave.PosicionYBotonera) {
				dialog.ShowInformation("Posición no permitida", "El vertical va de 20 a 32", configSave.LoginWindow)
				return
			}

			/*
				valor, numero := service.CheckLimitBotonera(secciones, configSave.PosicionXBotonera.Text)
				if !valor {
					dialog.ShowInformation("Posición no permitida", "El valor va de 0 a "+strconv.Itoa(numero), configSave.LoginWindow)
					return
				}
			*/
		}

		if configSave.Visor.Checked {
			if service.EmptyText(configSave.PosicionYVisor.Text) ||
				service.EmptyText(configSave.PosicionXVisor.Text) {
				dialog.ShowInformation("Campos vacíos", "El visor está activado y tienes que rellenar los campos correspondientes", configSave.LoginWindow)
				return
			}
			if !service.CheckIfPisitionsLoginIsNumber(configSave.PosicionYVisor) ||
				!service.CheckIfPisitionsLoginIsNumber(configSave.PosicionXVisor) {
				dialog.ShowInformation("Posiciones incorrectas", "La posición debe de ser un número", configSave.LoginWindow)
				return
			}
		}

		service.SendDatosConfig(configSave)

		ShowBotonera(myApp)
	}

	if !service.IsFileEmpty(service.GetRutaDependsModo()) {
		btncerrarApp.Show()
		config, err := service.LoadValuesFichero(service.GetRutaDependsModo())
		if err == nil {
			configUI := model.ConfigModeloUI{
				Ip: configSave.Ip, IpBalanza: configSave.IpBalanza, Token: configSave.Token,
				Usuario: configSave.Usuario, Password: configSave.Password,
				Botonera: configSave.Botonera, ColorRBotonera: configSave.ColorRBotonera,
				ColorGBotonera: configSave.ColorGBotonera, ColorBBotonera: configSave.ColorBBotonera,
				PosicionXBotonera: configSave.PosicionXBotonera,
				PosicionYBotonera: configSave.PosicionYBotonera,
				Visor:             configSave.Visor, PosicionXVisor: configSave.PosicionXVisor,
				PosicionYVisor: configSave.PosicionYVisor,
			}
			service.SetValuesUI(*config, configUI)
		}
	} else {
		btncerrarApp.Hide()
	}

}
