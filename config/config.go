package config

type Config struct {
	RandomStringLength  int
	RandomStringCharset string
}

var DefaultConfig = Config{
	RandomStringLength:  10,
	RandomStringCharset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
}
