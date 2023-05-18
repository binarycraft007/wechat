package wechat

import "errors"

var ErrAlreadyLoggedOut = errors.New("already logged out")
var ErrUnknownFileType = errors.New("unknown file type")
var ErrContactListEmpty = errors.New("contact list empty")
