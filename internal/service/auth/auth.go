package auth

type Service struct {
	// TODO: add repo
}

type LoginDTO struct {
	Email    string
	Password string
	AppId    int
}

type RegisterDTO struct {
	Email    string
	Password string
}

type IsAdminDTO struct {
	UserId int
}
