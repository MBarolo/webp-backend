package sessions

import "errors"

var ErrSessionFull = errors.New("the session does not allow for more players")
