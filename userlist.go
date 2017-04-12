package msgrelay

type UserList struct {
	ids []UserID
	c   chan *User
}

func NewUserList(ids ...UserID) (users *UserList) {
	return &UserList{
		ids: ids,
		c:   make(chan *User),
	}
}
