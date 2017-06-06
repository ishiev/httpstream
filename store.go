package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"time"

	"log"

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
		if err != nil {
			// Удаляем поток в случае наличия любых ошибок
			_ = DeleteStream(uustr)
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

// Clean очистка хранилища от устаревших данных
func Clean() error {
	files, err := ioutil.ReadDir(config.dataDir)
	if err != nil {
		log.Printf("TTL: Ошибка получения списка потоков хранилища: %s, очистка не выполнена\n", err.Error())
		return err
	}

	// Прохождение по списку файлов и удаление устаревших
	for _, file := range files {
		timeToDie := file.ModTime().Add(config.ttl)
		if timeToDie.Before(time.Now()) {
			// Удаление устаревшего потока
			path := filepath.Join(config.dataDir, file.Name())
			err = os.Remove(path)
			if err != nil {
				// TODO: инцидент безопасности, пока игнорируем, но пишем в лог
				log.Printf("TTL: Не удалось удалить поток: %s, время создания: %s, ошибка %s\n", path, file.ModTime().String(), err.Error())
			} else {
				log.Printf("TTL: Поток %s удален\n", path)
			}
		}
	}
	return nil
}
