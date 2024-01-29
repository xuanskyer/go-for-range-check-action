package biz

import "fmt"

func Test() {
	arr := make([]string, 0)
	for i := 0; i < 10; i++ {
		if i < 5 {
			continue
		}
		for _, item := range arr {
			fmt.Println(item)
		}

	}
}

func Test2() {
	arr := make([]string, 0)
	for i := 0; i < 10; i++ {
		if i < 5 {
			continue
		}
		for j := 0; j < 10; j++ {
			if j > 6 {
				continue
			}
			for _, item := range arr {
				fmt.Println(item)
			}
		}
	}
}

func Test3() {
	arr := make([]string, 0)
	for i := 0; i < 10; i++ {
		if i > 5 {
			continue
		}
		for j := 0; j < 10; j++ {
			if j > 6 {
				continue
			}
			for _, item := range arr {
				fmt.Println(item)
			}
			for i2 := 0; i2 < 10; i2++ {
				if i2 > 5 {
					continue
				}
				for j2 := 0; j2 < 10; j2++ {
					if j2 > 6 {
						continue
					}
					for _, item := range arr {
						fmt.Println(item)
					}
				}
			}
		}
	}
}
