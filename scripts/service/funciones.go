package service

import (
	"botonera-balanza/scripts/model"
	"errors"

	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/kbinani/screenshot"

	dialogFyne "fyne.io/fyne/v2/dialog"
	dialogSys "github.com/sqweek/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lxn/win"

	webview "github.com/webview/webview_go"
)

// Metodo para escribir en un fichero Log y saber los errores de la aplicacion
func LogError(mensaje string) {
	ruta := filepath.Join("..", "errores")

	err := os.MkdirAll(ruta, os.ModePerm)
	if err != nil {
		return
	}

	archivoLog := filepath.Join(ruta, "log.txt")

	f, err := os.OpenFile(archivoLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log := fmt.Sprintf("[%s] %s\n", timestamp, mensaje)
	f.WriteString(log)
}

// Metodo que comprueba si la Ip de tu terminal es igual a la Ip que escribiste en el modal
func CheckIpIsCorrect(ipText string, ipLocal string) bool {
	return ipText == ipLocal
}

// GetLocalIP devuelve la primera IP local válida (no loopback, no IPv6)
func GetLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			return ip.String()
		}
	}

	return ""
}

// metodo que comprueba si el archivo está vacío
func IsFileEmpty(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return true
	}
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		return true
	}
	return info.Size() == 0
}

func UpdateVarsEndPoints(subirTurno string, bajarTurno string, consulta string, urlGmedia string, confPath string) (string, string, string, string) {
	subirTurno = "http://" + GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia/api/upTurnDesktop?section=1&ip=" + GetValueNotIntDesdeConf(confPath, "ipbalanza=")
	bajarTurno = "http://" + GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia/api/downTurnDesktop?section=1&ip=" + GetValueNotIntDesdeConf(confPath, "ipbalanza=")
	consulta = "http://" + GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia/api/consultiva?ip=" + GetValueNotIntDesdeConf(confPath, "ipbalanza=")
	urlGmedia = "http://" + GetValueNotIntDesdeConf(confPath, "ip=") + "/gmedia"
	return subirTurno, bajarTurno, consulta, urlGmedia
}

func CheckLoginIpValidSecciones(ipServidor string, ipBalanza string) bool {
	auxEndPointConsultiva := "http://" + ipServidor + "/gmedia/api/consultiva?ip=" + ipBalanza
	secciones := GetConsultiva(auxEndPointConsultiva)
	return len(secciones) != 0
}

// metodo que devuelve una ruta de configuracion de dibal u otra dependiendo de si estamos en desarrollo o .exe
func GetRutaDependsModo() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	appConfigDir := filepath.Join(configDir, "Botonera")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return ""
	}
	return filepath.Join(appConfigDir, "dibal.conf")
}

// metodo que devuelve una ruta de supremo u otra dependiendo de si estamos en desarrollo o .exe
func GetRutaSupremoDependsModo() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	appConfigDir := filepath.Join(configDir, "Botonera")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return ""
	}
	return filepath.Join(appConfigDir, "configRutaSupremo.txt")
}

/*
Metodo que ejecuta un Supremo .exe siguiendo los siguientes pasos:
1. Comprueba si en el .txt de la ruta de supremo hay algo, si no hay nada, pasa al siguiente if
2. Comprueba en rutas comunes donde puede estar el .exe, si no es así, pasa al siguiente
3. Abre el explorador de archivos para que el usuario busque el .exe, si lo encuentra, se guarda la ruta en el .txt
*/
func ShowSupremo(window fyne.Window, buttonSupremo *widget.Button) {
	pathTextSupremo, err := GetRutaSupremo()
	if err != nil {
		pathTextSupremo, err = GetRutasComunesSupremo()
		if err != nil {
			pathTextSupremo, err = GetRutaSupremoExplorador()
			if err != nil || pathTextSupremo == "" {
				fmt.Println("No se seleccionó Supremo.exe")
				return
			}
			WriteNewRutaSupremo(pathTextSupremo)
		}
	}

	if pathTextSupremo != "" {
		err := exec.Command(pathTextSupremo).Start()
		if err != nil {
			fmt.Println("Error al ejecutar Supremo:", err)
		}
	} else {
		fmt.Println("No hay ruta válida para Supremo")
	}
}

