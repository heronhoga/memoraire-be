package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	if saveMemo.Error != nil || saveMemo.RowsAffected == 0 {
		if saveMemo.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_memos_date\" (SQLSTATE 23505)" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Memo already exists",
			})
			return
		}

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

func ReadMemo(w http.ResponseWriter, r *http.Request) {
	// Get query params
	totalItems := r.URL.Query().Get("items")
	if totalItems == "" {
		totalItems = "9"
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	totalItemsInt, err := strconv.Atoi(totalItems)
	if err != nil || totalItemsInt <= 0 {
		http.Error(w, "Invalid items value", http.StatusBadRequest)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Invalid page value", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	username, ok := r.Context().Value("user_id").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var existingUser models.User

	findUserId := config.DB.Where("username = ?", username).First(&existingUser)
	if findUserId.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User not available",
		})
		return
	}

	// Calculate offset
	offset := (pageInt - 1) * totalItemsInt

	// Fetch memos
	var memos []models.Memo
	result := config.DB.Where("user_id = ?", existingUser.ID).Limit(totalItemsInt).Offset(offset).Find(&memos)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve memos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Memo succesfully retrieved",
		"memo": memos,
	})
}

func UpdateMemo(w http.ResponseWriter, r *http.Request) {
	var newUpdateMemoRequest requests.RequestUpdateMemo

	err := json.NewDecoder(r.Body).Decode(&newUpdateMemoRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	//validate the request
	validate := validator.New()
	errValidate := validate.Struct(newUpdateMemoRequest)

	if errValidate != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Update memo failed",
		})
		return
	}

	// Get user ID from context
	username, ok := r.Context().Value("user_id").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var existingUser models.User

	findUserId := config.DB.Where("username = ?", username).First(&existingUser)
	if findUserId.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User not available",
		})
		return
	}

	//update the data
	updateMemo := config.DB.Exec(`
    UPDATE memos 
    SET date = ?, note = ? 
    WHERE memos.id = ?`,
    newUpdateMemoRequest.Date,
    newUpdateMemoRequest.Note,
    newUpdateMemoRequest.MemoId,
	)

	if updateMemo.Error != nil || updateMemo.RowsAffected == 0 {
		if updateMemo.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_memos_date\" (SQLSTATE 23505)" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Memo already exists",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Update memo failed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Update memo successful",
	})
}

func DeleteMemo(w http.ResponseWriter, r *http.Request) {
	var newDeleteMemoRequest requests.RequestDeleteMemo

	err := json.NewDecoder(r.Body).Decode(&newDeleteMemoRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	//validate the request
	validate := validator.New()
	errValidate := validate.Struct(newDeleteMemoRequest)

	if errValidate != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Delete memo failed",
		})
		return
	}

	// Get user ID from context
	username, ok := r.Context().Value("user_id").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var existingUser models.User

	findUserId := config.DB.Where("username = ?", username).First(&existingUser)
	if findUserId.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User not available",
		})
		return
	}

	deleteMemo := config.DB.Exec(`DELETE FROM memos WHERE memos.id = ? AND memos.user_id = ?`, newDeleteMemoRequest.MemoId, existingUser.ID)

	if deleteMemo.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Delete memo failed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Delete memo successful",
	})
}
