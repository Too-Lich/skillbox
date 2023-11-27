package user

type User struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

func NewUser(name string, age int) *User {
	return &User{
		Name:    name,
		Age:     age,
		Friends: make([]string, 0),
	}
}

func UpdateUser(name string, age int, friends []string) *User {
	return &User{
		Name:    name,
		Age:     age,
		Friends: friends,
	}
}

func NewFriend(name string, friends []string) *User {
	return &User{
		Name:    name,
		Friends: friends,
	}
}