// Lee la ruta del Supremo.exe desde el archivo
func GetRutaSupremo() (string, error) {
	data, err := os.ReadFile(GetRutaSupremoDependsModo())
	if err != nil {
		return "", err
	}
	path := string(data)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return "", errors.New("la ruta guardada ya no existe")
}

// Busca el Supremo.exe en ubicaciones comunes
func GetRutasComunesSupremo() (string, error) {
	rutas := []string{
		`C:\Program Files\Supremo\Supremo.exe`,
		`C:\Program Files (x86)\Supremo\Supremo.exe`,
		`C:\Users\TIC8\Downloads\Supremo.exe`,
	}

	for _, path := range rutas {
		if _, err := os.Stat(path); err == nil {
			_ = os.WriteFile(GetRutaSupremoDependsModo(), []byte(path), 0644)
			return path, nil
		}
	}
	return "", errors.New("no se encontró Supremo en rutas comunes")
}

// Abre el diálogo nativo de Windows para buscar el Supremo.exe
func GetRutaSupremoExplorador() (string, error) {
	path, err := dialogSys.File().
		Title("Selecciona Supremo.exe").
		Filter("Ejecutables", "exe").
		Load()
	if err != nil {
		return "", err
	}
	return path, nil
}

// Guarda la ruta seleccionada
func WriteNewRutaSupremo(path string) {
	if path != "" {
		err := os.WriteFile(GetRutaSupremoDependsModo(), []byte(path), 0644)
		if err != nil {
			fmt.Println("Error al guardar la ruta:", err)
		}
	}
}

// recoge los codigos de las secciones de los terminales
func RellenarIdsSecciones(secciones model.Seccion) int {
	return secciones.Code
}

// metodo para hacer arrays de los endPoints dependiendo del terminal
func MakeArrayEndPoints(valor int, endPoint string) string {
	num := strconv.Itoa(valor)
	return ReemplazarSectionEnURL(endPoint, num)
}

/*
Con este metodo hacemos un array de las secciones que tenga el terminal
1. Si tenemos 3 secciones, las recorremos y cogemos sus indices
2. Con esos indices hacemos el array de la url base de subir y bajar turno
3. Simplemente modificamos el string cambiado el número de seccion por el correspondiente
*/
func ReemplazarSectionEnURL(url string, nuevoValor string) string {
	sectionPrefix := "section="
	start := strings.Index(url, sectionPrefix)
	if start == -1 {
		return url
	}

	end := strings.Index(url[start:], "&")
	if end == -1 {
		return url[:start+len(sectionPrefix)] + nuevoValor
	}

	end += start
	return url[:start+len(sectionPrefix)] + nuevoValor + url[end:]
}

// Metodo para mirar el limite de la botonera y evitar negativos
func CheckLimitBotonera(secciones []model.Seccion, posX string) (bool, int) {
	diferenciaConst := 300
	numero, err := strconv.Atoi(posX)
	if err != nil {
		return false, 0
	}
	ancho := screenshot.GetDisplayBounds(0).Dx()
	espacioNecesario := len(secciones) * diferenciaConst
	espacioRestante := ancho - espacioNecesario

	posXMax := espacioRestante
	if posXMax < 0 {
		posXMax = 0
	}

	valido := numero >= 0 && numero <= posXMax
	return valido, posXMax
}

