// Code Written Manually

package v1

func (x *UserInternal) ToUser() *User {
	return &User{
		Uid:         x.GetUid(),
		Name:        x.GetName(),
		Email:       x.GetEmail(),
		Phone:       x.GetPhone(),
		Avatar:      x.GetAvatar(),
		AgentId:     x.AgentId,
		LoginStatus: x.GetLoginStatus(),
	}
}
