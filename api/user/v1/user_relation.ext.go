// Code Written Manually

package v1

// validUpdateAction is a list of valid actions that can be performed on a user relation
// {action: { curRelationStatus: newRelationStatus: true } }
var validUpdateAction = map[UpdateUserRelationAction]map[RelationStatus]RelationStatus{
	UpdateUserRelationAction_ACCEPT: {
		RelationStatus_REQUESTED: RelationStatus_FRIEND,
	},
	UpdateUserRelationAction_REJECT: {
		RelationStatus_REQUESTED: RelationStatus_STRANGER,
	},
	UpdateUserRelationAction_DELETE: {
		RelationStatus_FRIEND: RelationStatus_STRANGER,
	},
	UpdateUserRelationAction_BLOCK: {
		RelationStatus_FRIEND: RelationStatus_BLOCKED,
	},
	UpdateUserRelationAction_UNBLOCK: {
		RelationStatus_BLOCKED: RelationStatus_FRIEND,
	},
}

// CheckActionAndGetNewStatus checks if the action is valid and returns the new status
func (x UpdateUserRelationAction) CheckActionAndGetNewStatus(status RelationStatus) (RelationStatus, bool) {
	temp, ok := validUpdateAction[x]
	if !ok {
		return status, false
	}

	newStatus, ok := temp[status]
	return newStatus, ok
}
