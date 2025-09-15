package core

const (
	RoleClient = 1 + iota
	RoleAdmin
	RoleCashier
)

var rolePerms = map[int]int{
	RoleClient:  PermRequestAccount | PermViewAccount,
	RoleAdmin:   0,
	RoleCashier: 0,
}

func GetRolePerms(role int) (int, bool) {
    v, ok := rolePerms[role]
    return v, ok
}
