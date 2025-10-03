package ui

import (
	"botonera-balanza/scripts/model"
	"botonera-balanza/scripts/service"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// metodo que construye la ventana del Login
func ConstruirVentanaLogin(myApp fyne.App, btnCerrar *widget.Button, btnAceptar *widget.Button) model.ConfigModeloSave {
	loginWindow := myApp.NewWindow("Configuración Inicial")

	loginWindow.SetIcon(theme.SettingsIcon())
	textIp, containerIp := service.NewLabeledEntry("*Ip Servidor: ", "Ip aquí...", 220)
	textIpBalanza, containerIpBalanza := service.NewLabeledEntry("*Ip balanza: ", "Ip aquí...", 220)

	textToken, containerToken := service.NewLabeledEntry("Token: ", "Token aquí...", 220)
	textUsuario, containerUsuario := service.NewLabeledEntry("*Usuario", "Usuario aquí...", 220)
	textPassword, containerPassword := service.PasswordLeveledEntry("*Contraseña", "Contraseña aquí...", 220)

	botoneraCheck := widget.NewCheck("Botonera", nil)

	rText, rContainer := service.NewSizedEntry("Color r aquí...", 120)
	gText, gContainer := service.NewSizedEntry("Color g aquí...", 120)
	bText, bContainer := service.NewSizedEntry("Color b aquí...", 120)

	contentHBoxColores := container.NewHBox(
		rContainer, layout.NewSpacer(), gContainer, layout.NewSpacer(), bContainer,
	)

	verticalBotonera, verticalBotoneraContainer := service.NewSizedEntry("Vertical...", 120)
	horizontalBotonera, horizontalBotoneraContainer := service.NewSizedEntry("Horizontal...", 120)
	contentHBoxPosicionesBotonera := container.NewHBox(
		verticalBotoneraContainer, layout.NewSpacer(), horizontalBotoneraContainer, layout.NewSpacer(),
	)

	chechVisor := widget.NewCheck("Visor", nil)

	verticalVisor, verticalVisorContainer := service.NewSizedEntry("Vertical...", 120)
	horizontalVisor, horizontalVisorContainer := service.NewSizedEntry("Horizontal...", 120)
	contentHBoxPosicionesVisor := container.NewHBox(
		verticalVisorContainer, layout.NewSpacer(), horizontalVisorContainer, layout.NewSpacer(),
	)

	form := container.NewVBox(
		containerIp, containerIpBalanza, containerToken, containerUsuario, containerPassword,
		botoneraCheck, widget.NewLabel("*Color de fondo"), contentHBoxColores,
		widget.NewLabel("Posición"), contentHBoxPosicionesBotonera,
		chechVisor, widget.NewLabel("Posición"), contentHBoxPosicionesVisor,
		btnAceptar, btnCerrar,
	)

	loginWindow.SetContent(form)
	loginWindow.Resize(fyne.NewSize(360, 300))
	loginWindow.SetFixedSize(true)
	loginWindow.Show()

	return model.ConfigModeloSave{
		Ip: textIp, IpBalanza: textIpBalanza, Token: textToken, Usuario: textUsuario, Password: textPassword,
		Botonera:       botoneraCheck,
		ColorRBotonera: rText, ColorGBotonera: gText, ColorBBotonera: bText,
		PosicionXBotonera: horizontalBotonera, PosicionYBotonera: verticalBotonera,
		Visor:          chechVisor,
		PosicionXVisor: horizontalVisor, PosicionYVisor: verticalVisor,
		LoginWindow: loginWindow,
	}
}

// Metodo que construye la ventana de la botonera
func CrearVentanaBotonera(myApp fyne.App, valorBotonera string, saliBotonera bool, labelColores *widget.Label,
	seccion model.Seccion, ventanas *[]fyne.Window, urlSubirTurno string, urlBajarTurno string) fyne.Window {

	myWindow := myApp.NewWindow(seccion.Name)
	myWindow.SetIcon(theme.DesktopIcon())
	*ventanas = append(*ventanas, myWindow)
	service.OcultarBotoneraIfIsFalse(myWindow, valorBotonera)

	turnos := widget.NewLabel("000")
	service.SetTurnoActual(turnos, seccion)
	labelMensajesError := widget.NewLabel("ejemplo de error")

	modeloAPISubida := model.ModeloConexionAPI{
		Url:      endPointSubirTurno,
		LblTurno: turnos,
		LblError: labelMensajesError,
		Window:   myWindow,
	}

	modeloAPIBajada := model.ModeloConexionAPI{
		Url:      endPointBajarTurno,
		LblTurno: turnos,
		LblError: labelMensajesError,
		Window:   myWindow,
	}

	contenHBox := container.NewHBox(
		layout.NewSpacer(),
		container.NewMax(labelMensajesError),
		layout.NewSpacer(),
	)
	labelMensajesError.Hide()

	ruedaPersonalizada := crearBotonConIcono(myWindow, myApp, saliBotonera, ventanas)

	btnBajarTurno := widget.NewButton("-", func() {
		service.BajarTurno(modeloAPIBajada, ruedaPersonalizada, urlBajarTurno, seccion)
	})

	btnSubirTurno := widget.NewButton("+", func() {
		service.SubirTurno(modeloAPISubida, ruedaPersonalizada, urlSubirTurno, seccion)
	})

	botonBajarCss := estiloBotonBajarTurno(btnBajarTurno)
	botonSubirCss := estiloBotonSubirTurno(btnSubirTurno)
	turnoPersonalizado := estiloLabelTurnos(turnos)

	content := container.NewHBox(
		ruedaPersonalizada,
		layout.NewSpacer(),
		botonBajarCss,
		layout.NewSpacer(),
		turnoPersonalizado,
		layout.NewSpacer(),
		botonSubirCss,
	)

	contentPrincipal := container.NewVBox(
		content,
		contenHBox,
	)

	myWindow.SetContent(contentPrincipal)
	myWindow.Resize(fyne.NewSize(210, 40))
	myWindow.SetFixedSize(true)
	return myWindow
}

// Metodo que hace el estilo de bajar turno
func estiloBotonBajarTurno(bajarB *widget.Button) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	bg.SetMinSize(fyne.NewSize(80, 30))

	btnPersonalizado := container.NewMax(bg, bajarB)
	return btnPersonalizado
}

