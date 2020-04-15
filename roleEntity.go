package user

type UserRole int8

const (
	UserRoleAdmin  UserRole = 1 //管理员
	UserRoleWorker UserRole = 2 //业务操作员
)

var UserRoles = map[UserRole]string{
	UserRoleAdmin:  "管理员",
	UserRoleWorker: "操作员",
}

func (u UserRole) String() string {
	if name, ok := UserRoles[u]; ok {
		return name
	} else {
		return "-"
	}
}
