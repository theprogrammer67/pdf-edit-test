package main

import (
	"fmt"
	"log"

	_ "github.com/pdfcrowd/pdfcrowd-go"
	"seehuhn.de/go/pdf"
)

func main() {
	r, err := pdf.Open("/home/stoi/temp/braun_g1500.pdf", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	cat := r.GetMeta().Catalog
	p := cat.Pages
	fmt.Println(p)

	// // create the API client instance
	// client := pdfcrowd.NewPdfToHtmlClient("demo", "ce544b6ea52a5621fb9d55f8b542d14d")

	// // run the conversion and write the result to a file
	// err := client.ConvertFileToFile("/home/stoi/temp/braun_g1500.pdf", "/home/stoi/temp/braun_g1500.html")

	// // check for the conversion error
	// handleError(err)
}

// func handleError(err error) {
// 	if err != nil {
// 		why, ok := err.(pdfcrowd.Error)
// 		if ok {
// 			os.Stderr.WriteString(fmt.Sprintf("Pdfcrowd Error: %s\n", why))
// 		} else {
// 			os.Stderr.WriteString(fmt.Sprintf("Generic Error: %s\n", err))
// 		}

// 		panic(err.Error())
// 	}
// }
