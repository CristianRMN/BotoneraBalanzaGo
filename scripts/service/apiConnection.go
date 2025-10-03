package service

import (
	"botonera-balanza/scripts/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"fyne.io/fyne/v2/widget"
)

/*
Metodo que comprueba si el servidor está iniciado
*/
func IsServerAlive(url string) bool {
	fmt.Println("Verificando URL:", url)
	client := http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		LogError(err.Error())
		return false
	}
	defer resp.Body.Close()

	return true
}

/*
Metodo que comprueba las secciones que tienes
*/
func GetConsultiva(url string) []model.Seccion {
	token := "$2y$12$e15DVhqmVM2ZGGPCYVqo4um3cn0vSt/nd6ZEc8qFTraKggYK91oOK"
	var auxSeccion []model.Seccion

	if !IsServerAlive(url) {
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		LogError(err.Error())
		return nil
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		LogError(err.Error())
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError(err.Error())
		return nil
	}

	var respuesta model.RespuestaAPI
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		LogError(err.Error())
		return nil
	}

	for i := 0; i < len(respuesta.Result.Data.Secciones); i++ {
		if respuesta.Result.Data.Secciones[i].ControlSeccion == "T" {
			auxSeccion = append(auxSeccion, respuesta.Result.Data.Secciones[i])
		}
	}

	return auxSeccion

}

/*
Metodo que llama al endPoint de subir turno
Devuelve los mensajes correspondientes al resultado de la llamada
*/
func SubirTurno(modeloAPI model.ModeloConexionAPI, ruedaPersonal *widget.Button, urlSubirTurno string, seccion model.Seccion) {

	token := "$2y$12$e15DVhqmVM2ZGGPCYVqo4um3cn0vSt/nd6ZEc8qFTraKggYK91oOK"

	/*
		if !IsServerAlive(urlSubirTurno) {
			SetTextlabelError(modeloAPI.LblError, "Error al conectarse al servidor", 2, modeloAPI.Window, ruedaPersonal)
			return
		}
	*/
	client := &http.Client{}

	req, err := http.NewRequest("GET", urlSubirTurno, nil)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Error al conectarse al servidor", 2, modeloAPI.Window, ruedaPersonal)
		return
	}

	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Servidor no disponible o fuera de red", 2, modeloAPI.Window, ruedaPersonal)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Error al procesar la solicitud", 2, modeloAPI.Window, ruedaPersonal)
		return
	}

	var resultado model.TurnoResponse
	err = json.Unmarshal(body, &resultado)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Error al procesar la solicitud", 2, modeloAPI.Window, ruedaPersonal)
		return
	}

	if resultado.Status != "ok" {
		SwitchErrorSubirBajarTurno(modeloAPI, resp.StatusCode, ruedaPersonal, "aumentar")
		return
	}

	modeloAPI.LblTurno.Text = fmt.Sprintf("%03d", resultado.Result.Data.Turn)
	modeloAPI.LblTurno.Refresh()
	model.SetTurno(seccion.Name, fmt.Sprintf("%03d", resultado.Result.Data.Turn))
}

/*
Metodo que llama al endPoint de bajar turno
Devuelve los mensajes correspondientes al resultado de la llamada
*/
func BajarTurno(modeloAPI model.ModeloConexionAPI, ruedaPersonal *widget.Button, urlBajarTurno string, seccion model.Seccion) {

	token := "$2y$12$e15DVhqmVM2ZGGPCYVqo4um3cn0vSt/nd6ZEc8qFTraKggYK91oOK"

	/*
		if !IsServerAlive(urlBajarTurno) {
			SetTextlabelError(modeloAPI.LblError, "Error al conectarse al servidor", 2, modeloAPI.Window, ruedaPersonal)
			return
		}
	*/

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlBajarTurno, nil)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Error al conmectarse al servidor", 2, modeloAPI.Window, ruedaPersonal)
		return
	}

	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Servidor no disponible o fuera de red", 2, modeloAPI.Window, ruedaPersonal)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Error al procesar la solicitud", 2, modeloAPI.Window, ruedaPersonal)
		return
	}

	var resultado model.TurnoResponse
	err = json.Unmarshal(body, &resultado)
	if err != nil {
		SetTextlabelError(modeloAPI.LblError, "Error al procesar la solicitud", 2, modeloAPI.Window, ruedaPersonal)
		return
	}

	if resultado.Status != "ok" {
		SwitchErrorSubirBajarTurno(modeloAPI, resp.StatusCode, ruedaPersonal, "bajar")
		return
	}

	modeloAPI.LblTurno.Text = fmt.Sprintf("%03d", resultado.Result.Data.Turn)
	modeloAPI.LblTurno.Refresh()
	model.SetTurno(seccion.Name, fmt.Sprintf("%03d", resultado.Result.Data.Turn))

}

func SwitchErrorSubirBajarTurno(modeloAPI model.ModeloConexionAPI, resp int, ruedaPersonal *widget.Button, verbo string) {
	switch resp {
	case 400:
		SetTextlabelError(modeloAPI.LblError, "El parámetro 'section' está mal escrito o es null", 2, modeloAPI.Window, ruedaPersonal)
	case 401:
		SetTextlabelError(modeloAPI.LblError, "No autorizado", 2, modeloAPI.Window, ruedaPersonal)
	case 403:
		SetTextlabelError(modeloAPI.LblError, "No se puede "+verbo+" más turnos", 2, modeloAPI.Window, ruedaPersonal)
	case 404:
		SetTextlabelError(modeloAPI.LblError, "La IP no existe en la BD", 2, modeloAPI.Window, ruedaPersonal)
	case 405:
		SetTextlabelError(modeloAPI.LblError, "Método no permitido", 2, modeloAPI.Window, ruedaPersonal)
	case 422:
		SetTextlabelError(modeloAPI.LblError, "Alguno de los campos ingresados no son correctos", 2, modeloAPI.Window, ruedaPersonal)
	case 500:
		SetTextlabelError(modeloAPI.LblError, "Error interno del servidor al ejecutar la consulta.", 2, modeloAPI.Window, ruedaPersonal)
	case 503:
		SetTextlabelError(modeloAPI.LblError, "Error de conexión al conectarse al servidor", 2, modeloAPI.Window, ruedaPersonal)
	}

}
