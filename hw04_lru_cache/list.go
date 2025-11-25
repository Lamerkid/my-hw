package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstNode *ListItem
	lastNode  *ListItem
	length    int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstNode
}

func (l *list) Back() *ListItem {
	return l.lastNode
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v}
	// Инициализировать ноду, если список пустой
	if l.firstNode == nil {
		l.firstNode = newListItem
		l.lastNode = newListItem
	} else {
		// nil(ZeroValue) <- новый элемент <-> начало списка
		newListItem.Next = l.firstNode
		l.firstNode.Prev = newListItem
		// firstNode = новый элемент
		l.firstNode = newListItem
	}
	l.length++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v}
	// Инициализировать ноду, если список пустой
	if l.lastNode == nil {
		l.lastNode = newListItem
		l.firstNode = newListItem
	} else {
		// Конец списка <-> новый элемент -> nil(ZeroValue)
		newListItem.Prev = l.lastNode
		l.lastNode.Next = newListItem
		// lastNode = новый элемент
		l.lastNode = newListItem
	}
	l.length++
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	// При удалении последней ноды списка значения firstNode и lastNode становятся nil

	// Если i - начало списка, то firstNode = следующий элемент
	// Иначе i.prev -> i.next
	if i.Prev == nil {
		l.firstNode = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	// Если i - конец списка, то lastNode = предыдущий элемент
	// Иначе i.prev <- i.next
	if i.Next == nil {
		l.lastNode = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	// Если элемент в начале, то ничего не делаем
	if i.Prev == nil {
		return
	}
	// i.prev -> i.next (i.next может быть nil)
	i.Prev.Next = i.Next
	// Если i - конец списка, то lastNode = предыдущий элемент
	// Иначе i.prev <- i.next
	if i.Next == nil {
		l.lastNode = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	// nil <- i(firstNode) <-> бывшая firstNode
	i.Next = l.firstNode
	l.firstNode.Prev = i
	i.Prev = nil
	l.firstNode = i
}

func NewList() List {
	return new(list)
}
