package security

import "banking/core/security/perm"

type Role uint

const (
	RoleClient  = Role(1)
	RoleAdmin   = Role(2)
	RoleCashier = Role(3)
)

var rolePermMap = map[Role]uint{
	RoleAdmin: perm.OpenAccount | perm.CloseAccount | perm.WithdrawFunds,
	RoleCashier: perm.OpenAccount | perm.WithdrawFunds,
}

func (role Role) Can(perm uint) bool {
	if perms, ok := rolePermMap[role]; !ok {
		return false
	} else {
		return perms&perm > 0
	}
}
