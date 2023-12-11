package list

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrMismatchType = errors.New("mismatched type: the type of the provided value does not match the type of items already in the storage")

type List struct {
	len int64
	fn  *node
}

func Newlist() *List {
	return &List{len: 0, fn: nil}
}

func (l *List) Clear() {
	l.fn = nil
	l.len = 0
}

func (l *List) Len() int64 {
	return l.len
}

func (l *List) GetAll() (values []interface{}, ok bool) {
	curr := l.fn

	if curr == nil {
		return
	}

	for curr != nil {
		values = append(values, curr.value)
		curr = curr.next
	}
	return values, true
}

func (l *List) GetAllByValue(value interface{}) (ids []int64, ok bool) {
	valuesArr, _ := l.GetAll()
	if valuesArr != nil {
		for key, val := range valuesArr {
			if val == value {
				ids = append(ids, int64(key))
				ok = true
			}
		}
	}
	return
}

func CreateListFromSlice(inputSlice []int64) *List {
	var myList *List = Newlist()
	for _, val := range inputSlice {
		myList.Add(val)
	}
	return myList
}

func (l *List) Print() {
	var curr = l.fn
	for curr != nil {
		fmt.Print(curr.value, "[", curr.index, "]", "->")
		curr = curr.next
	}
	fmt.Println("nil")
}

func (l *List) getNodeByIndex(id int64) *node {
	if id < 0 || id > l.len-1 {
		return nil
	}
	var curr = l.fn
	for curr != nil {
		if curr.index == id {
			return curr
		}
		curr = curr.next
	}
	return nil
}

func (l *List) GetByIndex(id int64) (value interface{}, ok bool) {

	curr := l.getNodeByIndex(id)
	if curr == nil {
		return 0, false
	}
	return curr.value, true
}

func (l *List) GetByValue(value interface{}) (id int64, ok bool) {
	curr := l.fn
	for curr != nil {
		if curr.value == value {
			return curr.index, true
		}
		curr = curr.next
	}
	return 0, false
}

func (l *List) updateIndexes() {
	curr := l.fn
	var id int64 = 0
	for curr != nil {
		curr.index = id
		id++
		curr = curr.next
	}
}

func (l *List) Add(value interface{}) (int64, error) {
	newNode := &node{value: value, next: nil, index: l.len}

	if l.fn == nil {
		l.len++
		l.fn = newNode
		return 0, nil
	}

	if reflect.TypeOf(value) != reflect.TypeOf(l.fn.value) {
		return 0, ErrMismatchType
	}

	curr := l.fn
	for curr.next != nil {
		curr = curr.next
	}

	curr.next = newNode
	l.len++
	return l.len - 1, nil
}

func (l *List) RemoveByIndex(id int64) {
	switch {
	case id < 0 || id > l.len-1:
		return
	case id == l.len-1:
		if id == 0 {
			l.fn = nil
			l.updateIndexes()
			return
		} else {
			curr := l.getNodeByIndex(id - 1)
			curr.next = nil
			l.updateIndexes()
			return
		}
	case id == 0:
		l.fn = l.getNodeByIndex(1)
		l.updateIndexes()
		return

	default:
		l.getNodeByIndex(id - 1).next = l.getNodeByIndex(id + 1)
		l.updateIndexes()
		return
	}
}

func (l *List) Count(value interface{}) int {
	curr := l.fn
	var counter int = 0
	for curr != nil {
		if curr.value == value {
			counter++
		}
		curr = curr.next
	}
	return counter
}

func (l *List) RemoveAllByValue(value interface{}) {
	amount := l.Count(value)
	for amount != 0 {
		l.RemoveByValue(value)
		amount--
	}
}

func (l *List) RemoveByValue(value interface{}) {
	curr := l.fn
	if curr == nil {
		return
	}

	if curr.value == value {
		if curr.next == nil {
			l.fn = nil
			fmt.Println("List is empty")
			l.len--
			l.updateIndexes()
			return
		}
		l.fn = curr.next
		l.len--
		l.updateIndexes()
		return
	}

	for curr.next != nil {
		if curr.next.value == value {
			if curr.next.next != nil {
				curr.next = curr.next.next
				l.len--
				l.updateIndexes()
				return
			} else {
				curr.next = nil
				l.len--
				l.updateIndexes()
				return
			}

		} else {
			curr = curr.next
		}
	}
}

func (l *List) GetAllByValueSelectedFields(value interface{}) (ids []int64) {
	currNode := l.fn
	for currNode != nil {
		if IsEqualSelectedFields(currNode.value, value) {
			ids = append(ids, currNode.index)
		}
		currNode = currNode.next
	}
	return
}

func IsEqualSelectedFields(freeze interface{}, check interface{}) bool {
	freezeValue := reflect.ValueOf(freeze)
	checkValue := reflect.ValueOf(check)

	if freezeValue.Type() == checkValue.Type() {
		for i := 0; i < freezeValue.NumField(); i++ {

			fieldFreezeValue := freezeValue.Field(i)
			fieldCheckValue := checkValue.Field(i)

			if fieldFreezeValue.String() == "" {
				continue
			} else {
				if !reflect.DeepEqual(fieldFreezeValue.Interface(), fieldCheckValue.Interface()) {
					return false
				}
			}
		}
	} else {
		return false
	}
	return true
}
