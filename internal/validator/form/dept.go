package form

// CreateDepForm 新增部门
type CreateDeptForm struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Pid         uint   `json:"pid"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Sort        uint   `json:"sort"`
}



// UpdateDeptForm 更新部门
type UpdateDeptForm struct {
	ID uint `json:"id" binding:"required,gt=0"`
	Name string `json:"name" binding:"omitempty,min=1,max=60"`
	Pid *uint `json:"pid"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Sort *uint `json:"sort"`
}

// DeptListQuery 部门列表查询
type DeptListQuery struct {
	Paginate
	Name string `form:"name" json:"name"`
	Pid *uint `form:"pid" json:"pid"`
}

func NewDeptListQuery() *DeptListQuery {
	return &DeptListQuery{Paginate: Paginate{Page: 1, PerPage: 100}}
}

// DeptBindRoleForm 部门绑定角色
type DeptBindRoleForm struct {
	DeptId uint `json:"dept_id" binding:"required"`
	RoleIds []uint `json:"role_ids" binding:"required"`
}