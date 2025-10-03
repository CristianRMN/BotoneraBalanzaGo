package model

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// modelo de datos de las secciones
type Seccion struct {
	Code               int    `json:"code"`
	Name               string `json:"name"`
	ColorTexto         string `json:"color_texto"`
	ColorFondo         string `json:"color_fondo"`
	Icon               string `json:"icon"`
	Sonido             string `json:"sonido"`
	ActivoTurno        string `json:"activo_turno"`
	ActivoPubli        string `json:"activo_publi"`
	Description        string `json:"description"`
	SubeTurno          string `json:"sube_turno"`
	BajaTurno          string `json:"baja_turno"`
	SubeTurnoMando     string `json:"sube_turno_mando"`
	BajaTurnoMando     string `json:"baja_turno_mando"`
	CeroTurnoMando     string `json:"cero_turno_mando"`
	ActivoPubliBalanza string `json:"activo_publi_balanza"`
	TipoPubliBalanza   string `json:"tipo_publi_balanza"`
	Mute               string `json:"mute"`
	Ruta               string `json:"ruta"`
	DobleTicket        string `json:"dobleTicket"`
	Autoservicio       string `json:"autoservicio"`
	Balanza            int    `json:"balanza"`
	Orden              int    `json:"orden"`
	SeccionPrincipal   string `json:"seccionPrincipal"`
	TurnoActual        int    `json:"turnoActual"`
	TurnoSeccion       string `json:"turno_seccion"`
	PubliSeccion       string `json:"publi_seccion"`
	IdParent           *int   `json:"id_parent"`
	Pendings           int    `json:"pendings"`
	ControlSeccion     string `json:"control_seccion"`
}

// modelo de datos del resultado de los datos
type ResultadoData struct {
	ID        int       `json:"id"`
	Tipo      int       `json:"tipo"`
	OnOpen    bool      `json:"onOpen"`
	Secciones []Seccion `json:"secciones"`
}

// modelo de datos de resultado
type Resultado struct {
	Rows int           `json:"rows"`
	Data ResultadoData `json:"data"`
}

// modelo de datos de respuesta de la API
type RespuestaAPI struct {
	Status string    `json:"status"`
	Result Resultado `json:"result"`
}

// modelo de datos de la conexion a la API
type ModeloConexionAPI struct {
	Url      string
	LblTurno *widget.Label
	LblError *widget.Label
	Window   fyne.Window
}

// modelo de datos del back de la configuración de la aplicación
type ConfigModelo struct {
	Ip                string
	IpBalanza         string
	Token             string
	Usuario           string
	Password          string
	Botonera          bool
	ColorRBotonera    *int
	ColorGBotonera    *int
	ColorBBotonera    *int
	PosicionXBotonera *int
	PosicionYBotonera *int
	Visor             bool
	PosicionXVisor    *int
	PosicionYVisor    *int
}

// modelo de datos del front de la configuración de la aplicación
type ConfigModeloUI struct {
	Ip                *widget.Entry
	IpBalanza         *widget.Entry
	Token             *widget.Entry
	Usuario           *widget.Entry
	Password          *widget.Entry
	Botonera          *widget.Check
	ColorRBotonera    *widget.Entry
	ColorGBotonera    *widget.Entry
	ColorBBotonera    *widget.Entry
	PosicionXBotonera *widget.Entry
	PosicionYBotonera *widget.Entry
	Visor             *widget.Check
	PosicionXVisor    *widget.Entry
	PosicionYVisor    *widget.Entry
}

// Modelo de datos de la configuración a la hora de guardar
type ConfigModeloSave struct {
	Ip                *widget.Entry
	IpBalanza         *widget.Entry
	Token             *widget.Entry
	Usuario           *widget.Entry
	Password          *widget.Entry
	Botonera          *widget.Check
	ColorRBotonera    *widget.Entry
	ColorGBotonera    *widget.Entry
	ColorBBotonera    *widget.Entry
	PosicionXBotonera *widget.Entry
	PosicionYBotonera *widget.Entry
	Visor             *widget.Check
	PosicionXVisor    *widget.Entry
	PosicionYVisor    *widget.Entry
	LoginWindow       fyne.Window
}

// modelo de datos del settings de configuracion de ruedita
type ModalSettingsWidgets struct {
	Container     *fyne.Container
	TextUsuario   *widget.Entry
	TextPassword  *widget.Entry
	ButtonAceptar *widget.Button
	ButtonCerrar  *widget.Button
	ButtonSupremo *widget.Button
}

// modelo de datos del JSON que devuelve los endPoints de subir y bajar turno
type TurnoResponse struct {
	Status string `json:"status"`
	Result struct {
		Rows int `json:"rows"`
		Data struct {
			Turn int `json:"turn"`
		} `json:"data"`
	} `json:"result"`
}
