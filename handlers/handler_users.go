package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/heronhoga/memoraire-be/config"
	"github.com/heronhoga/memoraire-be/models"
	"github.com/heronhoga/memoraire-be/requests"
	"github.com/heronhoga/memoraire-be/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var userRequest requests.RequestRegister

	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	//validate the request
	validate := validator.New()
	errValidate := validate.Struct(userRequest)

	if errValidate != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
		})
		return
	}

	//encrypt password
	hashedPassword, errHash := utils.EncryptPassword(userRequest.Password)
	if errHash != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
		})
		return
	}


	//insert new data
	var newUser = models.User{
		Username: userRequest.Username,
		Password: hashedPassword,
		FirstName: userRequest.FirstName,
		LastName: userRequest.LastName,
	}

	insertNewUser := config.DB.Create(&newUser)

	if insertNewUser.Error != nil || insertNewUser.RowsAffected == 0 {
		if insertNewUser.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_users_username\" (SQLSTATE 23505)" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Username already exists",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
		})
		return
	}
	

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Registration successful",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest requests.RequestLogin

	err := json.NewDecoder(r.Body).Decode(&userLoginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	//validate the request
	validate := validator.New()
	errValidate := validate.Struct(userLoginRequest)

	if errValidate != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Login failed",
		})
		return
	}

}
