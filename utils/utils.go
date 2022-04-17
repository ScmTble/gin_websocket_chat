package utils

import "strconv"

func VerifyToken(token string) (uint, bool) {
	uid, err := strconv.ParseUint(token, 10, 0)
	if err != nil {
		return 0, false
	}
	return uint(uid), true
}
