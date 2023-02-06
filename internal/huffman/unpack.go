package huffman

import (
	"fmt"
	"strings"
)

// BytesToBits - перевод слайса байт в последовательность битов
func BytesToBits(data []byte) (string, error) {
	var result strings.Builder
	for _, b := range data {
		result.WriteString(fmt.Sprintf("%08b", b))
	}
	return result.String(), nil
}

// BitsToBytesUnpack - преобразовение последовательности бит в байты при распаковке файла
func BitsToBytesUnpack(data string, encodingDict map[string]byte) ([]byte, error) {
	// Определеям кол-во незначащих битов в конце последовательности
	emptyBitsNum, err := BitsToByte(data[:8])
	if err != nil {
		return []byte{}, fmt.Errorf("error getting emptyBitsNum: %v", err)
	}

	// Преобразуем исходную последовательность бит
	data = data[8:(len(data) - int(emptyBitsNum))]

	// Переводим преобразованную последовательность бит в слайс байт
	result, err := DecodeBits(data, encodingDict)
	if err != nil {
		return []byte{}, fmt.Errorf("error decoding bits while unpack: %v", err)
	}
	return result, err
}

// DecodeBits - декодирование битов с помощтю encodingDict (мапы кодов символов)
func DecodeBits(data string, encodingDict map[string]byte) ([]byte, error) {
	var currentCode strings.Builder
	var result []byte

	for _, b := range data {
		currentCode.WriteString(string(b))
		if v, ok := encodingDict[currentCode.String()]; ok {
			result = append(result, v)
			currentCode.Reset()
		}
	}
	return result, nil
}

// Unpack - Разархивирование файла
func Unpack(data []byte, encodingDict map[string]byte) ([]byte, error) {
	// Получение последовательности бит из слайса байт
	dataBits, err := BytesToBits(data)
	if err != nil {
		return []byte{}, fmt.Errorf("error getting bits from bytes: %v", err)
	}

	// Преобразуем последовательность бит в слайс байт с помощью кодов символов
	result, err := BitsToBytesUnpack(dataBits, encodingDict)
	if err != nil {
		return []byte{}, fmt.Errorf("error converting bit to slice byte: %v", err)
	}
	return result, nil
}
