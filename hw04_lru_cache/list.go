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
	Prev  *ListItem
	Next  *ListItem
}

type list struct {
	count     int
	firstElem *ListItem
	lastElem  *ListItem
}

func (l *list) init(item *ListItem) {
	l.count = 1
	l.firstElem = item
	l.lastElem = item
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.firstElem
}

func (l *list) Back() *ListItem {
	return l.lastElem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.count == 0 {
		l.init(item)
		return item
	}
	first := l.firstElem
	first.Prev = item
	item.Next = first
	l.firstElem = item
	l.count++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.count == 0 {
		l.init(item)
		return item
	}
	item.Prev = l.lastElem
	l.lastElem.Next = item
	l.lastElem = item
	l.count++
	return item
}

func (l *list) Remove(i *ListItem) {
	if l.count == 0 {
		return
	}
	if i == nil {
		return
	}
	if i.Prev == nil && i.Next == nil {
		return
	}
	prev := i.Prev
	next := i.Next
	i.Next, i.Prev = nil, nil

	if prev != nil {
		prev.Next = next
	} else {
		l.firstElem = next
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.lastElem = prev
	}

	l.count--
}

func (l *list) MoveToFront(item *ListItem) {
	if item == nil {
		return
	}
	if l.count == 0 {
		l.init(item)
		return
	}
	if item == l.firstElem {
		return
	}

	first := l.firstElem
	first.Prev = item

	prev := item.Prev
	next := item.Next
	prev.Next = next

	if next != nil {
		next.Prev = prev
	} else {
		l.lastElem = prev
	}

	item.Prev = nil
	item.Next = first

	l.firstElem = item
}

func NewList() List {
	return new(list)
}
