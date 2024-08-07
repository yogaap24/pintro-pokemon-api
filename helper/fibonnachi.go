package helper

import (
	"strconv"
	"strings"
	"sync"
)

var (
	fibMap = make(map[string]int)
	mu     sync.Mutex
)

func CalculateFibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func GetNextFibonacciValue(key string) int {
	mu.Lock()
	defer mu.Unlock()

	key = strings.TrimSpace(key)
	if key == "" {
		return 0
	}

	val, exists := fibMap[key]
	if !exists {
		val = 0
	}

	fibValue := CalculateFibonacci(val)
	fibMap[key] = val + 1

	return fibValue
}

func GenerateNickName(baseName string, fibValue int) string {
	baseName = strings.TrimSpace(baseName)
	baseName = strings.ReplaceAll(baseName, " ", "")

	if strings.Contains(baseName, "-") {
		parts := strings.Split(baseName, "-")
		if len(parts) > 1 {
			baseName = strings.Join(parts[:len(parts)-1], "-")
		}
	}

	return baseName + "-" + strconv.Itoa(fibValue)
}
