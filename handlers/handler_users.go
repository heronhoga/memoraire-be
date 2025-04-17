package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/heronhoga/memoraire-be/config"
	"github.com/heronhoga/memoraire-be/models"
	"github.com/heronhoga/memoraire-be/requests"
	"github.com/heronhoga/memoraire-be/utils"
	"golang.org/x/crypto/bcrypt"
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

	//find existing user
	var existingUser models.User

	findUsername := config.DB.Where("username = ?", userLoginRequest.Username).First(&existingUser)
	if findUsername.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Wrong username or password",
		})
		return
	}

	//match the password
	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userLoginRequest.Password))
	if errMatchPassword != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Wrong username or password",
		})
		return
	}

	//generate token
	jwtToken, errGenerateJwt := utils.GenerateJWT(existingUser.Username)
	if errGenerateJwt != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
		})
		return
	}

	//encrypt token
	encryptedToken, errEncryptToken := utils.Encrypt(jwtToken)
	if errEncryptToken != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
		})
		return
	}

	//store encrypted token to database
	existingUser.Session = encryptedToken
	updateUserToken := config.DB.Where("username = ?", existingUser.Username).Save(&existingUser)

	if updateUserToken.Error != nil || updateUserToken.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
		})
		return
	}

	//return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token": encryptedToken,
	})
}


func Logout(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.Context().Value("user_id").(string))

	//remove token from database
	userId := r.Context().Value("user_id").(string)
	removeToken := config.DB.Model(&models.User{}).Where("username = ?", userId).Update("session", "")


	if removeToken.RowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Logout failed",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}