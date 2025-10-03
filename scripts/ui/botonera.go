package ui

import (
	"botonera-balanza/scripts/model"
	"botonera-balanza/scripts/service"

	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lxn/win"
)

var confPath = service.GetRutaDependsModo()

// Endpoints para lalmar turnos y consultas
var (
	endPointSubirTurno = "http://" + service.GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia/api/upTurnDesktop?section=1&ip=" + service.GetValueNotIntDesdeConf(confPath, "ipbalanza=")
	endPointBajarTurno = "http://" + service.GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia/api/downTurnDesktop?section=1&ip=" + service.GetValueNotIntDesdeConf(confPath, "ipbalanza=")
	endPointConsultiva = "http://" + service.GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia/api/consultiva?ip=" + service.GetValueNotIntDesdeConf(confPath, "ipbalanza=")
	urlGmedia          = "http://" + service.GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia"
)

// metodo que crea la botonera
func ShowBotonera(myApp fyne.App) {
	endPointSubirTurno, endPointBajarTurno, endPointConsultiva, urlGmedia = service.UpdateVarsEndPoints(endPointSubirTurno, endPointBajarTurno, endPointConsultiva, urlGmedia, confPath)
	myApp.Settings().SetTheme(theme.LightTheme())

	secciones := service.GetConsultiva(endPointConsultiva)
	var idsSecciones []int
	var arrayEndPointsSubir []string
	var arrayEndPointsBajar []string

	numeroSeciones := len(secciones)

	if numeroSeciones == 0 {
		return
	}

	for i := 0; i < len(secciones); i++ {
		idsSecciones = append(idsSecciones, service.RellenarIdsSecciones(secciones[i]))
	}

	for i := 0; i < len(idsSecciones); i++ {
		model.UrlSubirTurnoArray = service.MakeArrayEndPoints(idsSecciones[i], endPointSubirTurno)
		model.UrlBajarTurnoArray = service.MakeArrayEndPoints(idsSecciones[i], endPointBajarTurno)
		arrayEndPointsSubir = append(arrayEndPointsSubir, model.UrlSubirTurnoArray)
		arrayEndPointsBajar = append(arrayEndPointsBajar, model.UrlBajarTurnoArray)
	}

	labelColores := widget.NewLabel("EjemploColores")

	posYBotonera := service.GetValorDesdeConf(confPath, "posicionBotoneraY=")
	valorBotonera := service.GetValueNotIntDesdeConf(confPath, "botonera=")
	botoneraColorR := service.GetValueNotIntDesdeConf(confPath, "colorR=")
	botoneraColorG := service.GetValueNotIntDesdeConf(confPath, "colorG=")
	botoneraColorB := service.GetValueNotIntDesdeConf(confPath, "colorB=")
	valorVisor := service.GetValueNotIntDesdeConf(service.GetRutaDependsModo(), "visor=")
	posXVisor := service.GetValueNotIntDesdeConf(service.GetRutaDependsModo(), "visorPosicionX=")
	posYVisor := service.GetValueNotIntDesdeConf(service.GetRutaDependsModo(), "visorPosicionY=")

	if model.CancelGourutine != nil {
		model.CancelGourutine()
	}

	service.Levantarnavegador(valorVisor, posXVisor, posYVisor, urlGmedia)
	model.GoroutineCtx, model.CancelGourutine = context.WithCancel(context.Background())
	var ventanas []fyne.Window

	var posicionesMejoradas []int
	diferencia := 300
	posXBotonera := service.GetValorDesdeConf(confPath, "PosicionBotoneraX=")

	for i := 0; i < numeroSeciones; i++ {
		posX := posXBotonera + i*diferencia
		posicionesMejoradas = append(posicionesMejoradas, posX)

		ventana := CrearVentanaBotonera(myApp, valorBotonera, false, labelColores, secciones[i], &ventanas,
			arrayEndPointsSubir[i], arrayEndPointsBajar[i])

		iCopy := i

		go func(ctx context.Context, seccion model.Seccion, posX int) {
			localSaliBotonera := false
			yaMostrada := false
			hwnd := service.GetHWND(seccion.Name)

			go func() {
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()

				for {

					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						if !service.IsCovered(hwnd) && !model.GetEntreSettings() {
							win.SetWindowPos(hwnd, win.HWND_NOTOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_NOACTIVATE)
							win.SetWindowPos(hwnd, win.HWND_TOPMOST, int32(posX), int32(posYBotonera), 0, 0, win.SWP_NOSIZE)
						}

						if valorBotonera == "true" {
							if hwnd != 0 {
								if !localSaliBotonera {
									style := win.GetWindowLong(hwnd, win.GWL_STYLE)
									style &^= win.WS_SYSMENU | win.WS_MINIMIZEBOX | win.WS_MAXIMIZEBOX
									win.SetWindowLong(hwnd, win.GWL_STYLE, style)
									win.SetWindowPos(hwnd, win.HWND_NOTOPMOST, int32(posX), int32(posYBotonera), 0, 0, win.SWP_NOSIZE|win.SWP_NOACTIVATE)
									localSaliBotonera = true
								}
								win.SetWindowPos(hwnd, win.HWND_TOPMOST, int32(posX), int32(posYBotonera), 0, 0, win.SWP_NOSIZE)
								left, top, right, bottom := service.GetWindowRect(hwnd)
								fmt.Println(left)
								midY := int((top + bottom) / 2)
								rR, gR, bR := service.GetPixelColor(int(right)+1, midY)

								fyne.DoAndWait(func() {
									labelColores.Text = fmt.Sprintf("%d, %d, %d", rR, gR, bR)
									labelColores.Refresh()
									if !model.GetEntreSettings() && !yaMostrada {
										if service.IsColorMatch(rR, gR, bR,
											service.ObtenerColoresBotoneraOcultar(botoneraColorR),
											service.ObtenerColoresBotoneraOcultar(botoneraColorG),
											service.ObtenerColoresBotoneraOcultar(botoneraColorB)) {

											win.ShowWindow(hwnd, win.SW_SHOWNOACTIVATE)
											ventana.Show()
											yaMostrada = true
										} else {
											//service.OcultarVentanaDelay(hwnd)
										}
									}
								})
							}
						}
					}
				}
			}()
		}(model.GoroutineCtx, secciones[iCopy], posX)
	}

}
