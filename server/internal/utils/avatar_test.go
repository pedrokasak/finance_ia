package utils

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAvatarDataURL(t *testing.T) {
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	valid := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngHeader)
	assert.NoError(t, ValidateAvatarDataURL(valid))

	assert.Error(t, ValidateAvatarDataURL("https://tracker.example/pixel.png"))
	assert.Error(t, ValidateAvatarDataURL("data:image/svg+xml;base64,PHN2Zz4="))
	assert.Error(t, ValidateAvatarDataURL("data:image/png;base64,not_base64"))

	invalidMagic := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte("hello"))
	assert.Error(t, ValidateAvatarDataURL(invalidMagic))
}

