package main

import "fmt"


type Employee struct {
	name string
	age uint8
	rank uint8
}

type gasEngine struct{
	mpg uint8
	gallons uint8
	Owner
}

type electricEngine struct{
	mpkwh uint8
	kwh uint8
	Owner
}

//compoition example
type Owner struct {
	Name string `json:"name"` // what does this tag do? answer: it is used to specify the name of the field when it is serialized to JSON. In this case, the field will be serialized as "name" instead of "Name".
	age uint8
}

// methods for gasEngine

func (e electricEngine) milesLeft() uint8 {
	return e.kwh * e.mpkwh
}

func (g gasEngine) milesLeft() uint8 {
	return g.gallons * g.mpg
}

//Interfaces -- contracts/method signatures that a type must implement to satisfy the interface
type Engine interface {
	milesLeft() uint8
}

func main (){
	var myEngine gasEngine = gasEngine{mpg: 30, gallons: 10}
	fmt.Printf("My engine has %d miles per gallon and %d gallons of fuel.\n", myEngine.mpg, myEngine.gallons)

	employee1 := Employee{name: "Ibrahim", age: 30, rank: 5}
	fmt.Printf("Employee: %s, Age: %d, Rank: %d\n", employee1.name, employee1.age, employee1.rank)

}