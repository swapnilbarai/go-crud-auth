package utils

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var mobileNumberRegex = regexp.MustCompile(`^[1-9][0-9]{9}$`)
var AccessSecret = []byte("Swapnil")
var RefreshSecret = []byte("barai")

const NanoSeconds = 1000000000

var RefreshTokenDuration = int64(NanoSeconds * 2592000) //month
var AccessTokenDuration = int64(NanoSeconds * 3600)     //hour

const RefrehTokenType = "refresh"
const AccessTokenType = "access"
