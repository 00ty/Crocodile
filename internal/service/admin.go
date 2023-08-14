package service

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// Login 管理员登录
func (ctr *AdminService) Login() {}

// List 卡密列表
func (ctr *AdminService) List() {}

// DelKami 删除卡密
func (ctr *AdminService) DelKami() {}

// Search 搜索规则
func (ctr *AdminService) Search() {}

// Generate 生成卡密
func (ctr *AdminService) Generate() {}
