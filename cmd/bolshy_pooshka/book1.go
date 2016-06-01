package main

import "math/rand"

//go:generate go-bindata -pkg $GOPACKAGE -o book1txt.go book1.txt

func GetRandText(maxLen int) string {
	txtLen := rand.Intn(maxLen)
	txtPos := rand.Intn(len(_book1Txt) - txtLen)
	return string(_book1Txt[txtPos : txtPos+txtLen])
}
