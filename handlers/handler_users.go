package handlers

import (
	"encoding/json"
	"fmt"
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

	fmt.Printf("Received: %+v\n", userRequest)

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Registration failed",
			"error": insertNewUser.Error,
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
