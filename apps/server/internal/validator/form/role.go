package form

type CreateRoleForm struct {
	Code        string `json:"code" binding:"required,min=1,max=60"`
	Name        string `json:"name" binding:"required,min=1,max=60"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Pid         uint   `json:"pid"`
	Sort        uint   `json:"sort"`
	Status      *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
}

// UpdateRoleForm 更新角色
type UpdateRoleForm struct {
	ID          uint   `json:"id" binding:"required"`
	Code        string `json:"code" binding:"omitempty,min=1,max=60"`
	Name        string `json:"name" binding:"omitempty,min=1,max=60"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Pid         *uint  `json:"pid"`
	Sort        *uint  `json:"sort"`
	Status      *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
}

// RoleListQuery 角色列表查询
type RoleListQuery struct {
	Paginate
	Name   string `form:"name" json:"name"`
	Status *uint8 `form:"status" json:"status"`
	Pid    *uint  `form:"pid" json:"pid"`
}

// NewRoleListQuery 创建角色列表查询
func NewRoleListQuery() *RoleListQuery {
	return &RoleListQuery{Paginate: Paginate{CurrentPage: 1, PageSize: 10}}
}
