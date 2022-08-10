package telegram

import "tdgames/logic"

// Sender can send files to users
type Sender interface {
	Send(user []*logic.User, path, caption string)
}
