package main

import (
	"errors"
	"fmt"
	"unsafe"
)

type LinearAllocator struct {
	data []byte
}

func NewLinearAllocator(capacity int) (LinearAllocator, error) {
	if capacity < 0 {
		return LinearAllocator{}, errors.New("capacity can't be less zero")
	}
	return LinearAllocator{
		data: make([]byte, 0, capacity),
	}, nil
}

func (a *LinearAllocator) Allocate(size int) (unsafe.Pointer, error) {
	oldLen := len(a.data)
	newLen := oldLen + size
	if newLen > cap(a.data) {
		return nil, errors.New("Size exceeds the permissible limit")
	}
	a.data = a.data[:newLen]
	return unsafe.Pointer(&a.data[oldLen]), nil
}

func (a *LinearAllocator) Free() {
	a.data = a.data[:0]
}

func store[T any](pointer unsafe.Pointer, value T) {
	*(*T)(pointer) = value
}

func load[T any](pointer unsafe.Pointer) T {
	return *(*T)(pointer)
}

func main() {
	const MB = 1 << 20
	// Init linear allocator
	allocator, err := NewLinearAllocator(MB)
	if err != nil {
		panic(err)
	}

	defer allocator.Free()

	// allocate pointer on value
	pointer1, _ := allocator.Allocate(8)
	pointer2, _ := allocator.Allocate(8)

	// add value in allocate pointer
	store[int](pointer1, 234)
	store[int](pointer2, 200)

	// get value form allocate
	value1 := load[int](pointer1)
	value2 := load[int](pointer2)

	// print value and addres
	fmt.Println(pointer1, value1)
	fmt.Println(pointer2, value2)
}
