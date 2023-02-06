package huffman

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// BitsToBytes - перевод двоичного числа в слайс байт
func BitsToBytes(val string) ([]byte, error) {
	// Заполняем 0-ями недостающие биты
	emptyBitsNum := 8 - len(val)%8
	if emptyBitsNum != 0 {
		for i := 0; i < emptyBitsNum; i++ {
			val += "0"
		}
	}

	// Определяем кол-во байт в заданном значении
	var bytesNum int
	bytesNum = len(val) / 8

	// Перевод двоичного значения в слайс байт
	var result []byte
	for i := 0; i < bytesNum; i++ {
		b, err := BitsToByte(val[8*i : 8*(i+1)])
		if err != nil {
			return result, fmt.Errorf("error converting provided value to bytes: %v", err)
		}
		result = append(result, b)
	}

	// В первом байте последовательности храним кол-во бит в его конце, не содержащих информацию из файла
	result = append([]byte{byte(emptyBitsNum)}, result...)
	return result, nil
}

// BitsToByte - перевод двоичного числа в байт
func BitsToByte(val string) (byte, error) {
	var intResult int8
	for num := range val {
		digit, err := strconv.Atoi(string(val[num]))
		if err != nil {
			return byte(intResult), fmt.Errorf("error converting provided value to byte: %v", err)
		}
		intResult += int8(math.Pow(2, float64(len(val)-1-num))) * int8(digit)
	}
	return byte(intResult), nil
}

// CompressData - Сжимает данные при помощи полученных кодов символов, хранящихся в encodingDict
func CompressData(data []byte, encodingDict map[byte]string) string {
	var result strings.Builder
	for _, b := range data {
		result.WriteString(encodingDict[b])
	}
	return result.String()
}

// Pack - архивирование файла
func Pack(data []byte) ([]byte, map[byte]string, error) {
	// Создаем мапу частот байтов
	frequencyMap := make(map[byte]int)
	for j := 0; j < len(data); j++ {
		frequencyMap[data[j]] += 1
	}

	// Создаем очередь с приоритетом из деревьев (чем меньше частота, тем больше приоритет)
	pq := make(PriorityQueue, 0)
	for k, v := range frequencyMap {
		tempNode := Node{
			Val:       k,
			Frequency: v,
		}
		tempTree := Tree{
			Root:     tempNode,
			Priority: tempNode.Frequency,
		}
		heap.Push(&pq, &tempTree) // Сложность O(log n)
	}

	// Собираем все деревья в одно большое дерево
	for len(pq) > 1 {
		tree1 := heap.Pop(&pq).(*Tree) // Сложность O(log n)
		tree2 := heap.Pop(&pq).(*Tree) // Сложность O(log n)

		tempNode := Node{
			Frequency: tree1.Root.Frequency + tree2.Root.Frequency,
			Left:      &tree1.Root,
			Right:     &tree2.Root,
		}
		tempTree := Tree{
			Root:     tempNode,
			Priority: tempNode.Frequency,
		}

		heap.Push(&pq, &tempTree) // Сложность O(log n)
	}

	// Создаем мапу сжатия (определяем коды символов)
	encodingDict := make(map[byte]string)
	finalTree := heap.Pop(&pq).(*Tree) // Сложность O(log n)
	finalTree.Root.CreateCompressionMap("", "", encodingDict)

	// Сжимаем данные при помощи полученных кодов символов
	compressedData := CompressData(data, encodingDict)

	// Преобразуем последовательнсоть бит в салйс байт
	result, err := BitsToBytes(compressedData)
	if err != nil {
		return []byte{}, nil, fmt.Errorf("error compressing file: %v", err)
	}
	return result, encodingDict, nil
}
