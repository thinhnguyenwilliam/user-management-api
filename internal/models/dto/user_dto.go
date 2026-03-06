package dto

type CreateUserRequest struct {
	Name        string
	Email       string
	Password    string
	PhoneNumber string
}

type UserResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
