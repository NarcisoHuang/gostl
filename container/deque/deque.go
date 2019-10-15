package deque

import (
	"errors"
	"fmt"
)

var ErrOutOffRange = errors.New("out off range")

type Deque struct {
	data  []interface{}
	begin int
	end   int
	size  int
}

func New(capacity int) *Deque {
	if capacity == 0 {
		capacity = 1
	}
	return &Deque{
		data: make([]interface{}, capacity, capacity),
	}
}

//Size return the size of deque
func (this *Deque) Size() int {
	return this.size
}

//Capacity return the capacity of deque
func (this *Deque) Capacity() int {
	return len(this.data)
}

func (this *Deque) Empty() bool {
	if this.Size() == 0 {
		return true
	}
	return false
}

func (this *Deque) expandIfNeeded() {
	if this.size == this.Capacity() {
		newCapacity := this.size * 2
		if newCapacity == 0 {
			newCapacity = 1
		}
		data := make([]interface{}, newCapacity, newCapacity)
		for i := 0; i < this.size; i++ {
			data[i] = this.data[(this.begin+i)%this.Capacity()]
		}
		this.data = data
		this.begin = 0
		this.end = this.size
	}
}

func (this *Deque) shrinkIfNeeded() {
	if int(float64(this.size*2)*1.2) < this.Capacity() {
		newCapacity := this.Capacity() / 2
		data := make([]interface{}, newCapacity, newCapacity)
		for i := 0; i < this.size; i++ {
			data[i] = this.data[(this.begin+i)%this.Capacity()]
		}
		this.data = data
		this.begin = 0
		this.end = this.size
	}
}

func (this *Deque) PushBack(value interface{}) {
	this.expandIfNeeded()
	this.data[this.end] = value
	this.end = this.nextIndex(this.end)
	this.size++
}

func (this *Deque) PushFront(value interface{}) {
	this.expandIfNeeded()
	this.begin = this.preIndex(this.begin)
	this.data[this.begin] = value
	this.size++
}

func (this *Deque) Insert(pos int, value interface{}) error {
	if pos < 0 || pos > this.size {
		return ErrOutOffRange
	}
	if pos == 0 {
		this.PushFront(value)
		return nil
	}
	if pos == this.size {
		this.PushBack(value)
		return nil
	}

	this.expandIfNeeded()
	if pos < this.size-pos {
		//move the front pos items
		idx := this.preIndex(this.begin)
		for i := 0; i < pos; i++ {
			this.data[idx] = this.data[this.nextIndex(idx)]
			idx = this.nextIndex(idx)
		}
		this.data[idx] = value
		this.begin = this.preIndex(this.begin)
	} else {
		//move the back pos items
		idx := this.end
		for i := 0; i < this.size-pos; i++ {
			this.data[idx] = this.data[this.preIndex(idx)]
			idx = this.preIndex(idx)
		}
		this.data[idx] = value
		this.end = this.nextIndex(this.end)
	}
	this.size++
	return nil
}

func (this *Deque) PopBack() interface{} {
	if this.Empty() {
		return nil
	}
	index := this.preIndex(this.end)
	val := this.data[index]
	this.data[index] = nil
	this.size--
	this.shrinkIfNeeded()
	return val
}

func (this *Deque) PopFront() interface{} {
	if this.Empty() {
		return nil
	}
	val := this.data[this.begin ]
	this.data[this.begin] = nil
	this.begin = this.nextIndex(this.begin)
	this.size--
	this.shrinkIfNeeded()
	return val
}

func (this *Deque) At(pos int) interface{} {
	if pos < 0 || pos >= this.size {
		return nil
	}
	return this.data[(pos+this.begin)%this.Capacity()]
}

func (this *Deque) Back() interface{} {
	return this.At(this.size - 1)
}

func (this *Deque) Front() interface{} {
	return this.At(0)
}

func (this *Deque) nextIndex(index int) int {
	return (index + 1) % this.Capacity()
}

func (this *Deque) preIndex(index int) int {
	return (index - 1 + this.Capacity()) % this.Capacity()
}

func (this *Deque) Erase(pos int) error {
	return this.EraseRange(pos, pos+1)
}

//EraseRange erase the data in the range [firstPos, lastPos), not include lastPos.
func (this *Deque) EraseRange(firstPos, lastPos int) error {
	if firstPos < 0 || lastPos > this.size {
		return ErrOutOffRange
	}
	if firstPos >= lastPos {
		return nil
	}
	eraseNum := lastPos - firstPos
	leftNum := firstPos
	rightNum := this.size - lastPos

	if leftNum <= rightNum {
		//move left data
		idx := (this.begin + this.preIndex(lastPos)) % this.Capacity()
		for i := 0; i < leftNum; i++ {
			tempIndex := (idx - eraseNum + this.Capacity()) % this.Capacity()
			this.data[idx] = this.data[tempIndex]
			idx = this.preIndex(idx)
		}
		this.begin = this.nextIndex(idx)
		for i := 0; i < eraseNum; i++ {
			this.data[idx] = nil
			idx = this.preIndex(idx)
		}

	} else {
		idx := (this.begin + firstPos) % this.Capacity()
		for i := 0; i < rightNum; i++ {
			tempIndex := (idx + eraseNum + this.Capacity()) % this.Capacity()
			this.data[idx] = this.data[tempIndex]
			idx = this.nextIndex(idx)
		}
		this.end = idx
		for i := 0; i < eraseNum; i++ {
			this.data[idx] = nil
			idx = this.nextIndex(idx)
		}
	}
	this.size -= eraseNum
	this.shrinkIfNeeded()
	return nil
}

func (this *Deque) String() string {
	str := "["
	for i := 0; i < this.size; i++ {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprintf("%v", this.data[(this.begin+i)%this.Capacity()])
	}
	str += "]"
	return str
}