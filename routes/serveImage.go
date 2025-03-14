package routes

import (
	"net/http"
	"os"
	"path/filepath"
)

/*
Utiliso ServeImage para servir imagenes desde el servidor
sin necesidad de exponer la ruta real del archivo.
*/
func ServeImage(w http.ResponseWriter, r *http.Request) {

	filePath := filepath.Join("uploads", r.URL.Path[len("/images/"):])

	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "Archivo o directorio no encontrado", http.StatusNotFound)
		return
	}
	if info.IsDir() {
		http.Error(w, "Acceso prohibido", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filePath)
}
