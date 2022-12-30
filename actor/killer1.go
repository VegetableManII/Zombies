package actor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	_ "image/png"
	"log"
)

const (
	killerFrameWidth  = 70
	killerFrameHeight = 70
)

type Killer1 struct {
	KillerBase
}

var killerImage *ebiten.Image

func init() {
	img, _, err := ebitenutil.NewImageFromFile("./Resources/killer1.png")
	if err != nil {
		log.Fatalf("Actor.%s", err)
	}
	killerImage = img
}

// GetSubImage 获killer的图像
func (k *Killer1) getSubImage() image.Image {
	var img image.Image
	if !k.attackModle {
		k.countKiller = 0
		img = killerImage.SubImage(image.Rect(0, k.Direction*killerFrameHeight, 0+killerFrameWidth, k.Direction*killerFrameHeight+killerFrameHeight))
	} else {
		if k.countKiller == k.RefreshRates*4 {
			k.countKiller = 0
			k.attackModle = false
		}
		pixCount := k.countKiller / k.RefreshRates
		img = killerImage.SubImage(image.Rect(pixCount*killerFrameWidth, k.Direction*killerFrameHeight, pixCount*killerFrameWidth+killerFrameWidth, k.Direction*killerFrameHeight+killerFrameHeight))
		k.countKiller++
	}
	return img
}

// SetMove x & y 是killer当前的位置
func (k *Killer1) SetMove(x, y float64) {
	k.movX, k.movY = x, y
	if x == -1 {
		k.Direction = 1
	}
	if x == 1 {
		k.Direction = 2
	}
	if y == -1 {
		k.Direction = 3
	}
	if y == 1 {
		k.Direction = 0
	}

}

func (k *Killer1) SelfUpdate(screen *ebiten.Image) {
	screen.DrawImage(k.getSubImage().(*ebiten.Image), k.getUpdateDrawImageOptions())
}

func (k *Killer1) HitArea(x, y float64) bool {
	rigth := k.PosX + 20*k.Scale // 35
	left := k.PosX - 20*k.Scale
	up := k.PosY - 10*k.Scale // 35
	down := k.PosY + 10*k.Scale
	return (x < rigth && x > left && y > up && y < down) && k.AttackModle()
}
