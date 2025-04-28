package rbac

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/casbin/casbin/v2"
)

type permManager struct {
	enforcer *casbin.Enforcer
}

func (pm *permManager) Enforcer() *casbin.Enforcer {
	return pm.enforcer
}

func NewPermManager(enforcer *casbin.Enforcer) *permManager {
	return &permManager{
		enforcer: enforcer,
	}
}

func identity[T defs.RBACObj | defs.RBACSubject | defs.RBACDomain, U string | int | uint | int64](rType T, objID U) string {
	return string(rType) + ":" + fmt.Sprint(objID)
}

func (pm *permManager) GrantGroupPermission(groupID string, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error) {
	groupSubject := identity(defs.RBACSubjectGroup, groupID)
	objSubject := identity(objType, objID)
	domain := identity(defs.RBACDomainTenant, tenantID)

	return pm.enforcer.AddPolicy(groupSubject, objSubject, string(action), domain)
}

func (pm *permManager) RevokeGroupPermission(groupID string, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error) {
	groupSubject := identity(defs.RBACSubjectGroup, groupID)
	objSubject := identity(objType, objID)
	domain := identity(defs.RBACDomainTenant, tenantID)

	return pm.enforcer.RemovePolicy(groupSubject, objSubject, string(action), domain)
}

func (pm *permManager) GrantUserPermission(userID int, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error) {
	userSubject := identity(defs.RBACSubjectUser, userID)
	objSubject := identity(objType, objID)
	domain := identity(defs.RBACDomainTenant, tenantID)

	return pm.enforcer.AddPolicy(userSubject, objSubject, string(action), domain)
}

func (pm *permManager) RevokeUserPermission(userID int, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error) {
	userSubject := identity(defs.RBACSubjectUser, userID)
	objSubject := identity(objType, objID)
	domain := identity(defs.RBACDomainTenant, tenantID)

	return pm.enforcer.RemovePolicy(userSubject, objSubject, string(action), domain)
}

func (pm *permManager) CheckPermission(userID int, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error) {
	userSubject := identity(defs.RBACSubjectUser, userID)
	objSubject := identity(objType, objID)
	domain := identity(defs.RBACDomainTenant, tenantID)

	return pm.enforcer.Enforce(userSubject, objSubject, string(action), domain)
}

func (pm *permManager) AddUserToGroup(userID int, groupID string, tenantID int) (bool, error) {
	userSub := identity(defs.RBACSubjectUser, userID)
	groupSub := identity(defs.RBACSubjectGroup, groupID)
	domain := identity(defs.RBACDomainTenant, tenantID)

	return pm.enforcer.AddGroupingPolicy(userSub, groupSub, domain)
}

func (pm *permManager) RemoveUserFromGroup(userID int, groupID string, tenantID int) (bool, error) {
	userSub := identity(defs.RBACSubjectUser, userID)
	groupSub := identity(defs.RBACSubjectGroup, groupID)
	domain := identity(defs.RBACDomainTenant, tenantID)

	return pm.enforcer.RemoveGroupingPolicy(userSub, groupSub, domain)
}
