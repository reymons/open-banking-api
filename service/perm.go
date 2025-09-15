package service

import (
	"banking/core"
	"context"
)

type Perm interface {
	Can(ctx context.Context, role int, perms int) bool
}

type permService struct{}

func NewPerm() Perm {
	return &permService{}
}

// TODO:
// Later, when we need to modify role permissions and
// add new roles at runtime, store perms in DB
func (s *permService) Can(ctx context.Context, role int, perms int) bool {
	rolePerms, ok := core.GetRolePerms(role)
	if !ok {
		return false
	}
	return (perms & rolePerms) == perms
}
