package utils

import "errors"

var ErrUserNotFoundInProducer = errors.New("user not found")

var ErrNoMessagesFound = errors.New("no messages found")
