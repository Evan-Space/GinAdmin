package form

// LoginForm 登陆表单
type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=16"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=18"`
	Nickname string `json:"nickname" binding:"required,min=1,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=120"`
	Phone    string `json:"phone_number" binding:"omitempty,max=15"`
	Status   *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
}

// CreateAdminUser 新增用户
type CreateAdminUser struct {
	Username string `json:"username" binding:"required,min=2,max=30"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname" binding:"required,min=1,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=120"`
	Phone    string `json:"phone_number" binding:"omitempty,max=15"`
	Status   *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
}


// UpdateAdminUser 更新用户（ID 必填，其余字段可选）
type UpdateAdminUser struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"omitempty,min=2,max=30"`
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=120"`
	Phone    string `json:"phone_number" binding:"omitempty,max=15"`
	Status   *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
}

// UpdateProfile 更新个人资料（不需要 ID，从 Token 获取）
type UpdateProfile struct {
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=30"`
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
	Email    string `json:"email" binding:"omitempty,email,max=120"`
	Phone    string `json:"phone_number" binding:"omitempty,max=15"`
}

// AdminUserList 用户列表查询
type AdminUserList struct {
	Paginate
	Username string `form:"username" json:"username"`
	Nickname string `form:"nickname" json:"nickname"`
	Status   *uint8 `form:"status" json:"status"`
	Email    string `form:"email" json:"email"`
}

// BindRoleForm 为用户绑定角色
type BindRoleForm struct {
	UserID  uint   `json:"user_id" binding:"required"`
	RoleIDs []uint `json:"role_ids" binding:"required"`
}