package pdfcpu

type Stamp struct {
	Paper      string `json:"paper"`
	Crop       string `json:"crop"`
	Origin     string `json:"origin"`
	ContentBox bool   `json:"contentBox"`
	Debug      bool   `json:"debug"`
	Guides     bool   `json:"guides"`
	Dirs       struct {
		Images string `json:"images"`
	} `json:"dirs"`
	Files struct {
		Logo1 string `json:"logo1"`
	} `json:"files"`
	Borders struct {
		MyBorder struct {
			Width int    `json:"width"`
			Col   string `json:"col"`
			Style string `json:"style"`
		} `json:"myBorder"`
	} `json:"borders"`
	Fonts struct {
		LabelFont struct {
			Name string `json:"name"`
			Col  string `json:"col"`
			Size int    `json:"size"`
		} `json:"labelFont"`
	} `json:"fonts"`
	Texts struct {
		TextHeaderLabel struct {
			Align string `json:"align"`
			Font  struct {
				Name string `json:"name"`
			} `json:"font"`
			Value string `json:"value"`
		} `json:"textHeaderLabel"`
		TextHeaderValue struct {
			Align string `json:"align"`
			Font  struct {
				Name string `json:"name"`
			} `json:"font"`
			Value string `json:"value"`
		} `json:"textHeaderValue"`
		TextClientLabel struct {
			Align string `json:"align"`
			Font  struct {
				Name string `json:"name"`
			} `json:"font"`
			Value string `json:"value"`
		} `json:"textClientLabel"`
		TextClientValue struct {
			Align string `json:"align"`
			Font  struct {
				Name string `json:"name"`
			} `json:"font"`
			Value string `json:"value"`
		} `json:"textClientValue"`
		TextDocumentLabel struct {
			Align string `json:"align"`
			Font  struct {
				Name string `json:"name"`
			} `json:"font"`
			Value string `json:"value"`
		} `json:"textDocumentLabel"`
		TextDocumentValue struct {
			Align string `json:"align"`
			Font  struct {
				Name string `json:"name"`
			} `json:"font"`
			Value string `json:"value"`
		} `json:"textDocumentValue"`
	} `json:"texts"`
	Pages struct {
		Num1 struct {
			Content struct {
				Box []struct {
					Pos     []int  `json:"pos"`
					Width   int    `json:"width"`
					Height  int    `json:"height"`
					FillCol string `json:"fillCol"`
					Border  struct {
						Name string `json:"name"`
					} `json:"border"`
					Rot int `json:"rot"`
				} `json:"box"`
				Text []struct {
					Name string `json:"name"`
					Hide bool   `json:"hide"`
					Pos  []int  `json:"pos"`
				} `json:"text"`
				Image []struct {
					Src    string `json:"src"`
					Pos    []int  `json:"pos"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"image"`
			} `json:"content"`
		} `json:"1"`
	} `json:"pages"`
}
