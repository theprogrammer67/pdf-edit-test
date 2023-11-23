package main

import (
	"bytes"
	"fmt"

	"github.com/dslipak/pdf"
)

func main() {
	// pdf.DebugOn = true
	content, err := readPdf("/home/stoi/temp/Loyalty system.pdf") // Read local pdf file
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}

func readPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	// remember close file

	totalPage := r.NumPage()
	fmt.Println(totalPage)

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
