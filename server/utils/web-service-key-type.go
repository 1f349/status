package utils

type WebServiceKeyType int

const (
	KeyUser = WebServiceKeyType(iota)
	KeyState
	KeyAccessToken
	KeyRefreshToken
)
