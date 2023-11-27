package service

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/font"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
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

// Stamp

func AddStamp(rs io.ReadSeeker, w io.Writer, text string, conf *model.Configuration) error {
	info, err := api.PDFInfo(rs, "", nil, nil)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New("empty file info")
	}
	var pages = []string{strconv.Itoa(info.PageCount)}

	wm := NewTextWatermark(text)

	return api.AddWatermarks(rs, w, pages, wm, conf)
}

func NewTextWatermark(text string) *model.Watermark {
	textAlign := types.AlignLeft
	bgColor := color.White
	borderColor := color.Blue
	textColor := color.Blue

	wm := model.DefaultWatermarkConfig()
	wm.OnTop = true // stamp
	wm.InpUnit = types.POINTS
	wm.Update = false
	wm.Mode = model.WMText
	wm.FontName = FontName
	wm.FontSize = 1
	wm.ScaledFontSize = 1
	wm.Pos = types.BottomCenter
	wm.Rotation = 0
	wm.UserRotOrDiagonal = true
	wm.Diagonal = model.NoDiagonal
	wm.HAlign = &textAlign
	wm.Color = textColor
	wm.StrokeColor = textColor
	wm.FillColor = textColor
	wm.FillColor = textColor
	wm.MLeft, wm.MRight = 3, 3
	wm.MTop, wm.MBot = 3, 3
	wm.BorderWidth = 2
	wm.BorderStyle = types.LJRound
	wm.BorderColor = &borderColor
	wm.BgColor = &bgColor
	wm.Scale = 1
	// wm.ScaleAbs = true
	wm.Opacity = 0.7
	// wm.Dx = 10
	wm.Dy = 10
	setTextWatermark(text, wm)

	return wm
}

func setTextWatermark(s string, wm *model.Watermark) {
	wm.TextString = s
	if font.IsCoreFont(wm.FontName) {
		bb := []byte{}
		for _, r := range s {
			// Unicode => char code
			b := byte(0x20) // better use glyph: .notdef
			if r <= 0xff {
				b = byte(r)
			}
			bb = append(bb, b)
		}
		s = string(bb)
	} else {
		bb := []byte{}
		u := utf16.Encode([]rune(s))
		for _, i := range u {
			bb = append(bb, byte((i>>8)&0xFF))
			bb = append(bb, byte(i&0xFF))
		}
		s = string(bb)
	}
	s = strings.ReplaceAll(s, "\\n", "\n")
	wm.TextLines = append(wm.TextLines, strings.FieldsFunc(s, func(c rune) bool { return c == 0x0a })...)
}
