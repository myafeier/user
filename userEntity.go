package user

import (
	"time"
)

type UserState int8

const (
	UserStateOK  UserState = 1  //正常
	UserStateDel UserState = -1 //停用
)

var UserStates = map[UserState]string{
	UserStateOK:  "正常",
	UserStateDel: "已停用",
}

func (u UserState) String() string {
	if name, ok := UserStates[u]; ok {
		return name
	} else {
		return "-"
	}
}

type UserEntity struct {
	Id         int64     `json:"id"`
	State      UserState `json:"state" xorm:"default 0 index"`                   //状态,default 1
	Passport   string    `json:"passport" xorm:"varchar(100) default '' unique"` //手机号码
	Password   string    `json:"-" xorm:"varchar(100) default ''"`               //
	Name       string    `json:"name"  xorm:"varchar(60) default ''"`            //姓名
	Sex        string    `json:"sex"  xorm:"varchar(10) default ''"`             //性别
	Avatar     string    `json:"avatar" xorm:"varchar(100) default ''"`          //头像
	RoleId     UserRole  `json:"role_id" xorm:"tinyint(2) default 0 index"`      //角色id
	MerchantId int64     `json:"merchant_id" xorm:"default 0 index"`             //所属商户id
	ShopId     int64     `json:"shop_id" xorm:"default 0 index"`                 //门店id
	WxId       int64     `json:"wx_id" xorm:"default 0 index"`                   //关联的顾客id
	Created    time.Time `json:"created" xorm:"created"`                         //
	Updated    time.Time `json:"updated" xorm:"updated"`                         //
	StateF     string    `json:"state_f,omitempty" xorm:"-"`                     //状态中文
	RoleF      string    `json:"role_f,omitempty" xorm:"-"`                      //角色中文
	ShopF      string    `json:"shop_f,omitempty" xorm:"-"`                      //门店名称
	MerchantF  string    `json:"merchant_f,omitempty" xorm:"-"`                  //商户名
	WxNickname string    `json:"wx_nickname,omitempty" xorm:"-"`                 //微信昵称
}

func (e *UserEntity) TableName() string {
	return "user"
}
func (e *UserEntity) Format() {
	e.RoleF = e.RoleId.String()
	e.StateF = e.State.String()
}
