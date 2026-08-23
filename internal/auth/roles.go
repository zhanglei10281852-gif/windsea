package auth

import "github.com/zhanglei10281852-gif/windsea/internal/domain"

const (
	RoleOperator   = "operator"
	RoleSupervisor = "supervisor"
	RoleContractor = "contractor"
	RoleViewer     = "viewer"
)

func CanCreateWorkOrder(role string) bool { return role == RoleOperator || role == RoleSupervisor }
func CanPublishCampaign(role string) bool { return role == RoleSupervisor }
func CanAcceptHandoff(role string) bool   { return role == RoleContractor || role == RoleSupervisor }
func CanRead(role string) bool {
	return role == RoleOperator || role == RoleSupervisor || role == RoleContractor || role == RoleViewer
}
func Authorize(role, action string) error {
	allowed := map[string]bool{"create_work_order": CanCreateWorkOrder(role), "publish_campaign": CanPublishCampaign(role), "accept_handoff": CanAcceptHandoff(role), "read": CanRead(role)}
	if !allowed[action] {
		return domain.ErrForbidden
	}
	return nil
}