// metodo principal para cargar los valores del fichero de configuracion en el modal
func LoadValuesFichero(path string) (*model.ConfigModelo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &model.ConfigModelo{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2) //comprueba que el formato no es diferente de (palabra=(palabra))
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		AuxLoadValuesFicheros(key, config, value)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// metodo auxiliar con el switch y los valores correspondientes
func AuxLoadValuesFicheros(key string, config *model.ConfigModelo, value string) {

	switch key {
	case "ip":
		config.Ip = value
	case "ipbalanza":
		config.IpBalanza = value
	case "token":
		config.Token = value
	case "usuario":
		config.Usuario = value
	case "password":
		config.Password = value
	case "botonera":
		config.Botonera = strings.ToLower(value) == "true"
	case "colorR":
		config.ColorRBotonera = ParseNullableInt(value)
	case "colorG":
		config.ColorGBotonera = ParseNullableInt(value)
	case "colorB":
		config.ColorBBotonera = ParseNullableInt(value)
	case "PosicionBotoneraX":
		config.PosicionXBotonera = ParseNullableInt(value)
	case "posicionBotoneraY":
		config.PosicionYBotonera = ParseNullableInt(value)
	case "visor":
		config.Visor = strings.ToLower(value) == "true"
	case "visorPosicionX":
		config.PosicionXVisor = ParseNullableInt(value)
	case "visorPosicionY":
		config.PosicionYVisor = ParseNullableInt(value)
	}

}

// Metodo que nos devuelve el turno actual en el que está la sección
func SetTurnoActual(labelTurno *widget.Label, seccion model.Seccion) {
	if model.GetTurno()[seccion.Name] == "" {
		labelTurno.Text = "000"
	} else {
		labelTurno.Text = model.GetTurno()[seccion.Name]
	}
}

// metodo que parsea un string a entero y comprueba si no es <nil>
func ParseNullableInt(valor string) *int {
	if valor == "<nil>" {
		return nil
	}
	numero, err := strconv.Atoi(valor)
	if err != nil {
		return nil
	}
	return &numero
}

// metodo para convertir valores enteros a string
func ValorAString(v interface{}) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok && s == "<nil>" {
		return ""
	}
	switch val := v.(type) {
	case int:
		return strconv.Itoa(val)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case string:
		return val
	default:
		return ""
	}
}

// metodo para setear los valores de la UI a través de objetos
func SetValuesUI(config model.ConfigModelo, configUI model.ConfigModeloUI) {
	configUI.Ip.SetText(ValorAString(config.Ip))
	configUI.IpBalanza.SetText(ValorAString(config.IpBalanza))
	configUI.Token.SetText(ValorAString(config.Token))
	configUI.Usuario.SetText(ValorAString(config.Usuario))
	configUI.Password.SetText(ValorAString(config.Password))
	configUI.Botonera.SetChecked(config.Botonera)
	configUI.ColorRBotonera.SetText(ValorAString(GetValueOrNil(config.ColorRBotonera)))
	configUI.ColorGBotonera.SetText(ValorAString(GetValueOrNil(config.ColorGBotonera)))
	configUI.ColorBBotonera.SetText(ValorAString(GetValueOrNil(config.ColorBBotonera)))
	configUI.PosicionXBotonera.SetText(ValorAString(GetValueOrNil(config.PosicionXBotonera)))
	configUI.PosicionYBotonera.SetText(ValorAString(GetValueOrNil(config.PosicionYBotonera)))
	configUI.Visor.SetChecked(config.Visor)
	configUI.PosicionXVisor.SetText(ValorAString(GetValueOrNil(config.PosicionXVisor)))
	configUI.PosicionYVisor.SetText(ValorAString(GetValueOrNil(config.PosicionYVisor)))
}

// metodo que obtiene la ventana usando librerias y apis
func GetHWND(title string) win.HWND {
	ret, _ := syscall.UTF16PtrFromString(title)
	return win.FindWindow(nil, ret)
}

// metodo que obtiene el color del pixel de la derecha de la UI de la botonera
func GetPixelColor(x, y int) (r, g, b byte) {
	hdc := win.GetDC(0)
	colorRef := win.GetPixel(hdc, int32(x), int32(y))
	win.ReleaseDC(0, hdc)
	r = byte(colorRef & 0xFF)
	g = byte((colorRef >> 8) & 0xFF)
	b = byte((colorRef >> 16) & 0xFF)
	return
}

// metodo de aproximacion para ocultar/mostrar ventana en base a colores
func IsColorMatch(r, g, b byte, valorR int, valorG int, valorB int) bool {
	baseR, baseG, baseB := valorR, valorG, valorB

	threshold := 35

	diff := AbsDiff(int(r), baseR) + AbsDiff(int(g), baseG) + AbsDiff(int(b), baseB)

	return diff <= threshold
}

// metodo que hace de aproximacion restando valores de los colores
func AbsDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

