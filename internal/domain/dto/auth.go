package dto

type AuthRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthRequestDTO) GetEmail() string {
	return r.Email
}

func (r *AuthRequestDTO) GetPassword() string {
	return r.Password
}
