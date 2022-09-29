package core

import (
	"cc_DavidGayle_BackendAPI/internal/app/common"
	"cc_DavidGayle_BackendAPI/internal/app/model"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

type userHandler struct {
	config *model.Config
	db     common.DB
}

type UserHandler interface {
	SignUp(respWriter http.ResponseWriter, req *http.Request)
	Login(respWriter http.ResponseWriter, req *http.Request)
	GetAllUsers(respWriter http.ResponseWriter, req *http.Request)
	UpdateUser(respWriter http.ResponseWriter, req *http.Request)
}

func NewUserHandler(config *model.Config) UserHandler {
	return &userHandler{
		config: config,
		db:     common.NewDb(config),
	}
}

// SignUp handle the POST /signup request
func (h *userHandler) SignUp(respWriter http.ResponseWriter, req *http.Request) {
	var user, existingUser model.User
	if ok := getPayload(respWriter, req, &user); !ok {
		return
	}

	connection := h.db.GetDatabase()
	defer h.db.CloseDatabase(connection)

	connection.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.Email != "" {
		errorResponse(respWriter, req, "Email address already in use")
		return
	}

	connection.Create(&user)
	generateJwtResponse(respWriter, req, user.Email)
}

// Login handle the POST /login request
func (h *userHandler) Login(respWriter http.ResponseWriter, req *http.Request) {
	var login model.LoginRequest
	var existingUser model.User
	if ok := getPayload(respWriter, req, &login); !ok {
		return
	}

	connection := h.db.GetDatabase()
	defer h.db.CloseDatabase(connection)

	connection.Where("email = ?", login.Email).First(&existingUser)
	if existingUser.Email == "" || existingUser.Password != login.Password {
		errorResponse(respWriter, req, "Unable to login")
		return
	}

	generateJwtResponse(respWriter, req, login.Email)
}

// GetAllUsers handle the GET /users request
func (h *userHandler) GetAllUsers(respWriter http.ResponseWriter, req *http.Request) {
	connection := h.db.GetDatabase()
	defer h.db.CloseDatabase(connection)

	var users []model.User
	connection.Find(&users)

	var displayUsers []model.DisplayUser
	for _, user := range users {
		displayUser := model.DisplayUser{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
		displayUsers = append(displayUsers, displayUser)
	}

	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(displayUsers)
}

// UpdateUser handle the PUT /users request
func (h *userHandler) UpdateUser(respWriter http.ResponseWriter, req *http.Request) {
	var update model.UserNameUpdate
	if ok := getPayload(respWriter, req, &update); !ok {
		return
	}

	email, err := getEmailAddress(req)
	if err != nil {
		fmt.Println("Error while getting email address from JWT token")
		errorResponse(respWriter, req, "error processing update request (auth token)")
		return
	}

	connection := h.db.GetDatabase()
	defer h.db.CloseDatabase(connection)

	var existingUser model.User
	connection.Where("email = ?", email).First(&existingUser)

	if existingUser.Email == "" {
		errorResponse(respWriter, req, "error processing update request (database)")
		return
	}

	existingUser.FirstName = update.FirstName
	existingUser.LastName = update.LastName
	connection.Save(&existingUser)

	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(map[string]string{"response": "update succeeded"})
}

func getPayload(respWriter http.ResponseWriter, req *http.Request, target interface{}) bool {
	err := json.NewDecoder(req.Body).Decode(target)
	if err != nil {
		errorResponse(respWriter, req, "Error in reading payload.")
		return false
	}
	return true
}

func getEmailAddress(req *http.Request) (string, error) {
	authTokenKey := "X-Authentication-Token"
	if req.Header[authTokenKey] == nil {
		return "", fmt.Errorf("header not found")
	}

	var mySigningKey = []byte(common.SecretKey)
	token, err := jwt.Parse(
		req.Header[authTokenKey][0],
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("encountered an error while parsing token")
			}
			return mySigningKey, nil
		},
	)

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if email := claims["email"]; email == nil {
			return "", fmt.Errorf("email not in token")
		} else {
			return email.(string), nil
		}
	} else {
		return "", fmt.Errorf("error while getting claims from token")
	}
}

func errorResponse(respWriter http.ResponseWriter, req *http.Request, msg string) {
	err := common.SetError(msg)
	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(err)
}

func generateJwtResponse(respWriter http.ResponseWriter, req *http.Request, email string) {
	respWriter.Header().Set("Content-Type", "application/json")
	if token, err := generateJWT(email); err == nil {
		json.NewEncoder(respWriter).Encode(map[string]string{"jwt_token": token})
	} else {
		errorResponse(respWriter, req, "Error generating token")
	}
}

func generateJWT(email string) (string, error) {
	var mySigningKey = []byte(common.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Printf("Something went Wrong: %s\n", err.Error())
		return "", err
	}

	return tokenString, nil
}
