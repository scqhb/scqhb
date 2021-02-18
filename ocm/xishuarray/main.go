package main

import "fmt"

type S0 struct {
	row   int
	col   int
	value int
}

func main() {

	var arrar [11][11]int
	arrar[0][2] = 1
	arrar[8][9] = 2
	for _, v := range arrar {
		for _, v2 := range v {
			fmt.Printf("%d \t", v2)
		}
		fmt.Println()
	}
	var SS []S0
	SS = append(SS, S0{11, 11, 0})
	for i, v := range arrar {
		for j, v2 := range v {
			if v2 != 0 {
				s1 := S0{
					row:   i,
					col:   j,
					value: v2,
				}
				SS = append(SS, s1)

			}
		}
	}
	for _, v := range SS {
		fmt.Printf("%d\t%d\t%d\n", v.row, v.col, v.value)
	}

	var s [11][11]int
	for i, v1 := range SS {
		if i != 0 {
			s[v1.row][v1.row] = v1.value
		}

	}

	fmt.Printf("############\n")
	for _, v3 := range s {
		for _, v4 := range v3 {
			fmt.Printf("%d\t", v4)
		}
		fmt.Println()
	}

}
