package user

import (
	"crypto/sha256"
	"fmt"

	"xorm.io/xorm"
)

type UserService struct {
	session *xorm.Session
}

func NewUserService(session *xorm.Session) *UserService {
	return &UserService{session: session}
}

func (u *UserService) GetOne(id int64) (user *UserEntity, has bool, err error) {
	user = new(UserEntity)
	has, err = u.session.ID(id).Get(user)
	return
}

type UserPostForm struct {
	State      UserState `json:"state" xorm:"default 0 index"`                   //状态,default 1
	Passport   string    `json:"passport" xorm:"varchar(100) default '' unique"` //手机号码
	Password   string    `json:"password" xorm:"varchar(100) default ''"`        //
	Name       string    `json:"name"  xorm:"varchar(60) default ''"`            //姓名
	Sex        string    `json:"sex"  xorm:"varchar(10) default ''"`             //性别
	Avatar     string    `json:"avatar" xorm:"varchar(100) default ''"`          //头像
	RoleId     UserRole  `json:"role_id" xorm:"tinyint(2) default 0 index"`      //角色id
	ShopId     int64     `json:"shop_id" xorm:"default 0 index"`                 //所属商户id
	MerchantId int64     `json:"merchant_id" xorm:"default 0 index"`             //所属商户id
	WxId       int64     `json:"wx_id" xorm:"default 0 index"`                   //关联的顾客id

}

func (u *UserService) Insert(user *UserPostForm, byUser *UserEntity) (userErr error, err error) {
	if user.Passport == "" || user.Password == "" {
		userErr = fmt.Errorf("用户名密码不可为空")
		return
	}
	ue := new(UserEntity)
	ue.Passport = user.Passport
	ue.Password = generatePassword(user.Password)
	ue.State = UserStateOK
	ue.Name = user.Name
	ue.Sex = user.Sex
	ue.Avatar = user.Avatar
	ue.RoleId = user.RoleId
	ue.ShopId = user.ShopId
	ue.WxId = user.WxId
	_, err = u.session.Insert(ue)
	return
}

func (u *UserService) Delete(userId int64) (userErr error, err error) {
	if userId < 1 {
		userErr = fmt.Errorf("invalid userId")
		return
	}
	_, err = u.session.ID(userId).Update(&UserEntity{State: UserStateDel})
	return
}

type PatchPwd struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (u *UserService) UpdatePwd(uid int64, oldPwd, newPwd string) (err error) {
	user := new(UserEntity)
	var has bool
	if has, err = u.session.ID(uid).Get(user); err != nil {
		return
	} else {
		if !has {
			err = fmt.Errorf("user:%d not fount", uid)
			return
		}
	}
	if user.Password != generatePassword(oldPwd) {
		err = fmt.Errorf("old password error: store:%s,input:%s", user.Password, generatePassword(oldPwd))
		return
	}
	user.Password = generatePassword(newPwd)
	affectNum, err := u.session.ID(uid).Cols("password").Update(user)
	if err != nil {
		return
	}
	if affectNum != 1 {
		err = fmt.Errorf("user not found")
	}
	return

}
func (u *UserService) Login(passport, password string) (user *UserEntity, err error) {
	user = new(UserEntity)
	has, err := u.session.Where("passport=?", passport).And("password=?", generatePassword(password)).Get(user)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("no user found")
	}
	return
}

type UserPutForm struct {
	UserPostForm
	Id int64 `json:"id"`
}

func (u *UserService) Update(user *UserPutForm, updateCols ...string) (userErr error, err error) {
	if user.Id < 1 {
		userErr = fmt.Errorf("invalid id")
		return
	}
	session := u.session.ID(user.Id)
	if len(updateCols) > 0 {
		session.Cols(updateCols...)
	}
	ue := new(UserEntity)
	ue.Passport = user.Passport
	ue.Password = ""
	ue.State = user.State
	ue.Name = user.Name
	ue.Sex = user.Sex
	ue.Avatar = user.Avatar
	ue.RoleId = user.RoleId
	ue.ShopId = user.ShopId
	ue.MerchantId = user.MerchantId
	ue.WxId = user.WxId
	fmt.Printf("ue: %+v\n", *ue)
	_, err = session.Update(ue)
	return
}

type UserFilter struct {
	State UserState `form:"state"` //用户状态
	Rid   UserRole  `form:"rid"`   //角色id
	Name  string    `form:"name"`  //姓名
	Mid   int64     `form:"mid"`   //商户id
	Page  int       `form:"page"`  //请求页码
}

func (u *UserService) ListMP(filter *UserFilter, limit int) (result []*UserEntity, total int64, err error) {

	if filter.Mid > 0 {
		u.session.Where("merchant_id=?", filter.Mid)
	}
	if filter.Name != "" {
		u.session.Where("name like ?", fmt.Sprintf("%%%s%%", filter.Name))
	}
	if filter.Rid > 0 {
		u.session.Where("role_id=?", filter.Rid)
	}
	if filter.State != 0 {
		u.session.Where("state=?", filter.State)
	}

	if limit > 0 {
		u.session.Limit(limit, (filter.Page-1)*limit)
		total, err = u.session.FindAndCount(&result)
	} else {
		err = u.session.Find(&result)
	}
	return
}

func generatePassword(password string) string {
	return fmt.Sprintf("%X", sha256.Sum256([]byte(password)))
}
