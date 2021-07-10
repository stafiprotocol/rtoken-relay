package config

type CallEnum uint8

var (
	Call         = CallEnum(0)
	DelegateCall = CallEnum(1)
)
