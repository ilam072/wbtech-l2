package storage

import (
	"github.com/ilam072/wbtech-l2/16/pkg/errutils"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func LocalPath(baseDir string, u *url.URL) string {
	// Убираем начальный слэш у пути
	cleanPath := u.Path
	if len(cleanPath) > 0 && cleanPath[0] == '/' {
		cleanPath = cleanPath[1:]
	}

	// Если путь пустой, используем index.html
	if cleanPath == "" {
		cleanPath = "index.html"
	}

	filePath := filepath.Join(baseDir, u.Host, cleanPath)

	// Если нет расширения, добавляем index.html
	if filepath.Ext(filePath) == "" {
		filePath = filepath.Join(filePath, "index.html")
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Ошибка создания директории %s: %v", dir, err)
	}

	return filePath
}

func Save(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return errutils.Wrap("storage.Save", err)
	}
	return nil
}
