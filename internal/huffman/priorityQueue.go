package huffman

// PriorityQueue - Очередь с приоритетом
type PriorityQueue []*Tree

// Len - Длина очереди
func (pq PriorityQueue) Len() int { return len(pq) }

// Less - Сравнение приоритетов
func (pq PriorityQueue) Less(i, j int) bool {
	// чтобы Pop давал самый низкий приоритет, используем оператор меньше
	return pq[i].Priority < pq[j].Priority
}

// Swap - Поменять элементы местами
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push - Добавить элемент в очередь
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Tree)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop - Достать элемент из очереди
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // избежать утечки памяти
	item.Index = -1 // для безопасности
	*pq = old[0 : n-1]
	return item
}
