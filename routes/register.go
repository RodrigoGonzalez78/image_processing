package routes

import (
	"encoding/json"
	"net/http"

	"github.com/RodrigoGonzalez78/db"
	"github.com/RodrigoGonzalez78/models"
	"github.com/RodrigoGonzalez78/utils"
)

// Register godoc
// @Summary      Registra un nuevo usuario
// @Description  Crea un nuevo usuario en la base de datos con nombre de usuario y contraseña.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body models.User true "Datos del usuario a registrar"
// @Success      201
// @Failure      400 {object} models.ErrorResponse
// @Router       /register [post]
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

	encrypt_password, err := utils.GenerateHashPassword(t.Password)

	if err != nil {
		http.Error(w, "Error al encriptar la contraseña!", 400)
		return
	}

	t.Password = encrypt_password

	encotrado, _ := db.IsUserNameUnique(t.UserName)

	if !encotrado {
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
