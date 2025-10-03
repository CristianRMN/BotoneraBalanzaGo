package utils

import (
	"botonera-balanza/scripts/service"

	"fyne.io/fyne/v2"
)

// Metodo que comprueba si el archivo de configuración está vacío
func InitApp(myApp fyne.App) bool {
	return service.IsFileEmpty(service.GetRutaDependsModo())
}
