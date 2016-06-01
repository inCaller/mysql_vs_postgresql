package main

import "math/rand"

//go:generate go-bindata -pkg $GOPACKAGE -o book1txt.go book1.txt

var txt = func() []byte {
	txt, err := book1TxtBytes()
	if err != nil {
		panic(err)
	}
	return txt
}()

func GetRandText(maxLen int) string {
	txtLen := rand.Intn(maxLen)
	txtPos := rand.Intn(len(txt) - txtLen)
	return string(txt[txtPos : txtPos+txtLen])
}
