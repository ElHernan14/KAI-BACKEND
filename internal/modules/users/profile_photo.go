package users

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	profileUploadDir       = "uploads/profiles"
	defaultProfilePhotoURL = "/uploads/profiles/default-avatar.webp"
	maxProfilePhotoMB      = 5
)

var allowedProfilePhotoExts = []string{".jpg", ".jpeg", ".png", ".webp"}

func findProfilePhotoURL(userID uuid.UUID) *string {
	for _, ext := range allowedProfilePhotoExts {
		path := filepath.Join(profileUploadDir, userID.String()+ext)
		if _, err := os.Stat(path); err == nil {
			url := "/" + filepath.ToSlash(path)
			return &url
		}
	}

	defaultURL := defaultProfilePhotoURL
	return &defaultURL
}

func removeProfilePhotos(userID uuid.UUID) error {
	for _, ext := range allowedProfilePhotoExts {
		path := filepath.Join(profileUploadDir, userID.String()+ext)
		if err := removeFileIfExists(path); err != nil {
			return err
		}
	}

	oldPattern := filepath.Join(profileUploadDir, userID.String()+"_*")
	oldFiles, err := filepath.Glob(oldPattern)
	if err != nil {
		return err
	}

	for _, path := range oldFiles {
		if err := removeFileIfExists(path); err != nil {
			return err
		}
	}

	return nil
}

func removeFileIfExists(path string) error {
	err := os.Remove(path)
	if err == nil || os.IsNotExist(err) {
		return nil
	}

	return err
}

func normalizedProfilePhotoExt(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func isAllowedProfilePhotoExt(ext string) bool {
	for _, allowed := range allowedProfilePhotoExts {
		if ext == allowed {
			return true
		}
	}

	return false
}
