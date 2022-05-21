// Code Written Manually

package v1

// updateStatusMap is a list of valid actions that can be performed on a user relation
// { { curFriendStatus: {newFriendStatus: true } }
var updateStatusMap = map[FriendStatus]map[FriendStatus]bool{
	FriendStatus_FRIEND: {
		FriendStatus_STRANGER: true,
		FriendStatus_BLOCKED:  true,
	},
	FriendStatus_STRANGER: {}, // you can't be change status when you are stranger
	FriendStatus_BLOCKED: {
		FriendStatus_UNBLOCKED: true, // you can unblock when you are blocked
	},
}

// CanUpdateStatus checks if a friend status can be updated to another status
func (x FriendStatus) CanUpdateStatus(targetStatus FriendStatus) bool {
	temp, ok := updateStatusMap[x]
	if !ok {
		return false
	}

	return temp[targetStatus]
}