// metodo que obtiene la derecha, izquierda, arriba y abajo de la UI de la botonera
func GetWindowRect(hwnd win.HWND) (left, top, right, bottom int32) {
	var rect win.RECT
	if win.GetWindowRect(hwnd, &rect) {
		return rect.Left, rect.Top, rect.Right, rect.Bottom
	}
	return 0, 0, 0, 0
}

// Metodo para saber si está o no cubierta la ventana
func IsCovered(hwnd win.HWND) bool {
	exStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	return exStyle&win.WS_EX_TOPMOST != 0
}

// metodo que establece un contenedor y texto con sus propiedades
func NewLabeledEntry(labelText, placeholder string, entryWidth float32) (*widget.Entry, fyne.CanvasObject) {
	label := widget.NewLabel(labelText)
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	entryContainer := container.NewGridWrap(fyne.NewSize(entryWidth, 37), entry)
	return entry, container.NewHBox(label, layout.NewSpacer(), entryContainer)
}

// metodo para establecer la password del login
func PasswordLeveledEntry(labelText, placeholder string, entryWidth float32) (*widget.Entry, fyne.CanvasObject) {
	label := widget.NewLabel(labelText)
	entry := widget.NewPasswordEntry()
	entry.SetPlaceHolder(placeholder)
	entryContainer := container.NewGridWrap(fyne.NewSize(entryWidth, 37), entry)
	return entry, container.NewHBox(label, layout.NewSpacer(), entryContainer)
}

// metodo que establece un contenedor y texto con sus propiedades
func NewSizedEntry(placeholder string, width float32) (*widget.Entry, *fyne.Container) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	container := container.NewGridWrap(fyne.NewSize(width, 37), entry)
	return entry, container
}

// Habilita todos los widget.Entry del slice
func HabilitarEntries(entries []*widget.Entry) {
	for _, entry := range entries {
		entry.Enable()
	}
}

// Deshabilita todos los widget.Entry del slice
func DeshabilitarEntries(entries []*widget.Entry) {
	for _, entry := range entries {
		entry.Disable()
	}
}

// metodo que deshabilita/habilita botones según botones pulsados
func DisableEntrysByCheck(checkBox *widget.Check, entries []*widget.Entry) {
	if checkBox.Checked {
		HabilitarEntries(entries)
	} else {
		DeshabilitarEntries(entries)
	}
	checkBox.OnChanged = func(valor bool) {
		if valor {
			HabilitarEntries(entries)
		} else {
			DeshabilitarEntries(entries)
		}
	}
}

// metodo que comprueba si el valor del textField es un numero
func CheckIfPisitionsLoginIsNumber(textfield *widget.Entry) bool {
	valor := textfield.Text

	num, err := strconv.Atoi(valor)
	if err != nil {
		return false
	}
	return num >= 0
}

// Metodo para restringir el valor vertical de la botonera
func CheckRestrictionPositionY(textfield *widget.Entry) bool {
	valor := textfield.Text

	num, err := strconv.Atoi(valor)
	if err != nil {
		return false
	}
	return num >= 20 && num <= 32
}

// metodo que comprueba si los valores de un text son numeros y están dentro del rango de los colores (0 - 255)
func IsValidRGBLoginEntry(entry *widget.Entry) bool {
	valor := entry.Text

	num, err := strconv.Atoi(valor)
	if err != nil {
		return false
	}

	return num >= 0 && num <= 255
}

// metodo que obtiene valores de los textfield y los pasa a numeros
func GetOptionalIntValue(enabled bool, entry *widget.Entry) *int {

	if !enabled {
		return nil
	}

	val, err := strconv.Atoi(entry.Text)
	if err != nil {
		return nil
	}

	return &val
}

// metodo que transforma punteros que nos da registros de memoria en sus valores
func GetValueOrNil(puntero *int) interface{} {
	if puntero == nil {
		return nil
	}
	return *puntero
}

