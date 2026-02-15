package structs

import "backend-commerce/models"

type UserLoginResponse struct {
	Id          uint            `json:"id"`
	Name        string          `json:"name"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
	Roles       []string        `json:"roles"`
	Permissions map[string]bool `json:"permissions"`
}

func ToUserLoginResponse(user models.User, permissions map[string]bool) UserLoginResponse {
	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	return UserLoginResponse{
		Id:          user.Id,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		Roles:       roles,
		Permissions: permissions,
	}
}
