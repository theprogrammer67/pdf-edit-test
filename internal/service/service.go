package service

import (
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/draw"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func CreatePdf(fileName string, text string) {
	mediaBox := types.RectForFormat("A4")

	xRefTable, err := pdfcpu.CreateDemoXRef()
	if err != nil {
		log.Fatalf("createXRefAndWritePDF: %v\n", err)
	}

	p := model.NewPage(mediaBox)
	var region *types.Rectangle
	writeText(text, xRefTable, p, region, types.AlignLeft, 0, 0, 0, 0, 0)

	rootDict, err := xRefTable.Catalog()
	if err != nil {
		log.Fatalf("createXRefAndWritePDF: %v\n", err)
	}
	if err = pdfcpu.AddPageTreeWithSamplePage(xRefTable, rootDict, p); err != nil {
		log.Fatalf("createXRefAndWritePDF: %v\n", err)
	}

	if err := api.CreatePDFFile(xRefTable, fileName, nil); err != nil {
		log.Fatalf("createXRefAndWritePDF: %v\n", err)
	}
}

func writeText(
	text string,
	xRefTable *model.XRefTable,
	p model.Page,
	region *types.Rectangle,
	hAlign types.HAlignment,
	w, mLeft, mRight, mTop, mBot float64) {

	buf := p.Buf
	mediaBox := p.MediaBox

	r := mediaBox
	if region != nil {
		r = region
	}

	// Courier
	// Courier-Bold
	// Courier-BoldOblique
	// Courier-Oblique
	// Helvetica
	// Helvetica-Bold
	// Helvetica-BoldOblique
	// Helvetica-Oblique
	// Times-Roman
	// Times-Bold
	// Times-Italic
	// Times-BoldItalic
	// Symbol
	// ZapfDingbats

	// fontName := "Times-BoldItalic"
	// fontName := "Times-Roman"
	fontName := "Roboto-Regular"
	k := p.Fm.EnsureKey(fontName)

	td := model.TextDescriptor{
		FontName: fontName,
		FontKey:  k,
		FontSize: 24,
		// ShowMargins:    true,
		MLeft:          mLeft,
		MRight:         mRight,
		MTop:           mTop,
		MBot:           mBot,
		Scale:          1.,
		ScaleAbs:       true,
		HAlign:         hAlign,
		RMode:          draw.RMFill,
		StrokeCol:      color.Black,
		FillCol:        color.Black,
		ShowBackground: true,
		BackgroundCol:  color.SimpleColor{R: 1., G: .98, B: .77},
		ShowBorder:     true,
		// ShowLineBB: true,
		ShowTextBB: true,
		// HairCross:      true,
	}

	// Multilines along the top of the page:
	td.VAlign, td.X, td.Y, td.Text = types.AlignTop, 0, r.Height(), text
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	// draw.DrawHairCross(buf, 0, 0, r)
}
