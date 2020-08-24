package models

type LoginType int

const (
	LoginNoType LoginType = iota
	LoginPlain
	LoginOAuth2
)

func LoginTypeToString(t LoginType) string {
	switch t {
	case LoginNoType:
		return "none"
	case LoginPlain:
		return "plain"
	case LoginOAuth2:
		return "oauth2"
	}
	return "none"
}
