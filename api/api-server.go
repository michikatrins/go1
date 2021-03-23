// main.go
package main
import (
    "fmt"
    "log"
    "net/http"
)
// IndexHandler nos permite manejar la petición a la ruta '/'
// y retornar "hola mundo" como respuesta al cliente.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hola mundo")
}
func main() {
    // Instancia de http.DefaultServerMux
    mux := http.NewServeMux()
    // Ruta a manejar
    mux.HandleFunc("/", IndexHandler)
    // Servidor escuchando en el puerto 8080
    http.ListenAndServe(":8080", mux)
}