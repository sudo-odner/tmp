package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 1. Сделать структуры Base и Child.
//  2. Структура Base должна содержать строковое поле name.
//  3. Структура Child должна содержать строковое поле lastName.
//  4. Сделать функцию Say у структуры Base, которая распечатывает
//     на экране: Hello, %name!
//  5. Пронаследовать Child от Base.
//  6. Инициализировать экземпляр b1 Base.
//     присвоить name значение Parent
//  7. Инициализировать экземпляр c1 Сhild.
//     присвоить name значение Child
//     присвоить lastName значение Inherited
//  8. Вызвать у обоих экземпляров метод Say.
//  9. Переопределить метод Say для структуры Child, чтобы он выводил
//     на экран: Hello, %lastName %name!
//  10. Сделать массив, содержащий b1 и c1.
//  11. Вызвать Say у всех элементов массива из шага 10.
//  12. Сделать метод NewObject для создания экземпляров Base и Child
//     в зависимости от входного параметра.
//  13. Написать юнит тесты для метода NewObject.
//  14. Сделать генератор объектов Base и Child такой, чтобы:
//     объекты Base создавались в фоновом потоке с задержкой 1 секунда;
//     объекты Child создавались в фоновом потоке с задержкой 2 секунды;
//     общее время генерации объектов не превышало 11 секунд;
//  15. Сделать асинхронный обработчик сгенерированных объектов такой, чтобы:
//     метод Say вызывался в порядке генерации объектов;
//     не приводил к утечкам памяти;

type Base struct {
	name string
}

func (b Base) Say() {
	fmt.Println("Hellow,", b.name)
}

type Child struct {
	Base
	lastName string
}

func (c Child) Say() {
	fmt.Println("Hellow,", c.lastName, c.name)
}

type Sayer interface {
	Say()
}

func NewObject(nameObject string) Sayer {
	switch nameObject {
	case "Base":
		return Base{
			name: "NewObjectBase",
		}
	case "Child":
		return Child{
			Base: Base{
				name: "NewObjectChild",
			},
			lastName: "NewObjectLastNameChild",
		}
	default:
		return nil
	}
}

func Workers(ch chan Sayer, f func()) {

}
func Generator() <-chan Sayer {
	ch := make(chan Sayer)
	ctx, chanel := context.WithTimeout(context.Background(), 11*time.Second)

	createObj := func(typeObject string) {
		select {
		case <-ctx.Done():
		case ch <- NewObject(typeObject):
		}
	}
	worker := func(ticker *time.Ticker, typeObject string) {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				createObj(typeObject)
			}
		}
	}
	go func() {
		defer func() {
			close(ch)
			chanel()
		}()
		var wg sync.WaitGroup

		wg.Go(func() {
			worker(time.NewTicker(1*time.Second), "Base")
		})
		wg.Go(func() {
			worker(time.NewTicker(2*time.Second), "Child")
		})
		wg.Wait()
	}()

	return ch
}

func main() {
	b1 := Base{
		name: "Parent",
	}

	c1 := Child{
		Base: Base{
			name: "Child",
		},
		lastName: "Inherited",
	}

	b1.Say()
	c1.Say()

	var arr = [2]Sayer{b1, c1}
	for _, el := range arr {
		el.Say()
	}
	c2 := NewObject("Child")

	c2.Say()

	for newObj := range Generator() {
		newObj.Say()
	}
}
