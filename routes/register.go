package routes

import (
	"encoding/json"
	"net/http"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var t models.User

	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		http.Error(w, "Error en los datos recibidos: "+err.Error(), 400)
		return
	}

	if len(t.UserName) == 0 {
		http.Error(w, "El nombre de usuario es requerido!", 400)
		return
	}

	if len(t.Password) < 8 {
		http.Error(w, "La contaseña debe tener almenos 8 caracteres!", 400)
		return
	}

	encotrado, _ := db.IsUserNameUnique(t.UserName)

	if encotrado {
		http.Error(w, "Ya esta registrado el nombre usuario!", 400)
		return
	}

	err = db.CreateUser(t)

	if err != nil {
		http.Error(w, "No se pudo registrar el usuario: "+err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
