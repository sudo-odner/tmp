package main

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

var (
	ErrOutOfMemory    = errors.New("out of memory")
	ErrInvalidPointer = errors.New("invalid pointer")
	ErrDoubleFree     = errors.New("double free")
)

// BlockAllocator представляет собой потокобезопасный блочный аллокатор памяти,
// который выделяет память блоками фиксированного размера из заранее созданного пула.
type BlockAllocator struct {
	mu        sync.Mutex
	memory    []byte
	blockSize int
	numBlocks int
	freeList  []int  // Список индексов свободных блоков
	isFree    []bool // Флаги для предотвращения двойного освобождения
}

// NewBlockAllocator создает новый аллокатор с заданным размером блока и их количеством.
func NewBlockAllocator(blockSize, numBlocks int) *BlockAllocator {
	alloc := &BlockAllocator{
		memory:    make([]byte, blockSize*numBlocks),
		blockSize: blockSize,
		numBlocks: numBlocks,
		freeList:  make([]int, numBlocks),
		isFree:    make([]bool, numBlocks),
	}

	// Изначально все блоки свободны
	for i := 0; i < numBlocks; i++ {
		alloc.freeList[i] = i
		alloc.isFree[i] = true
	}
	return alloc
}

// Alloc возвращает срез, представляющий свободный блок памяти.
func (a *BlockAllocator) Alloc() ([]byte, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.freeList) == 0 {
		return nil, ErrOutOfMemory
	}

	// Берем последний свободный блок из списка (pop)
	idx := a.freeList[len(a.freeList)-1]
	a.freeList = a.freeList[:len(a.freeList)-1]
	a.isFree[idx] = false

	offset := idx * a.blockSize
	return a.memory[offset : offset+a.blockSize], nil
}

// Free возвращает блок памяти обратно в аллокатор.
func (a *BlockAllocator) Free(block []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Проверяем, что размер блока соответствует ожиданиям
	if len(block) != a.blockSize || cap(block) == 0 {
		return ErrInvalidPointer
	}

	startAddr := uintptr(unsafe.Pointer(&a.memory[0]))
	blockAddr := uintptr(unsafe.Pointer(&block[0]))

	// Проверяем, что адрес блока находится в пределах нашего пула памяти
	if blockAddr < startAddr || blockAddr >= startAddr+uintptr(len(a.memory)) {
		return ErrInvalidPointer
	}

	offset := blockAddr - startAddr

	// Проверяем, что блок выровнен по размеру блока
	if offset%uintptr(a.blockSize) != 0 {
		return ErrInvalidPointer
	}

	idx := int(offset) / a.blockSize

	// Защита от двойного освобождения
	if a.isFree[idx] {
		return ErrDoubleFree
	}

	a.isFree[idx] = true
	a.freeList = append(a.freeList, idx)

	return nil
}

func main() {
	fmt.Println("Инициализация блочного аллокатора...")
	alloc := NewBlockAllocator(64, 3) // 3 блока по 64 байта

	b1, err := alloc.Alloc()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Выделен блок 1: %p\n", &b1[0])

	b2, err := alloc.Alloc()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Выделен блок 2: %p\n", &b2[0])

	b3, err := alloc.Alloc()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Выделен блок 3: %p\n", &b3[0])

	_, err = alloc.Alloc()
	if err != nil {
		fmt.Printf("Ожидаемая ошибка при выделении 4-го блока: %v\n", err)
	}

	fmt.Println("Освобождение блока 2...")
	err = alloc.Free(b2)
	if err != nil {
		panic(err)
	}

	b4, err := alloc.Alloc()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Выделен блок 4 (должен быть на месте блока 2): %p\n", &b4[0])
}
