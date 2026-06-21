package form

// 分页参数查询
type Paginate struct {
	Page    int `form:"page" json:"page" binding:"omitempty,gt=0"`
	PerPage int `form:"per_page" json:"per_page" binding:"omitempty,gt=0"`
}

// NewPaginate 创建一个新的分页查询
func NewPaginate() *Paginate {
	return &Paginate{}
}


// IDForm ID 查询/删除参数
type IDForm struct {
	ID uint `form:"id" json:"id" binding:"required"`
}

// NewIdForm 创建 ID 表单
func NewIdForm() *IDForm {
	return &IDForm{}
}