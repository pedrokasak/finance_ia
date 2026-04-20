package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

const (
	MaxAvatarBytes = 256 * 1024
)

// ValidateAvatarDataURL accepts only inline data URLs for png/jpeg/webp.
// This prevents external image trackers and validates payload by magic bytes.
func ValidateAvatarDataURL(v string) error {
	if v == "" {
		return nil
	}
	if len(v) > 500000 {
		return errors.New("avatar too large")
	}

	prefixes := []string{
		"data:image/png;base64,",
		"data:image/jpeg;base64,",
		"data:image/jpg;base64,",
		"data:image/webp;base64,",
	}

	var payload string
	validPrefix := false
	for _, p := range prefixes {
		if strings.HasPrefix(v, p) {
			payload = strings.TrimPrefix(v, p)
			validPrefix = true
			break
		}
	}
	if !validPrefix {
		return errors.New("invalid avatar format")
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return errors.New("invalid avatar base64")
	}
	if len(decoded) == 0 || len(decoded) > MaxAvatarBytes {
		return errors.New("avatar size is invalid")
	}
	if !hasAllowedMagicBytes(decoded) {
		return errors.New("unsupported avatar type")
	}

	return nil
}

func hasAllowedMagicBytes(b []byte) bool {
	// PNG: 89 50 4E 47 0D 0A 1A 0A
	if len(b) >= 8 &&
		b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
		b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A {
		return true
	}
	// JPEG: FF D8 FF
	if len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF {
		return true
	}
	// WEBP: RIFF....WEBP
	if len(b) >= 12 &&
		b[0] == 'R' && b[1] == 'I' && b[2] == 'F' && b[3] == 'F' &&
		b[8] == 'W' && b[9] == 'E' && b[10] == 'B' && b[11] == 'P' {
		return true
	}
	return false
}

