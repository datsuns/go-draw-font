package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

const (
	imagePath     = "./noto-emoji/png/128"
	fontSize      = 64  // point
	imageWidth    = 640 // pixel
	imageHeight   = 120 // pixel
	textTopMargin = 80  // fixed.I
)

func main() {
	// TrueType フォントの読み込み
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitFailure)
	}

	opt := truetype.Options{
		Size:              fontSize,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	text := "Hello, world! 👋"

	// 描画の初期位置
	dr.Dot.X = (fixed.I(imageWidth) - dr.MeasureString(text)) / 2
	dr.Dot.Y = fixed.I(textTopMargin)

	// 一文字ずつ描画していく
	for _, r := range text {
		path := fmt.Sprintf("%s/emoji_u%.4x.png", imagePath, r)
		_, err = os.Stat(path)

		// 画像ファイルが存在する場合に err == nil
		if err == nil {
			// 画像の読み込み
			fp, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			defer fp.Close()

			emoji, _, err := image.Decode(fp)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			// フォントのサイズと絵文字のサイズを合わせるためにリサイズ
			size := dr.Face.Metrics().Ascent.Floor() + dr.Face.Metrics().Descent.Floor()
			rect := image.Rect(0, 0, size, size)
			dst := image.NewRGBA(rect)
			draw.ApproxBiLinear.Scale(dst, rect, emoji, emoji.Bounds(), draw.Over, nil)

			// font.Drawer.Dot はグリフの baseline の座標を指している (だいたいグリフの左下の地点)
			// 一方 draw.Draw で画像を描画する場合には左上の座標が必要になる
			p := image.Pt(dr.Dot.X.Floor(), dr.Dot.Y.Floor()-dr.Face.Metrics().Ascent.Floor())
			draw.Draw(img, rect.Add(p), dst, image.ZP, draw.Over)
			dr.Dot.X += fixed.I(size)
		} else {
			// 対応するカラー絵文字が存在しないので TrueType フォントを使用する
			dr.DrawString(string(r))
		}
	}

	// JPEG に変換して stdout に出力
	buf := &bytes.Buffer{}
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitFailure)
	}

	_, err = io.Copy(os.Stdout, buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitFailure)
	}
}
