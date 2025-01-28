package dto

type UserRegisterReq struct {
	Email    string `json:"email,required"`
	Password string `json:"password,required"`
}

type UserRegisterResp struct{}

type UserLoginReq struct {
	Email    string `json:"email,required"`
	Password string `json:"password,required"`
}

type UserLoginResp struct {
	Token string `json:"token"`
}
