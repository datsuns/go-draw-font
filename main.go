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
	"path/filepath"
	"sync"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	//"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"gopkg.in/yaml.v2"
)

var (
	DestRoot    = "output"
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
		fixed.I(815 + 2),
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

func gen_png(ft *truetype.Font, opt *truetype.Options, cfg *Config, title string, list []string) {
	fmt.Printf("generate [%v] start\n", title)
	imageWidth := cfg.Image.Width
	imageHeight := cfg.Image.Height

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	face := truetype.NewFace(ft, opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	file, err := os.Create(filepath.Join(DestRoot, title+".png"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	h_idx := 0
	for i, c := range list {
		if c == "." {
			continue
		}
		dr.Dot.X = fixed.I(cfg.XPos[i%7])
		dr.Dot.Y = fixed.I(cfg.YPos[h_idx])
		buf.Reset()
		//fmt.Printf("%2v) x:%v, y:%v char[%v]\n", i, dr.Dot.X, dr.Dot.Y, c)
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

func day_exists(year, month, day int) bool {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if date.Year() == year && date.Month() == time.Month(month) && date.Day() == day {
		return true
	} else {
		return false
	}
}

func gen_month_text(year, month int) (string, []string) {
	body := []string{}
	title := fmt.Sprintf("%v-%02v", year, month)
	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	for i := time.Sunday; i < t.Weekday(); i++ {
		body = append(body, ".")
	}
	for i := 1; ; i++ {
		body = append(body, fmt.Sprintf("%v", i))
		if !day_exists(year, month, i+1) {
			break
		}
	}
	return title, body
}

func gen_day_list(cfg *Config) map[string][]string {
	if cfg.Output.Year == 0 {
		panic("plese set Year in config")
	}
	mlist := []int{}
	if cfg.Output.Month == 0 {
		mlist = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	} else {
		mlist = []int{cfg.Output.Month}
	}
	ret := map[string][]string{}
	for _, m := range mlist {
		title, body := gen_month_text(cfg.Output.Year, m)
		ret[title] = body
	}
	return ret
}

func main() {
	cfg, err := load_config("config.yaml")
	if err != nil {
		panic(err)
	}
	cfg.Dump()
	os.MkdirAll(DestRoot, 0777)
	ft := load_font(cfg.Font)
	opt := truetype.Options{
		Size:              cfg.Size,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	day_list := gen_day_list(cfg)
	wg := &sync.WaitGroup{}
	for n, t := range day_list {
		wg.Add(1)
		go func(title string, body []string) {
			gen_png(ft, &opt, cfg, title, body)
			wg.Done()
		}(n, t)
	}
	wg.Wait()
}
