package engines

import (
	"fmt"
)

type gasEngine struct {
	kpg     uint8
	gallons uint8
	owner   owner
	company
}
type electricEngine struct {
	kpkwh uint8
	kwh   uint8
	owner owner
	company
}
type engine interface {
	milesLeft() uint8
}

type owner struct {
	name string
}

type company struct {
	name string
}

func (e gasEngine) milesLeft() uint8 {
	return e.kpg * e.gallons
}
func (e electricEngine) milesLeft() uint8 {
	return e.kpkwh * e.kwh
}
func canMakeIt(e engine, miles uint8) {
	fmt.Println("test")
	// return true
}

func Test() {
	var gascar = gasEngine{kpg: 10, gallons: 120}
	fmt.Println("gallons: %v,kilometers per gallons: %v, owner name: %v, company: %v, miles left:%v", gascar.gallons, gascar.kpg, gascar.owner.name, gascar.name, gascar.milesLeft())
	canMakeIt(gascar, 100)
	// var electricCar = gasEngine{kpg: 10, gallons: 120}
	// fmt.Println("gallons: %v,kilometers per gallons: %v, owner name: %v, company: %v, miles left:%v", gascar.gallons, gascar.kpg, gascar.owner.name, gascar.name, electricCar.milesLeft())

}
