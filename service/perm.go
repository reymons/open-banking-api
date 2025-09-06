package service

import (
	"banking/core/security"
	"context"
)

type Perm interface {
	Can(ctx context.Context, role int, perms int) bool
}

type perm struct{}

// TODO:
// later, when we need to modify role permissions and
// add new roles at runtime, store perms in DB
func (s *perm) Can(ctx context.Context, role int, perms int) bool {
	rolePerms, ok := security.RolePermsMap[role]
	if !ok {
		return false
	}
	return perms&rolePerms > 0
}
