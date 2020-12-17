// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package fuzzystrmatch

import (
	"strings"
	"unicode"
)

func getLetters(runes []rune, i, l int) string {
	if i >= len(runes) {
		return ""
	}
	if i < 0 {
		i = 0
	}
	if i+l > len(runes) {
		l = len(runes) - i
	}
	return string(runes[i : i+l])
}

func isVowel(letter string) bool {
	switch letter {
	case "A", "E", "I", "O", "U":
		return true
	default:
		return false
	}
}

// Metaphone calculate the metaphone code of a source string.
// https://en.wikipedia.org/wiki/Metaphone
func Metaphone(source string, maxLength int) (string, error) {
	if len(source) == 0 {
		return "", nil
	}

	// Skip leading non-alphabetic characters
	source = strings.TrimLeftFunc(source, func(r rune) bool {
		if r <= unicode.MaxASCII {
			return !(unicode.IsUpper(r) || unicode.IsLower(r))
		}
		return false
	})
	if len(source) == 0 {
		return "", nil
	}

	phonedWord := make([]byte, 0)

	rIdx := 0
	source = strings.ToUpper(source)
	runes := []rune(source)

	letter0 := getLetters(runes, 0, 1)
	letter1 := getLetters(runes, 1, 1)
	switch letter0 {
	case "A":
		// AE becomes E
		if letter1 == "E" {
			phonedWord = append(phonedWord, 'E')
			rIdx += 2
		} else {
			phonedWord = append(phonedWord, 'A')
			rIdx++
		}
	// [GKP]N becomes N
	case "G", "K", "P":
		if letter1 == "N" {
			phonedWord = append(phonedWord, 'N')
			rIdx += 2
		}
	// WH becomes H, WR becomes R W if followed by a vowel
	case "W":
		if letter1 == "H" || letter1 == "R" {
			phonedWord = append(phonedWord, letter1[0])
			rIdx += 2
		} else if isVowel(letter1) {
			phonedWord = append(phonedWord, 'W')
			rIdx += 2
		}
	// X becomes S
	case "X":
		phonedWord = append(phonedWord, 'S')
		rIdx++
	// Vowels are kept
	case "E", "I", "O", "U":
		phonedWord = append(phonedWord, letter0[0])
		rIdx++
	}
	prevLetter := ""
	if rIdx > 0 {
		prevLetter = getLetters(runes, rIdx-1, 1)
	}
	for i := rIdx; i < len(runes); i++ {
		// Ignore non-alphas
		if !unicode.IsUpper(runes[i]) {
			prevLetter = ""
			continue
		}
		currLetter := getLetters(runes, i, 1)
		// Drop duplicates, except CC
		if currLetter == prevLetter && currLetter != "C" {
			prevLetter = currLetter
			continue
		}
		switch currLetter {
		// B -> B unless in MB
		case "B":
			if prevLetter != "M" {
				phonedWord = append(phonedWord, 'B')
			}
		// 'sh' if -CIA- or -CH, but not SCH, except SCHW. (SCHW is
		// handled in S) S if -CI-, -CE- or -CY- dropped if -SCI-,
		// SCE-, -SCY- (handed in S) else K
		case "C":
			l2 := getLetters(runes, i, 2)
			l3 := getLetters(runes, i, 3)
			pl3 := getLetters(runes, i-1, 3)
			if l3 == "CIA" || l2 == "CH" {
				phonedWord = append(phonedWord, 'X')
			} else if l2 == "CI" || l2 == "CE" || l2 == "CY" {
				phonedWord = append(phonedWord, 'S')
			} else if pl3 != "SCI" || l2 != "SCE" || l2 != "SCY" {
				phonedWord = append(phonedWord, 'K')
			}
		// J if in -DGE-, -DGI- or -DGY- else T
		case "D":
			l3 := getLetters(runes, i, 3)
			if l3 == "DGE" || l3 == "DGI" || l3 == "DGY" {
				phonedWord = append(phonedWord, 'J')
			} else {
				phonedWord = append(phonedWord, 'T')
			}
		case "F":
			phonedWord = append(phonedWord, 'F')
		case "G":
			l2 := getLetters(runes, i, 2)
			l3 := getLetters(runes, i, 3)
			l4 := getLetters(runes, i, 4)
			nnl1 := getLetters(runes, i+2, 1)
			if (l2 == "GH" && !isVowel(nnl1)) ||
				l2 == "GN" || l4 == "GNED" ||
				l3 == "GDE" || l3 == "GDI" || l3 == "GDY" {

			} else if prevLetter == "I" || prevLetter == "E" || prevLetter == "Y" {
				phonedWord = append(phonedWord, 'J')
			} else {
				phonedWord = append(phonedWord, 'K')
			}
		case "H":
			nl1 := getLetters(runes, i+1, 1)
			ppl2 := getLetters(runes, i-2, 2)
			if isVowel(nl1) && ppl2 != "CH" &&
				ppl2 != "SH" && ppl2 != "PH" &&
				ppl2 != "TH" && ppl2 != "GH" {
				phonedWord = append(phonedWord, 'H')
			}
		case "J":
			phonedWord = append(phonedWord, 'J')
		case "K":
			pl1 := getLetters(runes, i-1, 1)
			if pl1 != "C" {
				phonedWord = append(phonedWord, 'K')
			}
		case "L", "M", "N":
			phonedWord = append(phonedWord, currLetter[0])
		case "P":
			nl1 := getLetters(runes, i+1, 1)
			if nl1 == "H" {
				phonedWord = append(phonedWord, 'F')
			} else {
				phonedWord = append(phonedWord, 'P')
			}
		case "Q":
			phonedWord = append(phonedWord, 'K')
		case "R":
			phonedWord = append(phonedWord, 'R')
		case "S":
			nl1 := getLetters(runes, i+1, 1)
			l3 := getLetters(runes, i, 3)
			if nl1 == "H" || l3 == "SIO" || l3 == "SIA" {
				phonedWord = append(phonedWord, 'X')
			} else {
				phonedWord = append(phonedWord, 'S')
			}
		case "T":
			nl1 := getLetters(runes, i+1, 1)
			l3 := getLetters(runes, i, 3)
			if l3 == "TIO" || l3 == "TIA" {
				phonedWord = append(phonedWord, 'X')
			} else if nl1 == "H" {
				phonedWord = append(phonedWord, '0')
			} else if l3 == "TCH" {
				phonedWord = append(phonedWord, 'T')
			}
		case "V":
			phonedWord = append(phonedWord, 'F')
		case "W":
			nl1 := getLetters(runes, i+1, 1)
			if isVowel(nl1) {
				phonedWord = append(phonedWord, 'W')
			}
		case "X":
			phonedWord = append(phonedWord, 'K')
			phonedWord = append(phonedWord, 'S')
		case "Y":
			nl1 := getLetters(runes, i+1, 1)
			if isVowel(nl1) {
				phonedWord = append(phonedWord, 'Y')
			}
		case "Z":
			phonedWord = append(phonedWord, 'S')
		}
		prevLetter = currLetter
		if len(phonedWord) >= maxLength {
			break
		}
	}
	return string(phonedWord), nil
}
