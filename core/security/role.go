package security

import "banking/core/security/perm"

const (
	RoleClient  = 1
	RoleAdmin   = 2
	RoleCashier = 3
)

var RolePermsMap = map[int]int{
	RoleClient:  perm.OpenAccount,
	RoleAdmin:   perm.OpenAccount | perm.CloseAccount | perm.WithdrawFunds,
	RoleCashier: perm.OpenAccount | perm.WithdrawFunds,
}
