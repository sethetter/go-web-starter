package middleware

const (
	SessKey_Authenticated = "authenticated"
	SessKey_User          = "user"
	SessKey_Csrf          = "csrf"
)

// func isAuthenticated(session sessions.Session) bool {
// 	authd := session.Get(SessKey_Authenticated)
// 	userId := session.Get(SessKey_User)
//
// 	if authd, ok := authd.(bool); ok && authd {
// 		if userId, ok := userId.(string); ok && userId != "" {
// 			return true
// 		}
// 	}
// 	return false
// }
