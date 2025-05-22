package constants

const (
	MaxFileSize       = 10 * 1024 * 1024 // 10 MB
	MultipartFormFile = "images"
)

var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}
