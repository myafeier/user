package user

import "time"

type RolePermissionEntity struct {
	Id           int64     `json:"id"`
	RoleId       int64     `json:"role_id" xorm:"default 0 index"` //角色id
	PermissionId int64     `json:"permissions" xorm:"default 0"`
	Created      time.Time `json:"created" xorm:"created"` //添加时间
	Updated      time.Time `json:"updated" xorm:"updated"` //添加时间
}
