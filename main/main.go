package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zd4r/huffman"
)

// chooseAction - Выобр действия: архивация / разархивация
func chooseAction() int {
	var action int
	for action != 1 && action != 2 && action != 3 {
		fmt.Println("Выберите действие:\n1. Архивировать и разархивировать набор файлов\n2. Архивировать файл\n3. Разархивировать файл")
		_, err := fmt.Scanf("%d\n", &action)
		if err != nil {
			fmt.Printf("error choosing action: %v\n", err)
		}
	}
	return action
}

type NewEncodingDictStruct struct {
	NewEncodingDict map[string]byte `json:"new_encoding_dict"`
	FileExt         string          `json:"file_ext"`
}

// startArchiver - Запуск архиватора
func startArchiver() error {
	action := chooseAction()

	switch action {
	case 1:
		var path string
		fmt.Println("Введите путь до папки с файлами для архивации:")
		_, err := fmt.Scanf("%s\n", &path)
		if err != nil {
			return fmt.Errorf("error getting path: %v", err)
		}

		// Считываем файлы из папки с файлами для архивации
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error getting files from given path: %v", err)
		}

		// Итерируемся по файлам в указанной директории
		for _, f := range files {
			fileSize := f.Size()
			if !f.IsDir() && fileSize != 0 {
				// Читаем файл
				file, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
				if err != nil {
					return fmt.Errorf("error reading file: %v", err)
				}

				// Определяем расширение файла
				fileExt := filepath.Ext(f.Name())

				// Архивируем полученный файл
				packedFile, encodingDict, err := huffman.Pack(file)
				if err != nil {
					return fmt.Errorf("error packing file: %v", err)
				}
				fileNameWithoutExtension := f.Name()[:len(f.Name())-len(fileExt)]
				err = os.MkdirAll(filepath.Join(path, f.Name()+".Archive"), os.ModePerm)
				if err != nil {
					return fmt.Errorf("error creating archive folder: %v", err)
				}
				err = ioutil.WriteFile(filepath.Join(path, f.Name()+".Archive", fileNameWithoutExtension+"Packed.bin"), packedFile, 0644)
				if err != nil {
					return fmt.Errorf("error writing archive file: %v", err)
				}

				// Создаем новую мапу с кодами символов, где ключи и значения меняются местами для более эффективного поиска в дальнейшем
				newEncodingDict := make(map[string]byte)
				for k, v := range encodingDict {
					newEncodingDict[v] = k
				}

				// Инициализируем структуру хранящую мапу с кодами символов и расширение файла
				var NewEncodingDictStruct NewEncodingDictStruct
				NewEncodingDictStruct.NewEncodingDict = newEncodingDict
				NewEncodingDictStruct.FileExt = fileExt

				// Создаем Json файл с новой мапой кодов символов и расширением файла для дальнейшей разархивации файла
				newEncodingDictJson, err := json.Marshal(NewEncodingDictStruct)
				if err != nil {
					return fmt.Errorf("error creating json with encoding dict: %v", err)
				}
				err = ioutil.WriteFile(filepath.Join(path, f.Name()+".Archive", "encodingDict.json"), newEncodingDictJson, 0644)
				if err != nil {
					return fmt.Errorf("error writing encoding dict file: %v", err)
				}
				fmt.Printf("Архивация файла %v выполнена\n", f.Name())

				// Разархивируем файл
				unpackedFile, err := huffman.Unpack(packedFile, NewEncodingDictStruct.NewEncodingDict)
				if err != nil {
					return fmt.Errorf("error upacking file: %v", err)
				}
				err = ioutil.WriteFile(filepath.Join(path, f.Name()+".Archive", fileNameWithoutExtension+"Unpacked"+fileExt), unpackedFile, 0644)
				if err != nil {
					return fmt.Errorf("error writing unpacked file: %v", err)
				}
				fmt.Printf("Разархивация файла %v выполнена\n", fileNameWithoutExtension+"Unpacked"+fileExt)
			}
		}

	case 2:
		var filePath string
		fmt.Println("Введите путь до файла для архивации:")
		_, err := fmt.Scanf("%s\n", &filePath)
		if err != nil {
			return fmt.Errorf("error getting path: %v", err)
		}

		// Читаем файл
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		// Получение объекта FileInfo
		f, err := os.Lstat(filePath)
		if err != nil {
			return fmt.Errorf("error getting FileInfo of current file: %v", err)
		}

		// Проверка не является ли файл пустым
		if f.Size() == 0 {
			return errors.New("file is empty")
		}

		// Выделение пути до папки, где находится указанный файл
		path := filepath.Dir(filePath)

		// Определяем расширение файла
		fileExt := filepath.Ext(f.Name())

		// Архивируем полученный файл
		packedFile, encodingDict, err := huffman.Pack(file)
		if err != nil {
			return fmt.Errorf("error packing file: %v", err)
		}

		fileNameWithoutExtension := f.Name()[:len(f.Name())-len(fileExt)]
		err = os.MkdirAll(filepath.Join(path, f.Name()+".Archive"), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating archive folder: %v", err)
		}
		err = ioutil.WriteFile(filepath.Join(path, f.Name()+".Archive", fileNameWithoutExtension+"Packed.bin"), packedFile, 0644)
		if err != nil {
			return fmt.Errorf("error writing archive file: %v", err)
		}

		// Создаем новую мапу с кодами символов, где ключи и значения меняются местами для более эффективного поиска в дальнейшем
		newEncodingDict := make(map[string]byte)
		for k, v := range encodingDict {
			newEncodingDict[v] = k
		}

		// Инициализируем структуру хранящую мапу с кодами символов и расширение файла
		var NewEncodingDictStruct NewEncodingDictStruct
		NewEncodingDictStruct.NewEncodingDict = newEncodingDict
		NewEncodingDictStruct.FileExt = fileExt

		// Создаем Json файл с новой мапой кодов символов для дальнейшей разархивации файла
		newEncodingDictJson, err := json.Marshal(NewEncodingDictStruct)
		if err != nil {
			return fmt.Errorf("error creating json with encoding dict: %v", err)
		}
		err = ioutil.WriteFile(filepath.Join(path, f.Name()+".Archive", "encodingDict.json"), newEncodingDictJson, 0644)
		if err != nil {
			return fmt.Errorf("error writing encoding dict file: %v", err)
		}
		fmt.Printf("Архивация файла %v выполнена\n", f.Name())

	case 3:
		var filePath string
		fmt.Println("Введите путь до файла для разархивации (json файл для раскодировки должен находиться в той же директории):")
		_, err := fmt.Scanf("%s\n", &filePath)
		if err != nil {
			return fmt.Errorf("error getting path: %v", err)
		}

		// Читаем файл
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		// Получение объекта FileInfo
		f, err := os.Lstat(filePath)
		if err != nil {
			return fmt.Errorf("error getting FileInfo of current file: %v", err)
		}

		// Проверка не является ли файл пустым
		if f.Size() == 0 {
			return errors.New("file is empty")
		}

		// Выделение пути до папки, где находится указанный файл
		path := filepath.Dir(filePath)

		// Читаем файл с кодами символов для раскодировки
		newEncodingDictFile, err := ioutil.ReadFile(filepath.Join(path, "encodingDict.json"))
		if err != nil {
			return fmt.Errorf("error reading EncodingDict file: %v", err)
		}
		var NewEncodingDictStruct NewEncodingDictStruct
		err = json.Unmarshal(newEncodingDictFile, &NewEncodingDictStruct)
		if err != nil {
			return fmt.Errorf("error parsing json EncodingDict file: %v", err)
		}

		// Получаем расширение файла
		fileExt := NewEncodingDictStruct.FileExt

		// Название файла без расщирения
		fileNameWithoutExtension := f.Name()[:len(f.Name())-len(fileExt)]

		// Разархивирование файла
		unpackedFile, err := huffman.Unpack(file, NewEncodingDictStruct.NewEncodingDict)
		if err != nil {
			return fmt.Errorf("error upacking file: %v", err)
		}
		err = ioutil.WriteFile(filepath.Join(path, fileNameWithoutExtension+"Unpacked"+fileExt), unpackedFile, 0644)
		if err != nil {
			return fmt.Errorf("error writing unpacked file: %v", err)
		}
		fmt.Printf("Разрхивация файла %v выполнена\n", f.Name())

	}
	return nil
}

func main() {
	//defer time.Sleep(1 * time.Minute)
	defer fmt.Println("Application will be closed in 1 minute")
	err := startArchiver()
	if err != nil {
		fmt.Printf("archiver work error: %v\n", err)
	}
}
