package common

var SecretKey = "2022.09.28_secret_key"

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

func SetError(message string) Error {
	return Error{
		IsError: true,
		Message: message,
	}
}
