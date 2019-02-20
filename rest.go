package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Nota struct {
	// json:"titulo" esto es para que el json sepa que es en minuscula, y en la declaracion se pone el mayuscula Titulo
	titulo        string
	descripcion   string
	fechaCreacion time.Time
}

// Map de tipo Nota, que funciona como si se tratara de la base de datos
var mapNotas = make(map[string]Nota)
var id int

// getNotaHandler - GET

func getNotaHandler(w http.ResponseWriter, r *http.Request) {

	// Aqui declaramos una especie de objeto de tipo Nota
	var notas []Nota

	// Aqui se recorre el map que se declaro arriba mapNotas, con key, value
	// si no se utiliza en key o alguna variable se deja como _ y el no lo toma
	//y luego al arreglo le hacemos un append de los valores que hayan en el map
	for _, value := range mapNotas {
		notas = append(notas, value)
	}

	// Seteamos este header para que el navegador sepa que le estamos enviando algo en json
	w.Header().Set("Content-Type", "application/json")

	// la funcion marshal toma cualquier tipo de estructura go y la convierte a json, y se almacena en j
	j, err := json.Marshal(notas)

	// aqui se verifica si sale algun error al convertir la estructura a json, coon err,
	// el panic funciona para alertar al usuario de que ha ocurrido un fallo pero no mata el servidor,
	// y solo le sale a ese usuario
	if err != nil {
		panic(err)
	}

	// Aqui devolvemos la respuesta del servidor como un OK/200, ya que http tiene todos estos status
	w.WriteHeader(http.StatusOK)

	// Aqui retornamos el json que esta en j
	w.Write(j)
}

// postNotaHandler - POST
func postNotaHandler(w http.ResponseWriter, r *http.Request) {

	// Aqui declaramos una especie de objeto de tipo Nota
	var nota Nota

	// esta funcion de NewDecoder se encarga de tomalo que viene en el body
	// de la peticion del usuario y luego le dice que esa data la decodifique y la
	// almacena en el objeto notas de tipo Nota, con un puntero(&)

	// con err, es como se controlan los errores o se guarda si hay errores al momento
	// de que esta linea se ejecute, por si el usuario encia un json mal estructurado
	err := json.NewDecoder(r.Body).Decode(&nota)

	// Capturamos el error
	if err != nil {
		panic(err)
	}

	// esto es para setear la fecha de creacion en el momento que se hizo el post y fue exitoso
	// ya que el usuario no envia esta fecha si no que se setea cuando se realiza el post con exito
	nota.fechaCreacion = time.Now()
	id++
	// lo que hace este strconv es convertir un entero en un string con este Itoa
	k := strconv.Itoa(id)
	// luego al map en la posicion k, se le agrega la nota
	mapNotas[k] = nota

	// Seteamos este header para que el navegador sepa que le estamos enviando algo en json
	w.Header().Set("Content-Type", "application/json")

	// la funcion marshal toma cualquier tipo de estructura go y la convierte a json, y se almacena en j
	j, err := json.Marshal(nota)

	// aqui se verifica si sale algun error al convertir la estructura a json, coon err,
	// el panic funciona para alertar al usuario de que ha ocurrido un fallo pero no mata el servidor,
	// y solo le sale a ese usuario
	if err != nil {
		panic(err)
	}

	// Aqui devolvemos la respuesta del servidor como un 201, ya que http tiene todos estos status
	w.WriteHeader(http.StatusCreated)

	// Aqui retornamos el json que esta en j, y se devuleve la que inserto el usuario
	w.Write(j)

}

// putNotaHandler - PUT
func putNotaHandler(w http.ResponseWriter, r *http.Request) {

	// esto funciona para extraer los paramtros que esten viajando en la peticion, en este caso {id},
	// y le pasamos r que es rquest que es donde estan los parametros
	vars := mux.Vars(r)

	// Aqui almacenamos el parametro id que venia en la peticion por medio de vars
	k := vars["id"]

	// Aqui declaramos una especie de objeto de tipo Nota
	var notaUpdate Nota

	// se transforma el json en un objeto de go, de tipo Nota
	err := json.NewDecoder(r.Body).Decode(&notaUpdate)

	if err != nil {
		panic(err)
	}

	// esto funciona como un find, y lo que retorna es un boolean(true, false), indicando si encontro o no el id,
	// el ok sirve para que el nos retorne la nota si existe si no retorna vacio y esta nota queda almacenana en notaOld
	if notaOld, ok := mapNotas[k]; ok {
		// se guarda en la notaUpdate la fecha de creacion que tenia la notaOld
		notaUpdate.fechaCreacion = notaOld.fechaCreacion
		// luego se borra la nota
		delete(mapNotas, k)
		// y luego en la misma posicion insertamos la notaUpdate
		mapNotas[k] = notaUpdate
	} else {
		log.Printf("No Encontramos la nota que busca %s", k)
	}
	// Retorna un estado 204 de que no se encontraron datos
	w.WriteHeader(http.StatusNoContent)

}

// putNotaHandler - DELETE
func deleteNotaHandler(w http.ResponseWriter, r *http.Request) {

	// esto funciona para extraer los paramtros que esten viajando en la peticion, en este caso {id},
	// y le pasamos r que es rquest que es donde estan los parametros
	vars := mux.Vars(r)

	// Aqui almacenamos el parametro id que venia en la peticion por medio de vars
	k := vars["id"]

	// esto funciona como un find, y lo que retorna es un boolean(true, false), indicando si encontro o no el id,
	// el ok sirve para que el nos retorne la nota si existe si no retorna vacio y esta nota queda almacenana en notaOld
	if _, ok := mapNotas[k]; ok {
		// luego se borra la nota
		delete(mapNotas, k)
	} else {
		log.Printf("No Encontramos la nota que busca %s", k)
	}
	// Retorna un estado 204 de que no se encontraron datos
	w.WriteHeader(http.StatusNoContent)

}

func main() {

	// r es para crear las rutas o rutero
	r := mux.NewRouter().StrictSlash(false)

	// Api REST - (GET/POST/PUT/DELETE)
	r.HandleFunc("/api/notas", getNotaHandler).Methods("GET")
	r.HandleFunc("/api/notas", postNotaHandler).Methods("POST")
	r.HandleFunc("/api/notas/{id}", putNotaHandler).Methods("PUT")
	r.HandleFunc("/api/notas/{id}", deleteNotaHandler).Methods("DELETE")

	// creacion del server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Iniciar servidor
	log.Println("listening http://localhost:8080...")
	server.ListenAndServe()

}
