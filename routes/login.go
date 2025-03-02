package routes

import (
	"encoding/json"
	"net/http"

	"github.com/RodrigoGonzalez78/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var t models.User

	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		http.Error(w, "Usuario y/o contraseña invalidas"+err.Error(), 400)
		return
	}

	if len(t.UserName) == 0 {
		http.Error(w, "El username es requerido", 400)
		return
	}

	user, exist := bd.TriedLogin(t.Email, t.Password)

	if !exist {
		http.Error(w, "Usuario y/o contraseña invalidas", 400)
		return
	}

	jwtKey, err := jwt_p.GeneringJwt(document)

	if err != nil {
		http.Error(w, "Ocurrio un error"+err.Error(), 400)
		return
	}

	resp := models.ResponseLogin{
		UserId: document.ID,
		Token:  jwtKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
