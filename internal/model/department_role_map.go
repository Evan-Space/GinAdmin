package model

// DeptRoleMap 部门-角色关联

type DeptRoleMap struct {
	BaseModel
	DeptId uint `json:"dept_id" gorm:"column:dept_id;type:int unsigned;not null;default:0"`
	RoleId uint `json:"role_id" gorm:"column:role_id;type:int unsigned;not null;default:0"`
}

func (DeptRoleMap) TableName() string {
	return "department_role_map"
}