package v8repo

import (
	"errors"
	"fmt"
	"strings"
)

type ListTree struct {
	value    string
	elements []ListTree
	isValue  bool
}

func (l *ListTree) Load(data string) (*ListTree, error) {
	ReadListTreeData(data, l)
	return l, nil
}

func (l *ListTree) AppendValue(value string) {
	element := NewListTree()
	element.isValue = true
	element.value = value
	l.elements = append(l.elements, *element)
}

func (l *ListTree) AppendElement(element *ListTree) {
	l.elements = append(l.elements, *element)
}

func (l *ListTree) Length() int {
	return len(l.elements)
}

func (l *ListTree) ToString(params ...int) string {
	if l.isValue {
		return fmt.Sprintf("[%2d] %s\n", 0, l.value)
	}
	return elementToString(l, "", 0)
}

func (l *ListTree) Get(index ...int) (*ListTree, error) {
	if len(index) == 0 {
		return nil, errors.New("null index")
	}

	rl := l
	for _, i := range index {
		if len(rl.elements) < i {
			return nil, fmt.Errorf("inxdex %d is out of range %d", i, len(rl.elements))
		}
		rl = &rl.elements[i]
	}

	return rl, nil
}

func (l *ListTree) GetValue(index ...int) (string, error) {
	rl, err := l.Get(index...)
	if err != nil {
		return "", err
	}

	if !rl.isValue {
		return "", errors.New("element is not a value")
	}

	return rl.value, nil
}

func (l *ListTree) GetValues(index ...int) ([]string, error) {
	values := make([]string, 0)

	rl, err := l.Get(index...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < rl.Length(); i++ {
		fd, _ := rl.Get(i)
		if fd.isValue {
			values = append(values, fd.value)
		}
	}

	return values, nil
}

func NewListTree() *ListTree {
	list := new(ListTree)
	list.value = ""
	list.isValue = false
	list.elements = make([]ListTree, 0)

	return list
}

func elementToString(l *ListTree, prefix string, level int) string {
	var elementString string
	next := level + 1
	for i, item := range l.elements {
		levelPrefix := fmt.Sprintf("%s%2d", prefix, i)
		if item.isValue {
			// itemValue := fmt.Sprintf("%d", []byte(item.Value()))
			itemValue := item.value
			elementString += fmt.Sprintf("[%s] %s|%s%s\n", levelPrefix, strings.Repeat(".  ", 15-level), strings.Repeat(".  ", level), itemValue)
			continue
		}
		// fmt.Printf("[%s] %s%s\n", levelPrefix, strings.Repeat(". ", level), "+")
		nextPrefix := fmt.Sprintf("%s.", levelPrefix)
		elementString += elementToString(&item, nextPrefix, next)
	}
	return elementString
}

func ReadListTreeData(data string, list *ListTree) {
	data = strings.Trim(data, string(rune(65279)))
	data = strings.Trim(data, "{")
	data = strings.Trim(data, "}")

	ch := readSection(data, 0, list)
	<-ch
}

func readSection(data string, level int, list *ListTree) <-chan string {
	level++
	outCh := make(chan string)

	go func() {
		for i := 0; len(data) > 0; i++ {
			sym := data[0]

			if sym == ',' {
				data = data[1:]
				continue
			}

			if sym == '{' {
				data = data[1:]

				element := NewListTree()

				inCh := readSection(data, level, element)
				data = <-inCh

				list.AppendElement(element)

				continue
			}

			if sym == '}' {
				data = data[1:]
				outCh <- data
				break
			}

			if sym != '\n' && sym != '\r' {
				d := nextDelimeter(data)
				if d > 0 {
					value := data[0:d]
					data = data[d:]

					list.AppendValue(value)

					continue
				} else {
					value := data
					data = data[len(data):]

					list.AppendValue(value)

					continue
				}
			}

			data = data[1:]
		}
		close(outCh)
	}()

	return outCh
}

func nextDelimeter(data string) int {
	if data[0] == '"' {
		for i := 1; i < len(data); i++ {
			if data[i] == '"' {
				if len(data)-i > 4 && data[i+1] == '"' && data[i+2] == '"' && data[i+3] == '"' {
					i += 3
					continue
				}
				if len(data)-i > 2 && data[i+1] == '"' {
					i += 1
					continue
				}
				return i + 1
			}
		}
	}
	for i := 0; i < len(data); i++ {
		if data[i] == ',' || data[i] == '}' {
			return i
		}
	}

	return -1
}
