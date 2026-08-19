package domain

type BandMemberRole string

const (
	BandMemberRoleOwner BandMemberRole = "owner"
	BandRoleMember      BandMemberRole = "member"
)

type BandMember struct {
	userID UserID
	role   BandMemberRole
}

func LoadBandMember(userID UserID, role BandMemberRole) BandMember {
	return BandMember{
		userID: userID,
		role:   role,
	}
}

func (m BandMember) UserID() UserID {
	return m.userID
}

func (m BandMember) Role() BandMemberRole {
	return m.role
}
