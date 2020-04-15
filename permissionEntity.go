package user

type PermissionEntity struct {
	Id       int64  `json:"id"`
	Name     string `xorm:"varchar(200)"              json:"name"`     // 权限名称
	Url      string `xorm:"varchar(200)"               json:"url"`     // 权限url
	ParentId int64  `xorm:"default 0 index "         json:"parent_id"` // 父级id
	Sort     uint   `xorm:"default 0"              json:"sort"`        // 排序(数字越小越靠前)
	Created  uint   `xorm:"created"       json:"created"`              // 创建时间
	Updated  uint   `xorm:"updated"       json:"updated"`              // 更新时间
}
