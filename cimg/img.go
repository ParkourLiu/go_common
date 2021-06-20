package cimg

import (
	"fmt"
	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"os"
)

//图片色彩反转
func FZImage(m image.Image) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			r_uint8 := uint8(r >> 8)
			g_uint8 := uint8(g >> 8)
			b_uint8 := uint8(b >> 8)
			a_uint8 := uint8(a >> 8)
			r_uint8 = 255 - r_uint8
			g_uint8 = 255 - g_uint8
			b_uint8 = 255 - b_uint8
			newRgba.SetRGBA(i, j, color.RGBA{r_uint8, g_uint8, b_uint8, a_uint8})
		}
	}
	return newRgba
}

//图片灰化处理
func HDImage(m image.Image) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			_, g, _, a := colorRgb.RGBA()
			g_uint8 := uint8(g >> 8)
			a_uint8 := uint8(a >> 8)
			newRgba.SetRGBA(i, j, color.RGBA{g_uint8, g_uint8, g_uint8, a_uint8})
		}
	}
	return newRgba
}

//图片缩放, add at 2018-9-12
func RectImage(m image.Image, newdx int) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(image.Rect(0, 0, newdx, newdx*dy/dx))
	graphics.Scale(newRgba, m)
	return newRgba
}

//图片转为字符画（简易版）
func Ascllimage(m image.Image, filePath string) {
	if m.Bounds().Dx() > 300 {
		m = RectImage(m, 300)
	}
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	arr := []string{"M", "N", "H", "Q", "$", "O", "C", "?", "7", ">", "!", ":", "–", ";", "."}

	dstFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			colorRgb := m.At(j, i)
			_, g, _, _ := colorRgb.RGBA()
			avg := uint8(g >> 8)
			num := avg / 18
			dstFile.WriteString(arr[num])
			if j == dx-1 {
				dstFile.WriteString("\n")
			}
		}
	}

}

//图片合并,横向
func MergeX(imgs []*image.RGBA) *image.NRGBA {
	px := 0
	py := 0
	for _, v := range imgs {
		px = px + v.Bounds().Dx() //屏幕累加总长度
		dy := v.Bounds().Dy()     //屏幕最大高度
		if dy > py {
			py = dy
		}
	}
	newRgba := imaging.New(px, py, color.NRGBA{0, 0, 0, 0})
	x := 0
	for _, v := range imgs {
		newRgba = imaging.Paste(newRgba, v, image.Pt(x, 0))
		x = x + v.Bounds().Dx()
	}
	return newRgba
}

//图片合并,竖向
func MergeY(imgs []*image.RGBA) *image.NRGBA {
	px := 0
	py := 0
	for _, v := range imgs {
		py = py + v.Bounds().Dy() //屏幕累加总长度
		dx := v.Bounds().Dx()     //屏幕最大高度
		if dx > px {
			py = dx
		}
	}
	newRgba := imaging.New(px, py, color.NRGBA{0, 0, 0, 0})
	y := 0
	for _, v := range imgs {
		newRgba = imaging.Paste(newRgba, v, image.Pt(0, y))
		y = y + v.Bounds().Dy()
	}
	return newRgba
}