// Metodo para enviar los datos de configuración al fichero dibal.conf
func SendDatosConfig(configSave model.ConfigModeloSave) {
	auxColorR := GetOptionalIntValue(configSave.Botonera.Checked, configSave.ColorRBotonera)
	auxColorG := GetOptionalIntValue(configSave.Botonera.Checked, configSave.ColorGBotonera)
	auxColorB := GetOptionalIntValue(configSave.Botonera.Checked, configSave.ColorBBotonera)
	auxPosicionBotoneraVertical := GetOptionalIntValue(configSave.Botonera.Checked, configSave.PosicionYBotonera)
	auxPosicionBotoneraHorizontal := GetOptionalIntValue(configSave.Botonera.Checked, configSave.PosicionXBotonera)
	auxPosicionVisorVertical := GetOptionalIntValue(configSave.Visor.Checked, configSave.PosicionYVisor)
	auxPosicionVisorHorizontal := GetOptionalIntValue(configSave.Visor.Checked, configSave.PosicionXVisor)

	modelo := model.ConfigModelo{
		Ip:                configSave.Ip.Text,
		IpBalanza:         configSave.IpBalanza.Text,
		Token:             configSave.Token.Text,
		Usuario:           configSave.Usuario.Text,
		Password:          configSave.Password.Text,
		Botonera:          configSave.Botonera.Checked,
		ColorRBotonera:    auxColorR,
		ColorGBotonera:    auxColorG,
		ColorBBotonera:    auxColorB,
		PosicionXBotonera: auxPosicionBotoneraHorizontal,
		PosicionYBotonera: auxPosicionBotoneraVertical,
		Visor:             configSave.Visor.Checked,
		PosicionXVisor:    auxPosicionVisorHorizontal,
		PosicionYVisor:    auxPosicionVisorVertical,
	}

	SaveConfig(modelo, GetRutaDependsModo(), configSave.LoginWindow)
}

// metodo para guardar lo escrito en el modal de config, en un fichero
func SaveConfig(modelo model.ConfigModelo, file string, loginWindow fyne.Window) {
	content := fmt.Sprintf("ip=%s\nipbalanza=%s\ntoken=%s\nusuario=%s\npassword=%s\nbotonera=%t\n"+
		"colorR=%v\ncolorG=%v\ncolorB=%v\nPosicionBotoneraX=%v\nposicionBotoneraY=%v\nvisor=%t\n"+
		"visorPosicionX=%v\nvisorPosicionY=%v\n",
		modelo.Ip, modelo.IpBalanza, modelo.Token, modelo.Usuario, modelo.Password, modelo.Botonera,
		GetValueOrNil(modelo.ColorRBotonera), GetValueOrNil(modelo.ColorGBotonera),
		GetValueOrNil(modelo.ColorBBotonera), GetValueOrNil(modelo.PosicionXBotonera),
		GetValueOrNil(modelo.PosicionYBotonera), modelo.Visor,
		GetValueOrNil(modelo.PosicionXVisor), GetValueOrNil(modelo.PosicionYVisor))

	err := os.WriteFile(file, []byte(content), 0644)
	if err != nil {
		dialogFyne.ShowError(fmt.Errorf("no se pudo guardar la configuración: %v", err), loginWindow)
		return
	}

	dialogFyne.ShowInformation("Éxito", "Configuración guardada correctamente.", loginWindow)
	loginWindow.Close()
}

// metodo que comprueba si un texto está vacío
func EmptyText(textoCampos string) bool {
	return strings.TrimSpace(textoCampos) == ""
}

// metodo que comprueba valores del fichero de configuracion
func GetValorDesdeConf(path string, clave string) int {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return 0
	}
	file, err := os.Open(absPath)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(linea, clave) {
			valorStr := strings.TrimPrefix(linea, clave)
			valorStr = strings.TrimSpace(valorStr)
			valorInt, err := strconv.Atoi(valorStr)
			if err != nil {
				return 0
			}
			return valorInt
		}
	}
	return 0
}

// metodo que comprueba valores del fichero de configuracion sin enteros
func GetValueNotIntDesdeConf(path string, clave string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	file, err := os.Open(absPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(linea, clave) {
			valorStr := strings.TrimPrefix(linea, clave)
			return strings.TrimSpace(valorStr)
		}
	}
	return ""
}

// metodo para ocultar la botonera si establecemos como false la botonera en la configuracion
func OcultarBotoneraIfIsFalse(window fyne.Window, valor string) {
	if valor == "false" {
		window.Hide()
	} else {
		window.Show()
	}
}

