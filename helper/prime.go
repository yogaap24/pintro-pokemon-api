package helper

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"time"
)

const (
	MinPrime = 1
	MaxPrime = 20
)

var (
	PrimeNumbers     = []int{2, 3, 5, 7, 11, 13, 17, 19}
	NonPrimeNumbers  = []int{1, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20}
	DefaultAttempts  = 0
	DefaultThreshold = GenerateThreshold() // Set initial threshold
)

type PrimeError struct {
	Number int
}

func (e *PrimeError) Error() string {
	return fmt.Sprintf("Not a prime number (%d)", e.Number)
}

type PrimeGenerator struct {
	primes            []int
	nonPrimes         []int
	usedPrimes        map[int]bool
	attempts          int
	attemptsThreshold int
	mu                sync.Mutex
}

func NewPrimeGenerator() *PrimeGenerator {
	return &PrimeGenerator{
		primes:            append([]int(nil), PrimeNumbers...),
		nonPrimes:         append([]int(nil), NonPrimeNumbers...),
		usedPrimes:        make(map[int]bool),
		attempts:          DefaultAttempts,
		attemptsThreshold: DefaultThreshold,
	}
}

func (pg *PrimeGenerator) GetUniquePrime() (int, error) {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	var number int
	var err error

	log.Debug().Msgf("Attempts: %d, Threshold: %d", pg.attempts, pg.attemptsThreshold)

	if pg.attempts < pg.attemptsThreshold {
		number, err = pg.getRandomNumberFromSlice(pg.nonPrimes)
		if err != nil {
			return 0, err
		}
		pg.attempts++
		DefaultAttempts = pg.attempts
		return number, &PrimeError{Number: number}
	}

	if len(pg.primes) == 0 {
		pg.primes = append([]int(nil), PrimeNumbers...)
		pg.usedPrimes = make(map[int]bool)
	}

	number, err = pg.getRandomNumberFromSlice(pg.primes)
	if err != nil {
		return 0, err
	}

	pg.usedPrimes[number] = true
	pg.attempts = 0
	DefaultAttempts = 0
	pg.attemptsThreshold = GenerateThreshold() // Update threshold after successful prime generation

	return number, nil
}

func (pg *PrimeGenerator) getRandomNumberFromSlice(slice []int) (int, error) {
	if len(slice) == 0 {
		return 0, fmt.Errorf("No numbers available")
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(slice))
	number := slice[index]

	slice = append(slice[:index], slice[index+1:]...)
	return number, nil
}

func IsPrime(n int) error {
	if n < MinPrime || n > MaxPrime {
		return fmt.Errorf("Number out of valid range (%d)", n)
	}

	for _, prime := range PrimeNumbers {
		if n == prime {
			return nil
		}
	}

	return &PrimeError{Number: n}
}

func GenerateThreshold() int {
	rand.Seed(time.Now().UnixNano())
	return generatePatternedThreshold()
}

func generatePatternedThreshold() int {
	var thresholds []int
	lastValue := rand.Intn(MaxThreshold-MinThreshold+1) + MinThreshold
	thresholds = append(thresholds, lastValue)

	for i := 0; i < 5; i++ {
		step := rand.Intn(15) + 5

		if i%2 == 0 {
			newValue := lastValue + step
			if newValue > MaxThreshold {
				newValue = MaxThreshold
			}
			thresholds = append(thresholds, newValue)
			lastValue = newValue
		} else {
			newValue := lastValue - step
			if newValue < MinThreshold {
				newValue = MinThreshold
			}
			thresholds = append(thresholds, newValue)
			lastValue = newValue
		}
	}

	return thresholds[len(thresholds)-1]
}

const (
	MinThreshold = 10
	MaxThreshold = 100
)
