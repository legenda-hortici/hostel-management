package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// Функция для сохранения аватара
func SaveAvatar(file multipart.File) (string, error) {
	// Генерируем имя для файла (например, с помощью uuid или просто текущего времени)
	fileName := fmt.Sprintf("%d.png", time.Now().Unix()) // или используйте другой метод генерации имени

	// Указываем путь для сохранения
	savePath := filepath.Join("web", "static", "img", "avatars", fileName)

	// Создаем директорию, если она не существует
	dir := filepath.Dir(savePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directories: %v", err)
	}

	// Создаем файл
	outFile, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Копируем данные из загруженного файла в новый файл
	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", err
	}
	savePath = "/" + savePath[4:]
	// Возвращаем путь к файлу
	return savePath, nil
}
