package form

// CreateMenuForm 新增菜单
type CreateMenuForm struct {
	Code        string `json:"code" binding:"required,min=1,max=120"`
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Icon        string `json:"icon" binding:"omitempty,max=255"`
	Path        string `json:"path" binding:"omitempty,max=255"`
	Component   string `json:"component" binding:"omitempty,max=255"`
	Type        uint8  `json:"type" binding:"required,oneof=1 2 3"`
	Pid         uint   `json:"pid"`
	Sort        uint   `json:"sort"`
	IsShow      *uint8 `json:"is_show" binding:"omitempty,oneof=0 1"`
	Status      *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
	IsAuth      *uint8 `json:"is_auth" binding:"omitempty,oneof=0 1"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// UpdateMenuForm 更新菜单
type UpdateMenuForm struct {
	ID          uint   `json:"id" binding:"required"`
	Code        string `json:"code" binding:"omitempty,min=1,max=120"`
	Name        string `json:"name" binding:"omitempty,min=1,max=255"`
	Icon        string `json:"icon" binding:"omitempty,max=255"`
	Path        string `json:"path" binding:"omitempty,max=255"`
	Component   string `json:"component" binding:"omitempty,max=255"`
	Type        *uint8 `json:"type" binding:"omitempty,oneof=1 2 3"`
	Pid         *uint  `json:"pid"`
	Sort        *uint  `json:"sort"`
	IsShow      *uint8 `json:"is_show" binding:"omitempty,oneof=0 1"`
	Status      *uint8 `json:"status" binding:"omitempty,oneof=0 1"`
	IsAuth      *uint8 `json:"is_auth" binding:"omitempty,oneof=0 1"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// MenuListQuery 菜单列表查询
type MenuListQuery struct {
	Paginate
	Keyword string `form:"keyword" json:"keyword"`
	Status  *uint8 `form:"status" json:"status" binding:"omitempty,oneof=0 1"`
}

// NewMenuListQuery 创建菜单列表查询
func NewMenuListQuery() *MenuListQuery {
	return &MenuListQuery{Paginate: Paginate{CurrentPage: 1, PageSize: 10}}
}