// Metodo que hace el estilo de subir turno
func estiloBotonSubirTurno(bajarB *widget.Button) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255})
	bg.SetMinSize(fyne.NewSize(80, 30))

	btnPersonalizado := container.NewMax(bg, bajarB)
	return btnPersonalizado
}

// Metodo que hace el estilo de los labels de turnos
func estiloLabelTurnos(lbl *widget.Label) fyne.CanvasObject {
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	lbl.Alignment = fyne.TextAlignCenter
	return container.NewCenter(lbl)
}

// Metodo para crear el boton del icono de configuracion de la botonera
func crearBotonConIcono(window fyne.Window, myApp fyne.App, saliBotonera bool, ventanas *[]fyne.Window) *widget.Button {

	btn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		window.Resize(fyne.NewSize(210, 340))

		showModalSettings(window, myApp, saliBotonera, ventanas)
	})
	btn.Resize(fyne.NewSize(30, 30))
	return btn
}

// Metodo que crea la UI de los settings para acceder al fichero de configuracion
func CrearUiModalSettings(win fyne.Window) model.ModalSettingsWidgets {
	textUsuario := widget.NewEntry()
	textPassword := widget.NewPasswordEntry()

	buttonSupremo := widget.NewButton("(▲) Control remoto", nil)

	buttonSupremo.OnTapped = func() {
		service.ShowSupremo(win, buttonSupremo)
	}
	buttonAceptar := widget.NewButton("Aceptar", nil)
	buttonCerrar := widget.NewButton("Cerrar", nil)

	formContainer := container.NewVBox(
		widget.NewLabel("Usuario"),
		textUsuario,
		widget.NewLabel("Password"),
		textPassword,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer()),
	)

	dialogContainer := container.NewVBox(
		formContainer,
		container.NewHBox(layout.NewSpacer(), buttonAceptar, buttonCerrar),
		buttonSupremo,
	)

	return model.ModalSettingsWidgets{
		Container:     dialogContainer,
		TextUsuario:   textUsuario,
		TextPassword:  textPassword,
		ButtonAceptar: buttonAceptar,
		ButtonCerrar:  buttonCerrar,
		ButtonSupremo: buttonSupremo,
	}

}