// metodo que oculta la ventana en el for
func OcultarVentanaDelay(hwnd win.HWND) {
	go func(h win.HWND) {
		win.ShowWindow(h, win.SW_HIDE)
		win.RedrawWindow(h, nil, 0, win.RDW_ERASE|win.RDW_INVALIDATE|win.RDW_ALLCHILDREN|win.RDW_UPDATENOW)
	}(hwnd)
}

// metodo para obtener la posicion del visor establecido en el fichero de configuracion
func ObtenerPosicionVisor(posXStr string, posYStr string) (int32, int32) {
	x, errX := strconv.Atoi(posXStr)
	y, errY := strconv.Atoi(posYStr)

	if errX != nil || errY != nil {
		return 1014, 0
	}
	return int32(x), int32(y)
}

// metodo para obtener los valores de los colores del fichero de configuracion
func ObtenerColoresBotoneraOcultar(valor string) int {
	num, err := strconv.Atoi(valor)
	if err != nil {
		return 0
	}
	return int(num)
}

// Metodo para encontrar una ventana del sistema por el título
func FindWindow(title string) uintptr {
	t, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := model.ProcFindWindowW.Call(0, uintptr(unsafe.Pointer(t)))
	return hwnd
}

// Método para posicionar una ventana del sistema
func SetWindowPos(hwnd uintptr, x, y, w, h int32) {
	model.ProcSetWindowPos.Call(
		hwnd,
		0,
		uintptr(x),
		uintptr(y),
		uintptr(w),
		uintptr(h),
		model.SWP_NOZORDER|model.SWP_SHOWWINDOW,
	)
}

// Método que se encarga de levantar un webView
func AuxLevantarnavegador() {
	valorVisor := GetValueNotIntDesdeConf(GetRutaDependsModo(), "visor=")
	if valorVisor == "true" && len(os.Args) > 1 && os.Args[1] == "-webview" {
		ln, err := net.Listen("tcp", "127.0.0.1:9999")
		if err != nil {
			os.Exit(0)
		}
		defer ln.Close()

		title := "Información"
		go func() {
			maxIntentos := 20
			for i := 0; i < maxIntentos; i++ {
				time.Sleep(250 * time.Millisecond)
				hwnd := FindWindow(title)
				if hwnd != 0 {
					SetWindowPos(hwnd, 1025, 99, 1024, 680)
					return
				}
			}
		}()
		EjecutarWebView()
	}
}

// Metodo para ejecutar el lanzamiento del navegador
func EjecutarWebView() {
	url := os.Args[2]
	w := webview.New(false)
	if w == nil {
		return
	}
	defer w.Destroy()
	w.SetTitle("Información")
	w.SetSize(1024, 720, webview.HintNone)
	w.Navigate(url)
	w.Run()
	os.Exit(0)
}

// Metodo que ejecuta un proceso de levantar navegador
func Levantarnavegador(valorVisor string, posXVisor string, posYVisor string, urlGmedia string) {
	if valorVisor != "true" {
		return
	}

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err == nil {
		_ = conn.Close()
		return
	}

	cmd := exec.Command(os.Args[0], "-webview", urlGmedia)
	err = cmd.Start()
	if err == nil {
		model.NavegadorCmd = cmd
	}
}

// Metodo que cierra todas las ventanas de un array
func CerrarVentanas(ventanas []fyne.Window) {
	for _, ventana := range ventanas {
		if ventana != nil {
			ventana.Close()
		}
	}
}

// metodo que establece el valor de los labels de error al hacer llamadas de endPoints
func SetTextlabelError(label *widget.Label, mensaje string, segundos int, window fyne.Window, ruedaPersonal *widget.Button) {
	ruedaPersonal.Disable()
	label.SetText(mensaje)
	window.Resize(fyne.NewSize(210, 85))
	label.Show()
	label.Refresh()

	go func() {
		time.AfterFunc(time.Duration(segundos)*time.Second, func() {
			fyne.DoAndWait(func() {
				label.Hide()
				label.Refresh()
				ruedaPersonal.Enable()
				window.Resize(fyne.NewSize(210, 40))
			})
		})
	}()
}
