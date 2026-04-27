package dto

import "go-bin/internal/entity"

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=120"`
	Email string `json:"email" binding:"required,email,max=160"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=120"`
	Email string `json:"email" binding:"required,email,max=160"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToUserResponse(user entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToUserResponses(users []entity.User) []UserResponse {
	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, ToUserResponse(user))
	}

	return response
}
