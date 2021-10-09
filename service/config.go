package service

type Config struct {
	BaseURL   string
	HexSecret string
}

var (
	defaultSecret = []byte{18, 232, 139, 12, 216, 189, 22, 128, 122, 49, 246, 137, 191, 24, 38, 210}
)
