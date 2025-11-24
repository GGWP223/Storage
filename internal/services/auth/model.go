package auth

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Uid       string `json:"uid"`
	CreatedAt string `json:"create_data"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}
