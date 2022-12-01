package main

import "fmt"

//choose == 1, start cmd
//choose == 2, start App
func Init(choose int) {
	switch choose {
	case 1:
		commandLine()
	case 2:
		GUI()
	default:
		fmt.Println("wrong number!")
	}

}

func main() {
	fmt.Println("Please enter the number to choose starting way:(1.CMD, 2.GUI)")
	var num int
	fmt.Scanf("%d\n", &num)
	Init(num)
}
