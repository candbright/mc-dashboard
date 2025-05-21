package domain

// AllowList 白名单列表
type AllowList []AllowListUser

// AllowListUser 白名单用户
type AllowListUser struct {
	XUid               string `json:"xuid"`
	Name               string `json:"name"`
	IgnoresPlayerLimit bool   `json:"ignores_player_limit"`
}
