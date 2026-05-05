package main

// -l = disable inlining
// -m = print optimization decisions
// go build -gcflags '-l -m'

func getResultInStack() int {
	result := 200
	return result
}

func getResultInHeap() *int {
	result := 200
	return &result
}

func getResult(number *int) int {
	result := *number * 2
	return result
}

func printValue(v interface{}) {
	println(v)
}

func createPointer() *int {
	value2 := new(int)
	return value2
}

func main() {
	_ = getResultInStack()
	_ = getResultInHeap()

	number := 100
	_ = getResult(&number)

	var num1 int = 10
	var str1 string = "hello"

	printValue(num1)
	printValue(str1)

	var num2 int = 20
	var str2 string = "world"

	var i interface{}
	i = num2
	i = str2
	_ = i

	value1 := new(int) // stack
	_ = value1

	value2 := createPointer() // heap
	_ = value2
}
