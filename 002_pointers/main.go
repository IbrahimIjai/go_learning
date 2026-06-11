package main
import "fmt"

//pointers are variables that store the memory address of another variable. They allow us to indirectly access and manipulate the value stored at that address. In Go, we can declare a pointer using the * operator and assign it the address of a variable using the & operator.

func main() {	
	var x int = 10

	//stores address of x in ptr
	var ptr *int = &x

	fmt.Printf("Value of x: %d\n", x)
	fmt.Printf("Address of x: %p\n", &x)
	fmt.Printf("Value pointed to by ptr: %d\n", *ptr)


	// using slices

	var slice = []int{1, 2, 3, 4, 5}
	var sliceCopy = slice // this creates a copy of the slice header, but both slice and sliceCopy point to the same underlying array
	sliceCopy[0] = 10 // modifying sliceCopy will also modify the original slice
	fmt.Printf("Original slice: %v\n", slice)
	fmt.Printf("Copied slice: %v\n", sliceCopy)


	// using in functions

	var thing = [5]float64{1.0, 2.0, 3.0, 4.0, 5.0}
	fmt.Printf("Original array: %v\n", thing)
	var squaredThing = square(thing) 
	fmt.Printf("Squared array: %v\n", squaredThing)
	fmt.Printf("Original array after function call: %v\n", thing) // the original array remains unchanged

}


// function square (thin)

func square(thing2 [5]float64) [5]float64 {
	for i := range thing2{
		thing2[i] = thing2[i] * thing2[i]	
	}
	return thing2

}