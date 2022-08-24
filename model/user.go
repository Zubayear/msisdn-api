package model

type Role string

func (r Role) Values() (kinds []string) {
	for _, s := range []Role{Developer, ProductManager, Architect} {
		kinds = append(kinds, string(s))
	}
	return
}

const (
	Developer      Role = "Developer"
	ProductManager Role = "ProductManager"
	Architect      Role = "Architect"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Roles    Role   `json:"roles" binding:"required"`
}

func NewUser(username string, password string, roles Role) *User {
	return &User{Username: username, Password: password, Roles: roles}
}

type UserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserDto(username string, password string) *UserDto {
	return &UserDto{Username: username, Password: password}
}
