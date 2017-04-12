package msgrelay

import "testing"

func TestDistributor_AddUser(t *testing.T) {
	d := NewDistributor()
	u := d.AddUser("test_uid")
	u.AddSession("test_sid1")
	u.AddSession("test_sid2")
	u.AddSession("test_sid3")
	m := u.NewMsg("test_uid", "Hello, this is a test")
	u.ReceiveMsg(m)
}
