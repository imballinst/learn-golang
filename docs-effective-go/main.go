package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	fmt.Println("--- Allocation Test ---")
	runAllocationTest()

	fmt.Println("\n--- Loops Test ---")
	runLoopsTest()

	fmt.Println("\n--- Named Return Function Test ---")
	fmt.Println(runNamedReturnFunction(1, 2))

	fmt.Println("\n--- Array Slice Initializations Test ---")
	runArraySliceInitializations()

	fmt.Println("\n--- Comma Idiom Test ---")
	runCommaOkIdiom()

	fmt.Println("\n--- Slice Appending Test ---")
	runTestSliceAppending()

	fmt.Println("\n--- Stringifying Object Test ---")
	runTestObjectStringify()

	fmt.Println("\n--- Interface Checks Test ---")
	runTestInterfaceChecks()

	fmt.Println("\n--- Embedding Test ---")
	runTestEmbedding()

	fmt.Println("\n--- Goroutines Test ---")
	runTestGoroutinesWithChannels()

	fmt.Println("\n--- Recover from Panic Test ---")
	runTestRecoverFromPanic()
}

// Allocation test.
type TestAllocationObject struct {
	counter int
}

func runAllocationTest() {
	p := new(TestAllocationObject) // type *TestAllocationObject
	var v TestAllocationObject     // type  TestAllocationObject

	p.counter = 1
	v.counter = 2

	// Objects passed to a function is a fresh copy,
	// hence for `p` it will be updated, whereas `v`
	// will not.
	updatePointer(p)
	update(v)

	fmt.Printf("%+v %T, %+v %T\n", p, p, v, v)
}

func update(t TestAllocationObject) {
	t.counter += 2
}

func updatePointer(t *TestAllocationObject) {
	t.counter += 2
}

// Loops test.
func runLoopsTest() {
	i := 0

	for {
		fmt.Println("yo")

		i++

		if i == 3 {
			break
		}
	}
}

// Named result parameter function.
func runNamedReturnFunction(a int, b int) (sum int) {
	sum = a + b
	return sum
}

// Array/slice initializations.
func runArraySliceInitializations() {
	// These are the same.
	// Apparently `...` is to indicate a slice-ish array.
	array := [3]int{1, 2, 3}
	array2 := [...]int{1, 2, 3}

	fmt.Println(array, array2)
	fmt.Println(arraySum(&array), arraySum(&array2))
}

func arraySum(a *[3]int) (sum int) {
	for _, v := range *a {
		sum += v
	}
	return
}

// Test `value, ok` idiom.
func runCommaOkIdiom() {
	timeZone := map[string]int{
		"UTC": 0 * 60 * 60,
		"EST": -5 * 60 * 60,
		"CST": -6 * 60 * 60,
		"MST": -7 * 60 * 60,
		"PST": -8 * 60 * 60,
	}
	pst, ok := timeZone["PST"]
	wib, ok2 := timeZone["WIB"]

	fmt.Printf("PST: %d, %t\n", pst, ok)
	fmt.Printf("WIB: %d, %t\n", wib, ok2)
}

// Test slice appending.
func runTestSliceAppending() {
	slice1 := []int{1, 2, 3}
	fmt.Println(slice1)

	slice1 = appendSlice(slice1, []int{4, 5, 6})
	fmt.Println(slice1)

	slice1 = appendSlice(slice1, []int{4, 5, 6})
	fmt.Println(slice1)

	// Do the same, but with built-in `append`.
	slice1 = []int{1, 2, 3}
	fmt.Println(slice1)

	slice1 = append(slice1, 4, 5, 6)
	fmt.Println(slice1)

	slice1 = append(slice1, 4, 5, 6)
	fmt.Println(slice1)
}

func appendSlice(slice, data []int) []int {
	l := len(slice)
	if l+len(data) > cap(slice) { // reallocate
		// The docs multiplies this by 2. Let's try not to
		// and see what happens.
		// newSlice := make([]byte, (l+len(data))*2)
		newSlice := make([]int, (l + len(data)))
		// The copy function is predeclared and works for any slice type.
		copy(newSlice, slice)
		slice = newSlice
	}
	// The docs also needs us to "slice" the slice,
	// although the rest of the data is already 0.
	// The state of `slice`, despite the commented below
	// is uncommented, is the same. Probably it's unnecessary?
	// slice = slice[0 : l+len(data)]
	copy(slice[l:], data)
	return slice
}

// Test struct stringify.
// Had help from https://www.geeksforgeeks.org/interfaces-in-golang/.
// The idea is, interface can be "merged" into a struct,
// but the methods have to be implemented.
type TestObjectStringifyInterface interface {
	String() string
}

type TestObjectStringify struct {
	counter int
}

func (c TestObjectStringify) String() string {
	return fmt.Sprintf("The value is %d", c.counter)
}

func runTestObjectStringify() {
	var c TestObjectStringifyInterface = TestObjectStringify{counter: 2}

	fmt.Println(c)
}

// Test interface checks.
func runTestInterfaceChecks() {
	var m interface{} = (*json.RawMessage)(nil)

	if _, ok := m.(json.Marshaler); ok {
		fmt.Printf("value %v of type %T implements json.Marshaler\n", m, m)
	}
}

// Test embedding.
// Had some more help from https://www.geeksforgeeks.org/embedding-interfaces-in-golang/.
// It looks like for interfaces, we can "assign" values
// to it and execute stuff.
type EmbedInterfaceHello interface {
	Hello()
}

type EmbedInterfaceWorld interface {
	World()
}

type EmbedInterfaceHelloWorld interface {
	EmbedInterfaceHello
	EmbedInterfaceWorld
}

type EmbedMain struct {
	id string
}

func (em EmbedMain) Hello() {
	fmt.Println("Hello", em.id)
}

func (em EmbedMain) World() {
	fmt.Println("World", em.id)
}

func runTestEmbedding() {
	values := EmbedMain{id: "1234"}

	var emi EmbedInterfaceHelloWorld = values
	emi.Hello()
	emi.World()
}

// Test channeling with goroutines (again).
func runTestGoroutinesWithChannels() {
	// Use cases of channels:
	// 1. Channels of channels.
	// 2. Parallelization (not showcased in this file).
	ci := make(chan int)
	ci2 := make(chan int)

	go func() {
		ci2 <- <-ci * 2
	}()

	ci <- 123
	fmt.Println(<-ci2)
}

// Test recovering from a panic.
func runTestRecoverFromPanic() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
		}
	}()

	panic("test")
}
