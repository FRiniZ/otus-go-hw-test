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
	head *ListItem
	tail *ListItem
	size int
}

func (l *list) Len() int {
	return (l.size)
}

func (l *list) Front() *ListItem {
	return (l.head)
}

func (l *list) Back() *ListItem {
	return (l.tail)
}

func (l *list) PushFront(v interface{}) *ListItem {
	elm := &ListItem{v, nil, nil}

	if l.size == 0 {
		l.tail = elm
	} else {
		l.head.Prev = elm
		elm.Next = l.head
	}

	l.head = elm
	l.size++

	return (l.head)
}

func (l *list) PushBack(v interface{}) *ListItem {
	elm := &ListItem{v, nil, nil}

	if l.size == 0 {
		l.head = elm
	} else {
		l.tail.Next = elm
		elm.Prev = l.tail
	}

	l.tail = elm
	l.size++

	return (l.tail)
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	elmPrev := i.Prev
	elmNext := i.Next

	if elmPrev != nil {
		elmPrev.Next = elmNext
		if elmPrev.Next == nil {
			l.tail = elmPrev
		}
	}

	if elmNext != nil {
		elmNext.Prev = elmPrev
		if elmNext.Prev == nil {
			l.head = elmNext
		}
	}

	l.size--
	if l.size == 0 {
		l.head = nil
		l.tail = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}

	elmPrev := i.Prev
	elmNext := i.Next

	if elmPrev != nil {
		elmPrev.Next = elmNext
		if elmPrev.Next == nil {
			l.tail = elmPrev
		}
	}

	if elmNext != nil {
		elmNext.Prev = elmPrev
	}

	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i
}

func NewList() List {
	return new(list)
}
