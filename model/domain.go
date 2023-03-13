package model

type User struct {
	ID     string
	Answer string
}

type AccessMode string

func (m AccessMode) Grantable() bool {
	return m == AccessModeRead || m == AccessModeReadWrite
}

func (m AccessMode) CanRead() bool {
	return m.CanWrite() || m == AccessModeRead
}

func (m AccessMode) CanWrite() bool {
	return m.CanInvite() || m == AccessModeReadWrite
}
func (m AccessMode) CanInvite() bool {
	return m == AccessModeOwner
}

const (
	AccessModeRead      AccessMode = "R"
	AccessModeReadWrite AccessMode = "RW"
	AccessModeOwner     AccessMode = "O"
)

type ACLEntry struct {
	User     string
	Mode     AccessMode
	Alias    string
	Inviter  string
	Accepted bool
}
