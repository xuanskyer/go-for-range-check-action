package third

import "fmt"

func Test() {
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
