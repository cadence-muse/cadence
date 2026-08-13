package domain

type BandRole string

const (
	BandRoleOwner  BandRole = "owner"
	BandRoleMember BandRole = "member"
)

type BandMember struct {
	userID UserID
	role   BandRole
}

func LoadBandMember(userID UserID, role BandRole) BandMember {
	return BandMember{
		userID: userID,
		role:   role,
	}
}

func (m BandMember) UserID() UserID {
	return m.userID
}

func (m BandMember) Role() BandRole {
	return m.role
}
