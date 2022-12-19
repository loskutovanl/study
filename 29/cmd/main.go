package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {

	for {
		fmt.Print("Ввод: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		if input == "стоп" {
			break
		}

		num, err := strconv.Atoi(input)
		if err == nil {
			firstChan := generateChan(num)
			secondChan := squareNum(firstChan)
			thirdChan := multiplyByTwo(secondChan)
			fmt.Println("Произведение:", <-thirdChan)
		}
	}

	fmt.Println("Программа завершена.")
}

func generateChan(num int) <-chan int {
	out := make(chan int)
	go func() {
		out <- num
		close(out)
	}()
	return out
}

func squareNum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			result := int(math.Pow(float64(n), 2))
			fmt.Println("Квадрат:", result)
			out <- result
		}
		close(out)
	}()
	return out
}

func multiplyByTwo(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			result := 2 * n
			out <- result
		}
		close(out)
	}()
	return out
}
