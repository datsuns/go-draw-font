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
	"gopkg.in/yaml.v2"
)

var (
	DefaultXPos = []fixed.Int26_6{
		fixed.I(0),
		fixed.I(150 - 7),
		fixed.I(300 - 14),
		fixed.I(450 - 21),
		fixed.I(600 - 28),
		fixed.I(750 - 35),
		fixed.I(900 - 42),
	}
	DefaultYPos = []fixed.Int26_6{
		fixed.I(90),
		fixed.I(235),
		fixed.I(380 + 2),
		fixed.I(525 + 2),
		fixed.I(670 + 2),
	}
)

type Config struct {
	Output struct {
		Year  int
		Month int
	}
	Font  string
	Size  float64
	Image struct {
		Width  int `yaml:"width"`
		Height int `yaml:"height"`
	}
	XPos []int `yaml:"XPos,flow"`
	YPos []int `yaml:"YPos,flow"`
}

func (c *Config) Dump() {
	fmt.Printf("output : %v-%v\n", c.Output.Year, c.Output.Month)
	fmt.Printf("font file : %v\n", c.Font)
	fmt.Printf("font size : %v\n", c.Size)
	fmt.Printf("    image : %vx%v\n", c.Image.Width, c.Image.Height)
	fmt.Printf("    Xpos  : ")
	for i, p := range c.XPos {
		fmt.Printf("[%v:%v],", i, p)
	}
	fmt.Printf("\n")
	fmt.Printf("    Ypos  : ")
	for i, p := range c.YPos {
		fmt.Printf("[%v:%v],", i, p)
	}
	fmt.Printf("\n")
}

func load_config(path string) (*Config, error) {
	ret := Config{}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func load_font(path string) *truetype.Font {
	// フォントファイルを読み込み
	ftBinary, err := ioutil.ReadFile(path)
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
	cfg, err := load_config("config.yaml")
	if err != nil {
		panic(err)
	}
	cfg.Dump()
	ft := load_font(cfg.Font)
	opt := truetype.Options{
		Size:              cfg.Size,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	imageWidth := cfg.Image.Width
	imageHeight := cfg.Image.Height

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
		dr.Dot.X = fixed.I(cfg.XPos[i%7])
		dr.Dot.Y = fixed.I(cfg.YPos[h_idx])
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
