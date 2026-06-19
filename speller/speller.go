//go:build !solution

package speller

import "fmt"

var basicNumbersSpeller = map[int64]string{
	0: "zero", 1: "one", 2: "two", 3: "three", 4: "four", 5: "five",
	6: "six", 7: "seven", 8: "eight", 9: "nine", 10: "ten",
	11: "eleven", 12: "twelve", 13: "thirteen", 14: "fourteen", 15: "fifteen",
	16: "sixteen", 17: "seventeen", 18: "eighteen", 19: "nineteen", 20: "twenty",
	30: "thirty", 40: "forty", 50: "fifty", 60: "sixty", 70: "seventy",
	80: "eighty", 90: "ninety", 100: "hundred", 1000: "thousand",
	1_000_000: "million", 1_000_000_000: "billion",
}

func spellOnes(n int64) string {
	if n < 0 || n > 9 {
		return ""
	}
	return basicNumbersSpeller[n]
}

func spellTens(n int64) string {
	if n < 10 {
		return spellOnes(n)
	}

	if n >= 100 {
		return ""
	}

	if n%10 == 0 {
		return basicNumbersSpeller[n]
	}

	if n < 20 {
		return basicNumbersSpeller[n]
	}

	f := n / 10
	s := n % 10
	return fmt.Sprintf("%s-%s", basicNumbersSpeller[f*10], basicNumbersSpeller[s])
}

func spellHundreds(n int64) string {
	if n < 100 {
		return spellTens(n)
	}
	if n >= 1000 {
		return ""
	}
	f := n / 100
	s := n % 100
	if s == 0 {
		return fmt.Sprintf("%s hundred", basicNumbersSpeller[f])
	}
	return fmt.Sprintf("%s hundred %s", basicNumbersSpeller[f], spellTens(s))
}

func spellThousands(n int64) string {
	if n < 1000 {
		return spellHundreds(n)
	}
	if n >= 1_000_000 {
		return ""
	}

	f := n / 1000
	s := n % 1000
	if s == 0 {
		return fmt.Sprintf("%s thousand", spellHundreds(f))
	}
	return fmt.Sprintf("%s thousand %s", spellHundreds(f), spellHundreds(s))
}

func spellMillions(n int64) string {
	if n < 1_000_000 {
		return spellThousands(n)
	}
	if n >= 1_000_000_000 {
		return ""
	}
	f := n / 1_000_000
	s := n % 1_000_000
	if s == 0 {
		return fmt.Sprintf("%s million", spellHundreds(f))
	}
	return fmt.Sprintf("%s million %s", spellHundreds(f), spellThousands(s))
}

func spellBillions(n int64) string {
	if n < 1_000_000_000 {
		return spellMillions(n)
	}
	if n >= 1_000_000_000_000 {
		return ""
	}
	f := n / 1e9
	s := n % 1e9
	if s == 0 {
		return fmt.Sprintf("%s billion", spellHundreds(f))
	}
	return fmt.Sprintf("%s billion %s", spellHundreds(f), spellMillions(s))
}

func Spell(n int64) (s string) {
	if n < 0 {
		return "minus " + spellBillions(-n)
	}
	return spellBillions(n)
}
