package calculator

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"time"
)

//	type ExpressionParallel struct {
//		Expression string
//		IndexB     int
//		IndexE     int
//	}
//
//	type ExpressionParallelResult struct {
//		Result float64
//		IndexB int
//		IndexE int
//	}
type MoveType struct {
	Type      string
	Index     int
	Prioritet int
}
type ExpressionOutput struct {
	Num   string
	Index int
}
type ExpressionInput struct {
	Move MoveType
	Chan chan ExpressionOutput
}

var TIME_ADDITION_MS int = 0
var TIME_SUBTRACTION_MS int = 0
var TIME_MULTIPLICATIONS_MS int = 0
var TIME_DIVISIONS_MS int = 0

func NoSpaces(nums string) string {
	var out []string
	for _, c := range nums {
		if c != ' ' {
			out = append(out, string(c))
		}
	}
	return strings.Join(out, "")
}
func FinalCalc(input ExpressionInput, expression []string) {
	a, _ := strconv.ParseFloat(expression[input.Move.Index-1], 64)
	b, _ := strconv.ParseFloat(expression[input.Move.Index+1], 64)
	result := 0.0
	if input.Move.Type == "+" {
		timer := time.NewTimer(time.Duration(TIME_ADDITION_MS) * time.Millisecond)
		result = a + b
		<-timer.C
	} else if input.Move.Type == "-" {
		timer := time.NewTimer(time.Duration(TIME_SUBTRACTION_MS) * time.Millisecond)
		result = a - b
		<-timer.C
	} else if input.Move.Type == "*" {
		timer := time.NewTimer(time.Duration(TIME_MULTIPLICATIONS_MS) * time.Millisecond)
		result = a * b
		<-timer.C
	} else if input.Move.Type == "/" {
		if b != 0 {
			timer := time.NewTimer(time.Duration(TIME_DIVISIONS_MS) * time.Millisecond)
			result = a / b
			<-timer.C
		} else {
			input.Chan <- ExpressionOutput{
				"error",
				input.Move.Index,
			}
			return
		}
	}
	input.Chan <- ExpressionOutput{
		fmt.Sprintf("%f", result),
		input.Move.Index,
	}
}
func CalcExpression(expression string) (string, error) {
	//fmt.Println(expression)
	expression = NoSpaces(expression)
	expression = strings.Replace(expression, "/", " / ", -1)
	expression = strings.Replace(expression, "*", " * ", -1)
	expression = strings.Replace(expression, "+", " + ", -1)
	expression = strings.Replace(expression, "-", " - ", -1)
	nums := strings.Split(expression, " ")
	if strings.Contains("+-/*", nums[0]) {
		return "", errors.New("Невалидное выражение")
	}
	if strings.Contains("+-/*", nums[len(nums)-1]) {
		return "", errors.New("Невалидное выражение")
	}
	if len(nums) == 1 {
		return expression, nil
	}
	//fmt.Println(nums)
	n := 0
	moves := 0
	prioritet := -1
	proshloe := -1
	movesSlice := make([]MoveType, 0)
	for i, el := range nums {
		//fmt.Println(el)
		if strings.Contains("+-/*", el) {
			if proshloe != -1 {
				if proshloe == 1 {
					return "", errors.New("Невалидное выражение")
				}
			}
			moves++
			proshloe = 1
			if prioritet == -1 && strings.Contains("/*", el) {
				prioritet = i
			}
			movesSlice = append(movesSlice, MoveType{
				Type:      el,
				Index:     i,
				Prioritet: prioritet,
			})
		} else {
			for _, c := range el {
				if !strings.Contains("1234567890.", string(c)) {
					return "", errors.New("Невалидные символы")
				}
			}
			if proshloe != -1 {
				if proshloe == 0 {
					return "", errors.New("Невалидное выражение")
				}
			}
			proshloe = 0
			n++
		}
	}
	if len(movesSlice) > 1 {
		//fmt.Println(movesSlice)
		parralelMoves := make([]MoveType, 0)
		i := 0
		for {
			if i >= len(movesSlice) {
				break
			}
			if i == 0 {
				if movesSlice[i].Prioritet >= movesSlice[i+1].Prioritet {
					parralelMoves = append(parralelMoves, movesSlice[i])
					i += 2
				} else {
					i++
				}
			} else if i == len(movesSlice)-1 {
				if movesSlice[i].Prioritet >= movesSlice[i-1].Prioritet {
					parralelMoves = append(parralelMoves, movesSlice[i])
					i += 2
				} else {
					i++
				}
			} else {
				if movesSlice[i-1].Prioritet <= movesSlice[i].Prioritet && movesSlice[i].Prioritet >= movesSlice[i+1].Prioritet {
					parralelMoves = append(parralelMoves, movesSlice[i])
					i += 2
				} else {
					i++
				}
			}
		}
		//fmt.Println(parralelMoves)
		ch := make(chan ExpressionOutput)
		for _, mov := range parralelMoves {
			go FinalCalc(ExpressionInput{
				mov,
				ch,
			}, nums)
		}
		for i := 0; i < len(parralelMoves); i++ {
			select {
			case x, ok := <-ch:
				if ok {
					if x.Num == "error" {
						return "", errors.New("Деление на ноль")
					}
					nums[x.Index] = x.Num
				}
			}
		}
		new_nums := make([]string, 0)
		if strings.Contains("+-/*", nums[1]) {
			new_nums = append(new_nums, nums[0])
		}
		for i := 1; i+1 < len(nums); i++ {
			if !strings.Contains("+-/*", nums[i-1]) && !strings.Contains("+-/*", nums[i+1]) {
				new_nums = append(new_nums, nums[i])
			} else if strings.Contains("+-/*", nums[i-1]) && strings.Contains("+-/*", nums[i+1]) {
				new_nums = append(new_nums, nums[i])
			} else if !strings.Contains("+-/*", nums[i]) && !strings.Contains("+-/*", nums[i-1]) && !strings.Contains("+-/*", nums[i+1]) {
				new_nums = append(new_nums, nums[i])
			}
		}
		if strings.Contains("+-/*", nums[len(nums)-2]) {
			new_nums = append(new_nums, nums[len(nums)-1])
		}
		//fmt.Println(nums, new_nums)
		return CalcExpression(strings.Join(new_nums, " "))
	}
	if n-moves != 1 {
		return "", errors.New("Невалидное выражание")
	}
	out := 0.0
	if prioritet != -1 {
		a, _ := strconv.ParseFloat(nums[prioritet-1], 64)
		b, _ := strconv.ParseFloat(nums[prioritet+1], 64)
		if nums[prioritet] == "*" {
			out = a * b
		} else {
			out = a / b
		}
		if len(nums)-2 != 1 {
			return CalcExpression(fmt.Sprintf("%s%f%s", strings.Join(nums[:prioritet-1], ""), out, strings.Join(nums[prioritet+2:], "")))
		}
	} else {
		a, _ := strconv.ParseFloat(nums[0], 64)
		b, _ := strconv.ParseFloat(nums[2], 64)
		if nums[1] == "+" {
			out = a + b
		} else {
			out = a - b
		}
		if len(nums)-2 != 1 {
			return CalcExpression(fmt.Sprintf("%f%s", out, strings.Join(nums[3:], "")))
		}
	}
	return fmt.Sprintf("%f", out), nil
}

