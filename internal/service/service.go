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

func AddWatermark(rs io.ReadSeeker, w io.Writer, text string, conf *model.Configuration) error {
	info, err := api.PDFInfo(rs, "", nil, nil)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New("empty file info")
	}
	var pages = []string{strconv.Itoa(info.PageCount)}

	// font := "Roboto-Regular"
	// desc := fmt.Sprintf("font:%s, rtl:off, align:l, scale:1.0 rel, rot:0, fillc:#000000, bgcol:#ab6f30, margin:10, border:10 round, opacity:.7", font)
	// unit := types.POINTS
	// wm, err := api.TextWatermark(text, desc, true, false, unit)
	// if err != nil {
	// 	return err
	// }
	// wm.Pos = types.BottomCenter

	wm := TextWatermark(text)

	return api.AddWatermarks(rs, w, pages, wm, conf)
}

func TextWatermark(text string) *model.Watermark {
	textAlign := types.AlignLeft
	bgColor := color.White
	borderColor := color.Blue
	textColor := color.Blue

	wm := model.DefaultWatermarkConfig()
	wm.OnTop = true
	wm.InpUnit = types.POINTS
	wm.Update = false
	wm.Mode = model.WMText
	wm.FontName = "Roboto-Regular"
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
