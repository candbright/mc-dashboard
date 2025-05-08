package model

type AllowList []AllowListUser

type AllowListUser struct {
	XUid               string `json:"xuid"`
	Name               string `json:"name"`
	IgnoresPlayerLimit bool   `json:"ignoresPlayerLimit"`
}
