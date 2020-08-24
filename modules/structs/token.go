package structs

type Token struct {
	Token     string `json:"token"`
	ExpiresIn uint64 `json:"expires_in"`
}
