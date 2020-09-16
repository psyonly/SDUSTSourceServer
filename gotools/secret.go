package gotools

import (
	"fmt"
	"golang.org/x/crypto/md4"
	"io"
)

const (
	salt1 = "kcats"
	salt2 = "stmtn"
	salt3 = "myktm"
)

var NUMHASH = map[rune]string{
	'0': "h0",
	'1': "1Q",
	'5': "7T",
	'6': "v88",
	'2': "e1",
	'7': "C5",
	'8': "w1",
	'3': "NM",
	'4': "s1",
	'9': "KU",
}

var CHARHASH = map[rune]string{
	'a': "AA", 'z': "6c", 'M': "yq", 'N': "h2", 'Y': "un",
	'b': "2t", 'y': "mn", 'L': "go",
	'c': "Pu", 'x': "ki", 'K': "km", 'O': "cl", 'X': "r3", 'Z': "41",
	'd': "Zn", 'w': "H20", 'J': "tj",
	'e': "a1", 'v': "Cu", 'I': "i1", 'P': "big", 'W': "9u",
	'f': "34", 'u': "1zt", 'H': "d1",
	'g': "5y", 't': "3r", 'G': "AN", 'Q': "NE", 'V': "6c",
	'h': "18", 's': "3", 'F': "b",
	'i': "kr", 'r': "ml", 'E': "v7", 'R': "0O", 'U': "550",
	'j': "KR", 'q': "09", 'D': "9i", 'S': "L",
	'k': "CJ", 'p': "uo", 'C': "N0",
	'l': "45", 'o': "t8", 'B': "cA", 'T': "0IO",
	'm': "qo", 'n': "5", 'A': "x",
}

func HashSecret(psw string) (ans string) {
	//lens := len(psw)
	numCount, charCount := 0, 0
	tans := []rune{}
	for _, v := range psw {
		if v >= '0' && v <= '9' {
			ans += NUMHASH[v]
			numCount++
		} else if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			ans += CHARHASH[v]
			charCount++
		}
	}
	tans = []rune(ans)
	if numCount == charCount {
		numCount += charCount
		numCount %= len(tans)
	}

	if numCount > charCount {
		numCount ^= charCount
		charCount ^= numCount
		numCount ^= charCount
	}

	for i := 0; i < numCount; i++ {
		strReverse(&tans, numCount, charCount, len(tans))
		fmt.Println(string(tans))
	}
	return string(tans)
}

func strReverse(str *[]rune, start, end int, lens int) {
	p1, p2 := 0, len(*str)-1
	for p1 < start {
		(*str)[p1] ^= (*str)[start]
		(*str)[start] ^= (*str)[p1]
		(*str)[p1] ^= (*str)[start]
		p1++
		start--
	}
	for end < p2 {
		(*str)[end] ^= (*str)[p2]
		(*str)[p2] ^= (*str)[end]
		(*str)[end] ^= (*str)[p2]
		end++
		p2--
	}
	p1, p2 = 0, len(*str)-1
	for p1 < p2 {
		(*str)[p1] ^= (*str)[p2]
		(*str)[p2] ^= (*str)[p1]
		(*str)[p1] ^= (*str)[p2]
		p1++
		p2--
	}
}

// 返回pwd的加密码
func OnLock(pwd string, uName string) string {
	h := md4.New()
	io.WriteString(h, pwd)

	io.WriteString(h, salt1)
	io.WriteString(h, salt2)
	io.WriteString(h, uName)
	io.WriteString(h, salt3)

	fin := fmt.Sprintf("%X", h.Sum(nil))
	return fin
}
