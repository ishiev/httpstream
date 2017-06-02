package main

import (
	"io"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

// SaveStream сохраняет новый поток и возвращает его уникальный идентификатор
// При ошибке возвращается пустая строка
func SaveStream(stream io.Reader) (string, error) {
	uustr := uuid.NewV4().String()

	out, err := os.Create(GetStreamPath(uustr))
	if err != nil {
		return "", err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, stream); err != nil {
		return "", err
	}
	err = out.Sync()

	return uustr, err
}

// GetStream открывает и копирует сохраненный поток. Возвращает nil или ошибку
/*
func GetStream(out io.Writer, id string) error {
	in, err := os.Open(id)
	if err != nil {
		return err
	}
	defer func() {
		cerr := in.Close()
		if err == nil {
			err = cerr
		}
	}()
	_, err = io.Copy(out, in)
	return err
}
*/

// DeleteStream удаляет поток
func DeleteStream(id string) error {
	return os.Remove(GetStreamPath(id))
}

// GetStreamPath возвращает путь к потоку данных по его идентификатору
func GetStreamPath(id string) string {
	return filepath.Join(config.dataDir, id) + ".data"
}
