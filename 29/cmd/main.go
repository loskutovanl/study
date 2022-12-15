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
	askNum()
	fmt.Println("Программа завершена.")
}

func askNum() {
	fmt.Print("Ввод: ")
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	for scanner.Scan() {
		input := scanner.Text()
		if input == "стоп" {
			break
		}
		if num, err := strconv.Atoi(input); err != nil {
			fc := getNum(num)
			wg.Add(1)
			sc := squareNum(fc, &wg)
			wg.Add(1)
			multiplyByTwo(sc, &wg)
		}
	}

	wg.Wait()
}

func getNum(num int) chan int {
	firstChan := make(chan int)
	firstChan <- num
	return firstChan
}

func squareNum(squareChan chan int, wg *sync.WaitGroup) chan int {
	defer wg.Done()
	secondChan := make(chan int)
	go func() {
		num := <-squareChan
		result := int(math.Sqrt(float64(num)))
		fmt.Println("Квадрат:", result)
		secondChan <- result
	}()
	return secondChan
}

func multiplyByTwo(multiplyChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		num := <-multiplyChan
		result := num * 2
		fmt.Println("Произведение:", result)
	}()
}
