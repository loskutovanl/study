package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	firstChan := generateChan(&wg)
	secondChan := squareNum(firstChan, &wg)
	multiplyByTwo(secondChan, &wg)

	wg.Wait()
	fmt.Println("Программа завершена.")
}

func generateChan(wg *sync.WaitGroup) <-chan int {
	out := make(chan int)
	wg.Add(1)
	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()

		for {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input := scanner.Text()

			if input == "стоп" {
				break
			}

			num, err := strconv.Atoi(input)
			if err == nil {
				out <- num
			} else {
				fmt.Printf("Ввод %s не может быть преобразован в число.\n", input)
				continue
			}
		}
	}()
	return out
}

func squareNum(in <-chan int, wg *sync.WaitGroup) <-chan int {
	out := make(chan int, 5)
	wg.Add(1)
	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()

		for n := range in {
			result := int(math.Pow(float64(n), 2))
			fmt.Println("Квадрат:", result)
			out <- result
		}
	}()
	return out
}

func multiplyByTwo(in <-chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		for n := range in {
			result := 2 * n
			fmt.Println("Произведение:", result)
		}
	}()
}
