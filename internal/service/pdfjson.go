package service

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func CreatePdfFromJson(inFileJSON, outFile string, conf *model.Configuration) error {
	err := api.CreateFile("", inFileJSON, outFile, conf)

	return err
}
