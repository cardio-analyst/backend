package http

import "regexp"

const (
	passwordPattern           = `"password":.*".+"`
	passwordReplacementString = `"password": "******"`
)

var rePassword = regexp.MustCompile(passwordPattern)

// hidePassword hides password from string.
func hidePassword(stringWithPassword string) string {
	return rePassword.ReplaceAllString(stringWithPassword, passwordReplacementString)
}
