package service

import (
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/draw"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

var FontName string = "Roboto-Regular"
var stampDim types.Dim = types.Dim{Width: 595, Height: 50}

func CreatePdfStamp(w io.Writer, text string, conf *model.Configuration) error {
	mediaBox := types.RectForDim(stampDim.Width, stampDim.Height)
	xRefTable, err := pdfcpu.CreateDemoXRef()
	if err != nil {
		return err
	}

	p := model.NewPage(mediaBox)
	var region *types.Rectangle
	writeText(text, xRefTable, p, region, types.AlignLeft, stampDim.Width, 3, 3, 3, 3)

	rootDict, err := xRefTable.Catalog()
	if err != nil {
		return err
	}
	if err = pdfcpu.AddPageTreeWithSamplePage(xRefTable, rootDict, p); err != nil {
		return err
	}

	err = createPdf(xRefTable, w, nil)
	if err != nil {
		return err
	}

	return nil
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

	k := p.Fm.EnsureKey(FontName)

	td := model.TextDescriptor{
		FontName: FontName,
		FontKey:  k,
		FontSize: 8,
		// ShowMargins:    true,
		MLeft:          mLeft,
		MRight:         mRight,
		MTop:           mTop,
		MBot:           mBot,
		Scale:          1.,
		ScaleAbs:       true,
		HAlign:         hAlign,
		RMode:          draw.RMFill,
		StrokeCol:      color.Blue,
		FillCol:        color.Blue,
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
}

func createPdf(xRefTable *model.XRefTable, w io.Writer, conf *model.Configuration) error {
	ctx := pdfcpu.CreateContext(xRefTable, conf)
	return api.WriteContext(ctx, w)
}
