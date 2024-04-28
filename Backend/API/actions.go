package API

import (
	"encoding/json"
	"net/http"
	"paquetes/analizador"
	"paquetes/comandos"
	"paquetes/enviroment"
)

type data struct {
	Carnet int    `json:"carnet"`
	Nombre string `json:"nombre"`
}

type dataFront struct {
	Content string `json:"text"`
}

type respuesta struct {
	Consola  string                    `json:"consola"`
	Graficas []enviroment.ContentGraph `json:"graficas"`
}

type login struct {
	Id       string `json:"id"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type respuestaLogin struct {
	Message string `json:"message"`
	Consola string `json:"consola"`
}

func datos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var myData = data{201901103, "Josué Daniel Rojché García"}
	json.NewEncoder(w).Encode(myData)
}

func _dataFront(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	contenidoFront := dataFront{}

	decoder := json.NewDecoder(r.Body)

	if decoder.Decode(&contenidoFront) != nil {
		http.Error(rw, "Error en el JSON", http.StatusBadRequest)
		return
	}

	enviroment.ContentConsola = ""
	enviroment.ArrGraph = make([]enviroment.ContentGraph, 0)
	// fmt.Println(contenidoFront.Content)
	analizador.AnalizarPeticion(contenidoFront.Content)

	res := respuesta{}

	res.Consola = enviroment.ContentConsola
	res.Graficas = enviroment.ArrGraph

	json.NewEncoder(rw).Encode(res)
}

func doLogin(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	dataLogin := login{}

	decoder := json.NewDecoder(r.Body)

	if decoder.Decode(&dataLogin) != nil {
		http.Error(rw, "Error en el JSON", http.StatusBadRequest)
		return
	}

	enviroment.ContentLogin = ""
	enviroment.ContentConsola = ""
	comandos.Login(dataLogin.User, dataLogin.Password, dataLogin.Id)

	respuestaLogin := respuestaLogin{}
	respuestaLogin.Message = enviroment.ContentLogin
	respuestaLogin.Consola = enviroment.ContentConsola
	json.NewEncoder(rw).Encode(respuestaLogin)

	// fmt.Println(dataLogin.Id, dataLogin.User, dataLogin.Password)
}

func doLogout(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	enviroment.ContentLogin = ""
	enviroment.ContentConsola = ""

	comandos.Logout()

	respuestaLogin := respuestaLogin{}
	respuestaLogin.Message = enviroment.ContentLogin
	respuestaLogin.Consola = enviroment.ContentConsola

	json.NewEncoder(rw).Encode(respuestaLogin)
}

func doLoginFront(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	dataLogin := login{}

	decoder := json.NewDecoder(r.Body)

	if decoder.Decode(&dataLogin) != nil {
		http.Error(rw, "Error en el JSON", http.StatusBadRequest)
		return
	}

	enviroment.ContentLogin = ""
	enviroment.ContentConsola = ""
	comandos.LoginFront(dataLogin.User, dataLogin.Password, dataLogin.Id)

	respuestaLogin := respuestaLogin{}
	respuestaLogin.Message = enviroment.ContentLogin
	respuestaLogin.Consola = enviroment.ContentConsola
	json.NewEncoder(rw).Encode(respuestaLogin)

	// fmt.Println(dataLogin.Id, dataLogin.User, dataLogin.Password)
}
