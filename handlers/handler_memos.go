package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/heronhoga/memoraire-be/config"
	"github.com/heronhoga/memoraire-be/models"
	"github.com/heronhoga/memoraire-be/requests"
)

func CreateMemo(w http.ResponseWriter, r *http.Request) {
	var memoRequest requests.RequestCreateMemo

	err := json.NewDecoder(r.Body).Decode(&memoRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	//validate the request
	validate := validator.New()
	errValidate := validate.Struct(memoRequest)

	if errValidate != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Create memo failed",
		})
		return
	}

	//find user id
	userId := r.Context().Value("user_id").(string)
	var existingUser models.User

	findUser := config.DB.Where("username = ?", userId).First(&existingUser)

	if findUser.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "No user available",
		})
		return
	}

	//save memo
	newMemo := models.Memo{
		Date: memoRequest.Date,
		Note: memoRequest.Note,
		UserID: existingUser.ID,
	}

	saveMemo := config.DB.Create(&newMemo)
	if saveMemo.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "New Memo successfully created",
	})

}