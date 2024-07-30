package main

import (
	"errors"
	"fmt"
	"go_helloworld/engines"
	"strings"
	"unicode/utf8"

	"rsc.io/quote"
)

func main() {
	engines.Test()

}

func Basics() {
	fmt.Println("hello")
	fmt.Println(utf8.RuneCountInString("hello"))
	fmt.Println(quote.Go())
	var char rune = 'a'
	var str string = "134234"

	fmt.Printf("%c \t %v", char, str)
	print("%c \t %v", char, str)
	var _, _, err = divide(1, 0)
	print(err.Error())
	var slice = []int{1, 2, 3, 4, 5}
	var slice2 = []int{1, 2, 3, 4, 5}
	print("%v", slice)
	slice = append(slice, slice2...)
	print("%v", slice)
	var slice3 = make([]int, 5, 8)
	print("%v", slice3)
	for i, value := range slice {
		print("%v,%v ", value, i)
	}
	var _string = "hello"
	print("\n%v, %T", _string[0], _string[0])
	for _, v := range _string {
		print("\n%v %v,", v, v)
	}
	var str_slice = []string{"h", "w", "e", "l", "l", "o"}
	var strBuilder strings.Builder
	for i := range str_slice {
		strBuilder.WriteString(str_slice[i])
	}
	print("\n %v", strBuilder.String())
}

func print(value string, args ...interface{}) {
	fmt.Printf(value, args...)
}
func divide(numerator int, denominator int) (int, int, error) {
	var er error
	if denominator == 0 {
		er = errors.New("Denominator cannot be zero")
		return numerator, denominator, er
	}
	fmt.Println(numerator / denominator)
	return numerator, denominator, er
}
