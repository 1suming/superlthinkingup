package checker

import (
	"encoding/json"

	"github.com/answerdev/answer/configs"
)

var (
	reservedUsernameMapping = make(map[string]bool)
)

func init() {
	var usernames []string
	_ = json.Unmarshal(configs.ReservedUsernames, &usernames)
	for _, username := range usernames {
		reservedUsernameMapping[username] = true
	}
}

// IsReservedUsername checks whether the username is reserved
func IsReservedUsername(username string) bool {
	return reservedUsernameMapping[username]
}
