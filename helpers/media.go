package helpers

import (
	"fmt"
	"galaxy/backend-api/config"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

// Helper to generate unique filename
func GenerateFilename(original string) string {
	fileExt := filepath.Ext(original)
	originalFileName := strings.TrimSuffix(filepath.Base(original), fileExt)
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	return filename
}

// Helper to save image, optionally resize
func SaveMedia(file io.Reader, filename string, resize bool, destDir string) error {
	if err := EnsureFolderExists(destDir); err != nil {
		return err
	}

	path := filepath.Join(destDir, filename)
	if resize {
		file, _, err := image.Decode(file)
		if err != nil {
			return err
		}
		src := imaging.Resize(file, 1000, 0, imaging.Lanczos)
		return imaging.Save(src, path)
	} else {
		out, err := os.Create(path)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		return err
	}
}

// Helper to construct image URL from baseAppUrl, folder, and filename
func GetMediaURL(folder, filename string) string {
	baseAppUrl := config.GetEnv("APP_URL", "http://localhost:3000")
	return fmt.Sprintf("%s/%s/%s", baseAppUrl, folder, filename)
}

// Helper to delete media file from hard-drive
func DeleteMedia(folder, filename string) error {
	path := filepath.Join(folder, filename)
	return os.Remove(path)
}

// Helper to ensure folder exists, create if not
func EnsureFolderExists(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return os.MkdirAll(folder, os.ModePerm)
	}
	return nil
}