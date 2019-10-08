// ref1) https://qiita.com/n-noguchi/items/566e83c5cc0d3b80852c
//		ざっとしたコードのベース
// ref2) https://qiita.com/uobikiemukot/items/11dac0f1418492493226
//		基礎を書いてくれてそう

// 左よせ:
//   1文字入れて固定幅でオフセットしていくイメージ
// センタリング:
//   さて頑張る

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	//"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
)

func load_font() *truetype.Font {
	// フォントファイルを読み込み
	ftBinary, err := ioutil.ReadFile("GN-KillGothic-U-KanaNA.ttf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ft, err := truetype.Parse(ftBinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return ft
}

func main() {
	ft := load_font()

	opt := truetype.Options{
		//Size:              60,
		Size:              30,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	imageWidth := 1500
	imageHeight := 1000

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	text := []string{
		".", ".", "1", "2", "3", "4", "5",
		"6", "7", "8", "9", "10", "11", "12",
		"13", "14", "15", "16", "17", "18", "19",
		"20", "21", "22", "23", "24", "25", "26",
		"27", "28", "29", "30", "31",
	}
	xpos := []fixed.Int26_6{
		fixed.I(0),
		fixed.I(150 - 7),
		fixed.I(300 - 14),
		fixed.I(450 - 21),
		fixed.I(600 - 28),
		fixed.I(750 - 35),
		fixed.I(900 - 42),
	}
	hpos := []fixed.Int26_6{
		fixed.I(90),
		fixed.I(235),
		fixed.I(380 + 2),
		fixed.I(525 + 2),
		fixed.I(670 + 2),
	}

	file, err := os.Create(`test.png`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	h_idx := 0
	for i, c := range text {
		if c == "." {
			continue
		}
		dr.Dot.X = xpos[i%7]
		dr.Dot.Y = hpos[h_idx]
		buf.Reset()
		fmt.Printf("%2v) x:%v, y:%v char[%v]\n", i, dr.Dot.X, dr.Dot.Y, c)
		dr.DrawString(c)
		err = png.Encode(buf, img)
		if (i > 0) && (i%7 == 6) {
			h_idx += 1
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	file.Write(buf.Bytes())
}