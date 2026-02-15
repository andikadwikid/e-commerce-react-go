package structs

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required, min=6"`
	RoleIDs  []uint `json:"role_ids" binding:"required"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required, min=6"`
	RoleIDs  []uint `json:"role_ids" binding:"required"`
}

type UserDetailResponse struct {
	Id       uint           `json:"id"`
	Name     string         `json:"name"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Roles    []RoleResponse `json:"roles"`
}
