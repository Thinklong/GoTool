package main

import (
	"fmt"
)

func main() {
	fmt.Printf("thinklong")
	monthdays := map[string]int{
		"Jan": 31, "Feb": 28, "Mar": 31,
		"Apr": 30, "May": 31, "Jun": 30,
		"Jul": 31, "Aug": 31, "Sep": 30,
		"Oct": 31, "Nov": 30, "Dec": 31,
	}
	key := "Jan"
	if _, ok := monthdays[key]; ok {
		fmt.Println(true)
	} else {
		fmt.Println(false)
	}

	ok := is_set(monthdays, "Jan")
	fmt.Println(ok)
	a := "A"
	for i := 1; i < 10; i++ {
		for j := i; j > 0; j-- {

			fmt.Print(a)
		}
		fmt.Print("\n")
	}

	avg_list := []float64{111.323, 222.23422, 2342.232221}
	switch len(avg_list) {
	case 0:
		fmt.Println(0)
	default:
		num := 0.0
		for _, value := range avg_list {
			num += value
		}
		avg := num / float64(len(avg_list))
		fmt.Printf("%f", avg)
	}

	fmt.Println("abc")
}

func is_set(data map[string]int, _key string) bool {

	if _, ok := data[_key]; ok {
		return true
	}
	return false
}
