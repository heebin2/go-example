package main

import (
	"fmt"
	"strconv"

	"github.com/eiannone/keyboard"
)

func main() {

	index := make(map[rune]int)
	index['k'] = 0
	index['l'] = 1
	index[';'] = 2
	index['\''] = 3

	index['1'] = 0
	index['2'] = 1
	index['3'] = 2

	index['4'] = 0
	index['5'] = 1
	index['6'] = 2

	index['m'] = 0
	index[','] = 1
	index['.'] = 2
	index['/'] = 3

	updown := make(map[rune]bool)
	updown['k'] = true
	updown['l'] = true
	updown[';'] = true
	updown['\''] = true

	updown['m'] = false
	updown[','] = false
	updown['.'] = false
	updown['/'] = false

	updown['1'] = false
	updown['2'] = false
	updown['3'] = false

	updown['4'] = true
	updown['5'] = true
	updown['6'] = true

	ary := make([]int, 4)

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Add   : [k] [l] [;] [']")
	fmt.Println("      : [4] [5] [6]")
	fmt.Println("Sub   : [m] [,] [.] [/]")
	fmt.Println("      : [1] [2] [3]")
	fmt.Println("Quit  : ESC")
	fmt.Println("Reset : Delete ")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Key == keyboard.KeyEsc {
			break
		}

		if event.Key == keyboard.KeyDelete {
			ary = []int{0, 0, 0, 0}
			fmt.Println("[  0  0  0  0  ] delete")
			continue
		}

		i, exist := index[event.Rune]
		if !exist {
			fmt.Println("not found key : ", event.Rune)
			continue
		}
		ud := updown[event.Rune]
		if ud {
			ary[i] += 1
		} else if ary[i] > 0 {
			ary[i] -= 1
		}
		noti := "["
		for j := range ary {
			if i == j {
				if ud {
					noti += " +"
				} else {
					noti += " -"
				}
			} else {
				noti += "  "
			}
			noti += strconv.Itoa(ary[j])
		}
		noti += "  ]"
		fmt.Println(noti)
	}
}