func Calc(expression string) (float64, error) {
	//fmt.Println(expression)
	open := 0
	begin := -1
	end := -1
	for i, c := range expression {
		if c == '(' {
			open++
			begin = i
		} else if c == ')' {
			open--
			end = i
			if open == -1 {
				return 0, errors.New("Закрывается никогда не открытая скобка")
			}
			if end-begin == 1 {
				return 0, errors.New("Пустое выражение в скобках")
			}
			res, err := CalcExpression(expression[begin+1 : end])
			if err != nil {
				return 0, err
			}
			return Calc(expression[:begin] + res + expression[end+1:])
		}
	}
	//opened := make([]int, 0)
	//expressionsParallel := make([]ExpressionParallel, 0)
	//for i, c := range expression {
	//	if c == '(' {
	//		opened = append(opened, i)
	//	} else if c == ')' {
	//		if len(opened) == 0 {
	//			return 0, errors.New("Закрывается никогда не открытая скобка")
	//		}
	//		if i-opened[len(opened)-1] == 1 {
	//			return 0, errors.New("Пустое выражение в скобках")
	//		}
	//		expressionsParallel = append(expressionsParallel,
	//			ExpressionParallel{
	//				Expression: expression[opened[len(opened)]-1 : i],
	//				IndexB:     opened[len(opened)-1],
	//				IndexE:     i,
	//			})
	//		opened = opened[:len(opened)-1]
	//	}
	//}
	//ch := make(chan ExpressionParallelResult)

	if open > 0 {
		return 0, errors.New("Скобка открылась, но так и не закрылась")
	}
	out, err := CalcExpression(expression)
	if err != nil {
		return 0, err
	}
	out1, _ := strconv.ParseFloat(out, 64)
	return out1, nil
}

func Initial() {
	godotenv.Load("../.env")
	value := os.Getenv("TIME_ADDITION_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_ADDITION_MS")
			os.Exit(0)
		}
		TIME_ADDITION_MS = intvalue
	} else {
		TIME_ADDITION_MS = 0
	}

	value = os.Getenv("TIME_SUBTRACTION_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_SUBTRACTION_MS")
			os.Exit(0)
		}
		TIME_SUBTRACTION_MS = intvalue
	} else {
		TIME_SUBTRACTION_MS = 0
	}

	value = os.Getenv("TIME_MULTIPLICATIONS_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_MULTIPLICATIONS_MS")
			os.Exit(0)
		}
		TIME_MULTIPLICATIONS_MS = intvalue
	} else {
		TIME_MULTIPLICATIONS_MS = 0
	}

	value = os.Getenv("TIME_DIVISIONS_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_DIVISIONS_MS")
			os.Exit(0)
		}
		TIME_DIVISIONS_MS = intvalue
	} else {
		TIME_DIVISIONS_MS = 0
	}
	fmt.Printf("TIME_ADDITION_MS: %d\n", TIME_ADDITION_MS)
	fmt.Printf("TIME_SUBTRACTION_MS: %d\n", TIME_SUBTRACTION_MS)
	fmt.Printf("TIME_MULTIPLICATIONS_MS: %d\n", TIME_MULTIPLICATIONS_MS)
	fmt.Printf("TIME_DIVISIONS_MS: %d\n", TIME_DIVISIONS_MS)

}

func test() {
	fmt.Println(Calc("2 + 2 + 2 + 2 + 2 + 2 + (2 + (2 + (2 + 2)))"))
	fmt.Println(Calc("1+1"))
	fmt.Println(Calc("(2+2)*2"))
	fmt.Println(Calc("2+2*2"))
	fmt.Println(Calc("1+1*"))
}
