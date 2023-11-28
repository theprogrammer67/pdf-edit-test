package pdfcpu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type PdfCpu struct {
	resDir    string
	conf      *model.Configuration
	fontName  string
	stampDim  types.Dim
	stampFile string
	template  *Stamp
}

type StampParams struct {
	Header   string
	Client   string
	Document string
}

func New(resDir string) *PdfCpu {
	var err error
	var b []byte

	p := &PdfCpu{
		resDir:    resDir,
		fontName:  "Roboto-Regular",
		stampDim:  types.Dim{Width: 595, Height: 50},
		stampFile: "stamp.pdf",
		template:  &Stamp{},
	}
	p.conf = api.LoadConfiguration()

	// fonts
	font.UserFontDir = filepath.Join(p.resDir, "fonts")
	err = api.InstallFonts([]string{filepath.Join(font.UserFontDir, p.fontName+".ttf")})
	check(err)

	// template
	b, err = os.ReadFile(filepath.Join(p.resDir, "template.json"))
	check(err)

	err = json.Unmarshal(b, p.template)
	check(err)
	p.template.Dirs.Images = filepath.Join(p.resDir, "images")

	return p
}

func (p *PdfCpu) AddPdfStamp(inFile, outFile string, params *StampParams) error {
	var err error
	var b []byte
	var info *pdfcpu.PDFInfo

	info, err = p.getPdfInfo(inFile)
	if err == nil {
		err = checkPaperSize(info)
		if err == nil {
			pages := []string{strconv.Itoa(info.PageCount)}

			stampJsonFile := outFile + ".json"
			stampPdfFile := outFile + ".stamp.pdf"
			stamp := *p.template
			stamp.Texts.TextHeaderValue.Value = params.Header
			stamp.Texts.TextClientValue.Value = params.Client
			stamp.Texts.TextDocumentValue.Value = params.Document

			b, err = json.Marshal(&stamp)
			if err == nil {
				err = os.WriteFile(stampJsonFile, b, 0644)

				if err == nil {
					defer func() {
						os.Remove(stampJsonFile)
					}()

					err = api.CreateFile("", stampJsonFile, stampPdfFile, p.conf)
					if err == nil {
						defer func() {
							os.Remove(stampPdfFile)
						}()

						var wm *model.Watermark
						wm, err = api.PDFWatermark(stampPdfFile, "sc:1.0 abs, rotation:0", true, false, types.POINTS)
						if err == nil {
							err = api.AddWatermarksFile(inFile, outFile, pages, wm, p.conf)
						}

					}
				}
			}
		}
	}

	return err
}

func (p *PdfCpu) getPdfInfo(inFile string) (*pdfcpu.PDFInfo, error) {
	var info *pdfcpu.PDFInfo
	var err error
	var f *os.File

	f, err = os.Open(inFile)
	if err == nil {
		defer f.Close()

		info, err = api.PDFInfo(f, inFile, nil, p.conf)
		if (err == nil) && (info == nil) {
			err = errors.New("missing PDF Info")
		}
	}

	return info, err
}

func checkPaperSize(info *pdfcpu.PDFInfo) error {
	// A4 210 x 297 mm
	for d := range info.PageDimensions {
		dc := d.ToMillimetres()
		if (dc.Width < 210) || (dc.Height < 297) {
			e := fmt.Sprintf("page format %.2f x %.2f mm is unsupported", dc.Width, dc.Height)
			return errors.New(e)
		}
	}

	return nil
}

/////////////

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

func AddPdfStamp(inFile, outFile, wmFile string) error {
	var err error
	var wm *model.Watermark
	onTop := false
	update := false

	wm, err = api.PDFWatermark(wmFile, "sc:1.0 abs, rotation:0", onTop, update, types.POINTS)
	if err == nil {
		err = api.AddWatermarksFile(inFile, outFile, nil, wm, nil)
	}

	return err
}

func AddTextStamp(rs io.ReadSeeker, w io.Writer, text string, conf *model.Configuration) error {
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
