package main

import (
	"fmt"
	"math"
)

func main() {
	var low, high int
	fmt.Print("Enter the lower bound: ")
	fmt.Scan(&low)
	fmt.Print("Enter the upper bound: ")
	fmt.Scan(&high)

	guessedNumbers := make([]int, 0)
	attempts := 0
	currentGuess := (low + high) / 2

	for {
		attempts++
		fmt.Printf("Is your number %d? (higher/lower/correct): ", currentGuess)
		var response string
		fmt.Scan(&response)

		if containsNumber(guessedNumbers, currentGuess) {
			fmt.Println("I've already guessed this number before. You must be lying or mistaken.")
			return
		}

		guessedNumbers = append(guessedNumbers, currentGuess)
		fmt.Println("COM guessed answers: ", guessedNumbers)

		switch response {
		case "higher":
			if currentGuess >= high {
				fmt.Println("You said 'higher', but there are no higher numbers in the range. You must be lying.")
				return
			}
			low = currentGuess + 1
		case "lower":
			if currentGuess <= low {
				fmt.Println("You said 'lower', but there are no lower numbers in the range. You must be lying.")
				return
			}
			high = currentGuess - 1
		case "correct":
			fmt.Printf("Hurray! I guessed your number in %d attempts.\n", attempts)
			return
		default:
			fmt.Println("Invalid response. Please enter 'higher', 'lower', or 'correct'.")
			attempts-- // Don't count invalid responses
			continue
		}

		if low > high {
			fmt.Println("The range is now invalid. You must have made a mistake or lied.")
			return
		}

		currentGuess = (low + high) / 2

		possibleNumbers := high - low + 1
		if possibleNumbers <= 0 {
			fmt.Println("There are no more possible numbers. You must have made a mistake or lied.")
			return
		}
		numberx := math.Log2(float64(high-low+2)) * 2
		fmt.Println(numberx)
		if float64(attempts) > numberx {
			fmt.Println("I've made too many guesses. You must be lying.")
			return
		}
	}
}

func containsNumber(numbers []int, num int) bool {
	for _, n := range numbers {
		if n == num {
			return true
		}
	}
	return false
}
