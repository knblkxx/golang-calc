package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func delsp(nums string) string {
	var out []string
	for _, c := range nums {
		if c != ' ' {
			out = append(out, string(c))
		}
	}
	return strings.Join(out, "")
}
func CalcExpression(expression string) (string, error) {
	expression = delsp(expression)
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
	n := 0
	moves := 0
	prioritet := -1
	proshloe := -1
	for i, el := range nums {
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
				return 0, errors.New("Закывающих скобок больше")
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
	if open > 0 {
		return 0, errors.New("Не закыта скобка")
	}
	out, err := CalcExpression(expression)
	if err != nil {
		return 0, err
	}
	out1, _ := strconv.ParseFloat(out, 64)
	return out1, nil
}

//Тесты ниже
func test() {
	fmt.Println(Calc("2 + 2 + (2 + (2 + (2 + 2)))"))
	fmt.Println(Calc("2+3"))
	fmt.Println(Calc("(2+7)*3"))
	fmt.Println(Calc("2+2*2"))
	fmt.Println(Calc("1+1*"))
}
