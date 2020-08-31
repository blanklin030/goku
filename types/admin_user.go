package types

const COOKIE_TOKEN = "GOKU_SSO_TOKEN"
const COOKIE_USER = "GOKU_SSO_USER"
const COOKIE_REQUEST = "GOKU_REQUEST_ID"

type AdminLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
