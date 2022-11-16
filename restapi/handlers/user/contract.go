package user

type ChangeNicknameRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}
