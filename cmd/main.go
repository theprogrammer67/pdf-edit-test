package main

import (
	"fmt"
	"log"
	"path/filepath"
	"pdf-edit/internal/service"
	"pdf-edit/pkg/filebuffer"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/font"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/draw"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func main() {
	var err error

	fileDir, _ := filepath.Abs("../internal/samples")

	conf := api.LoadConfiguration()
	// font.UserFontDir, _ = filepath.Abs("../internal/fonts")
	font.UserFontDir = filepath.Join(fileDir, "fonts")
	fmt.Printf("Fonts dir: %s\n", font.UserFontDir)

	fp := filepath.Join(font.UserFontDir, service.FontName+".ttf")
	if err := api.InstallFonts([]string{fp}); err != nil {
		log.Printf("Error install font: %v\n", err)
	}

	// JSON to PDF
	fileName := "JsonPdf"
	inFile := filepath.Join(fileDir, fileName+".json")
	outFile := filepath.Join(fileDir, fileName+".pdf")
	err = service.CreatePdfFromJson(inFile, outFile, conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	// PDF stamp
	fileName = "Форма договора для юридических лиц"
	wmFileName := "JsonPdf"
	inFile = filepath.Join(fileDir, fileName+".pdf")
	outFile = filepath.Join(fileDir, fileName+"_stamp.pdf")
	wmFile := filepath.Join(fileDir, wmFileName+".pdf")

	err = service.AddPdfStamp(inFile, outFile, wmFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	return

	// Text stamp
	fileName = "Форма договора для юридических лиц"
	inFile = filepath.Join(fileDir, fileName+".pdf")
	outFile = filepath.Join(fileDir, fileName+"_stamp.pdf")

	inBuff, err := filebuffer.ReadFile(inFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	outBuff := filebuffer.NewFileBuffer(nil)

	stamp := "Документ подписан электронной подписью 30.10.2923 16:10 (МСК)\nКлиент    Курбатов Андрей Алексеевич\nЭлектронный документ    A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5A5"
	err = service.AddTextStamp(inBuff, outBuff, stamp, conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = outBuff.WriteFile(outFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// PDF stamp
	outBuff = filebuffer.NewFileBuffer(nil)
	outFile = filepath.Join(fileDir, "PDF_stamp.pdf")
	err = service.CreatePdfStamp(outBuff, stamp, conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = outBuff.WriteFile(outFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	return

	fp = filepath.Join(fileDir, "Certificate.pdf")
	service.CreatePdf(fp, "This is certificate\nSecond line\nЭто сертификат")

	align, rtl := "l", "off"
	desc := fmt.Sprintf("font:%s, rtl:%s, align:%s, scale:1.0 rel, rot:0, fillc:#000000, bgcol:#ab6f30, margin:10, border:10 round, opacity:.7", "Roboto-Regular", rtl, align)
	var pages = []string{"7"}
	err = api.AddTextWatermarksFile(inFile, outFile, pages, true, "Документ подписан электронной подписью\nКлиент\nЭлектронный документ", desc, nil)
	if err != nil {
		log.Fatalf("AddTextWatermarksFile %s: %v\n", outFile, err)
	}

	// createPdf("/home/stoi/temp/Text_demo.pdf")

	// Merge inFiles by concatenation in the order specified and write the result to out.pdf.
	// out.pdf will be overwritten.

	// inFiles := []string{"/home/stoi/temp/Loyalty system.pdf", "/home/stoi/temp/Стовпец Игорь Александрович.pdf"}
	// err := api.MergeCreateFile(inFiles, "/home/stoi/temp/out.pdf", nil)
	// if err != nil {
	// 	log.Println(err)
	// 	log.Fatal(err.Error())
	// }

	// f, err := os.Open("/home/stoi/temp/braun_g1500.pdf")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// defer f.Close()

	b, err := filebuffer.ReadFile(inFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	info, err := api.PDFInfo(b, "braun_g1500.pdf", nil, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	if info == nil {
		log.Fatal("Empty file info")
	}

	fmt.Printf("Page count: %d\n", info.PageCount)
	fmt.Printf("Page dimensions: %v\n", info.PageDimensions)
	fmt.Printf("Info: %v\n", info)

}

// func createPdf() {
// 	mediaBox := types.RectForFormat("A4")
// 	p := model.Page{MediaBox: mediaBox, Fm: model.FontMap{}, Buf: new(bytes.Buffer)}
// 	pdfcpu.CreateTestPageContent(p)

// 	xRefTable, err := pdfcpu.CreateDemoXRef()
// 	if err != nil {
// 		log.Fatalf("createPdf: %v\n", err)
// 	}
// 	rootDict, err := xRefTable.Catalog()
// 	if err != nil {
// 		log.Fatalf("createPdf: %v\n", err)
// 	}
// 	if err = pdfcpu.AddPageTreeWithSamplePage(xRefTable, rootDict, p); err != nil {
// 		log.Fatalf("createPdf: %v\n", err)
// 	}

// 	if err := api.CreatePDFFile(xRefTable, "/home/stoi/temp/Test.pdf", nil); err != nil {
// 		log.Fatalf("createPdf: %v\n", err)
// 	}
// }

func createPdf(fileName string) {
	mediaBox := types.RectForDim(float64(600), float64(600))
	createXRefAndWritePDF(fileName, mediaBox, createTextDemoAlignLeft)

}

func createXRefAndWritePDF(fileName string, mediaBox *types.Rectangle, f func(xRefTable *model.XRefTable, mediaBox *types.Rectangle) model.Page) {
	xRefTable, err := pdfcpu.CreateDemoXRef()
	if err != nil {
		log.Fatalf("createXRefAndWritePDF: %v\n", err)
	}

	p := f(xRefTable, mediaBox)

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

func createTextDemoAlignLeft(xRefTable *model.XRefTable, mediaBox *types.Rectangle) model.Page {
	return createTextDemoAlignedWidthAndMargin(xRefTable, mediaBox, types.AlignLeft, 0, 0, 0, 0, 0)
}

func createTextDemoAlignedWidthAndMargin(xRefTable *model.XRefTable, mediaBox *types.Rectangle, hAlign types.HAlignment, w, mLeft, mRight, mTop, mBot float64) model.Page {
	p := model.NewPage(mediaBox)
	var region *types.Rectangle
	writeTextDemoAlignedWidthAndMargin(xRefTable, p, region, hAlign, w, mLeft, mRight, mTop, mBot)
	region = types.RectForWidthAndHeight(50, 70, 200, 200)
	writeTextDemoAlignedWidthAndMargin(xRefTable, p, region, hAlign, w, mLeft, mRight, mTop, mBot)
	return p
}

func writeTextDemoAlignedWidthAndMargin(
	xRefTable *model.XRefTable,
	p model.Page,
	region *types.Rectangle,
	hAlign types.HAlignment,
	w, mLeft, mRight, mTop, mBot float64) {

	buf := p.Buf
	mediaBox := p.MediaBox

	mediaBB := true

	var cr, cg, cb float32
	cr, cg, cb = .5, .75, 1.
	r := mediaBox
	if region != nil {
		r = region
		cr, cg, cb = .75, .75, 1
	}
	if mediaBB {
		draw.FillRectNoBorder(buf, r, color.SimpleColor{R: cr, G: cg, B: cb})
	}

	fontName := "Helvetica"
	k := p.Fm.EnsureKey(fontName)

	td := model.TextDescriptor{
		FontName:       fontName,
		FontKey:        k,
		FontSize:       24,
		ShowMargins:    true,
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
		ShowLineBB:     true,
		ShowTextBB:     true,
		HairCross:      true,
	}

	td.VAlign, td.X, td.Y, td.Text = types.AlignBaseline, -1, r.Height()*.75, "M\\u(lti\nline\n\nwith empty line"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	td.VAlign, td.X, td.Y, td.Text = types.AlignBaseline, r.Width()*.75, r.Height()*.25, "Arbitrary\ntext\nlines"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	// Multilines along the top of the page:
	td.VAlign, td.X, td.Y, td.Text = types.AlignTop, 0, r.Height(), "0,h (topleft)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	td.VAlign, td.X, td.Y, td.Text = types.AlignTop, -1, r.Height(), "-1,h (topcenter)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	td.VAlign, td.X, td.Y, td.Text = types.AlignTop, r.Width(), r.Height(), "w,h (topright)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	// Multilines along the center of the page:
	// x = 0 centers the position of multilines horizontally
	// y = 0 centers the position of multilines vertically and enforces alignMiddle
	td.VAlign, td.X, td.Y, td.Text = types.AlignBaseline, 0, -1, "0,-1 (left)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	td.VAlign, td.X, td.Y, td.Text = types.AlignMiddle, -1, -1, "-1,-1 (center)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	td.VAlign, td.X, td.Y, td.Text = types.AlignBaseline, r.Width(), -1, "w,-1 (right)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	// Multilines along the bottom of the page:
	td.VAlign, td.X, td.Y, td.Text = types.AlignBottom, 0, 0, "0,0 (botleft)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	td.VAlign, td.X, td.Y, td.Text = types.AlignBottom, -1, 0, "-1,0 (botcenter)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	td.VAlign, td.X, td.Y, td.Text = types.AlignBottom, r.Width(), 0, "w,0 (botright)\nand line2"
	model.WriteColumn(xRefTable, buf, mediaBox, region, td, w)

	draw.DrawHairCross(buf, 0, 0, r)
}
