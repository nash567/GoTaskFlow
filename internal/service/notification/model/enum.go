package model

type Status int

const (
	StatusRead Status = iota
	StatusUnread
)

//go:generate enumer -type=Status -text -json  -trimprefix=Status -transform=snake -output=enum_status_gen.go
