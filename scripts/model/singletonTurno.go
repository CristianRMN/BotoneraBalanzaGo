package model

// Map con los turnos de las secciones de la aplicación
var turnosActuales map[string]string

// Metodo para iniciar los turnos
func InitTurnos() {
	turnosActuales = make(map[string]string)
}

// variable para saber si entraste en los settings
var entreEnSettings bool

// Setter y getter de que entre en los settings
func SetEntreSettings(valor bool) {
	entreEnSettings = valor
}

func GetEntreSettings() bool {
	return entreEnSettings
}

// Setter y getter de los turnos maps para las secciones
func SetTurno(seccionName string, turnoActual string) {
	turnosActuales[seccionName] = turnoActual
}

func GetTurno() map[string]string {
	return turnosActuales
}
