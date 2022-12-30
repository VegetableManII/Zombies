package actor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

const (
	killer2FrameWidth  = 80
	killer2FrameHeight = 80
)

type Killer2 struct {
	PosX, PosY   float64
	Speed        float64
	Scale        float64
	RefreshRates int
	Direction    int
	movX, movY   float64
	attackModle  bool
	countKiller  int
}

var killer2Image *ebiten.Image

func init() {
	img, _, err := ebitenutil.NewImageFromFile("./Resources/killer2.png")
	if err != nil {
		log.Fatalf("Actor.%s", err)
	}
	killer2Image = img
}

// 应该不用上锁
func (k *Killer2) Attack() {
	k.attackModle = true
}
func (k *Killer2) AttackModle() bool {
	return k.attackModle
}

// GetSubImage 获killer的图像
func (k *Killer2) getSubImage() image.Image {
	var img image.Image
	if !k.attackModle {
		k.countKiller = 0
		img = killer2Image.SubImage(image.Rect(0, k.Direction*killer2FrameHeight, 0+killer2FrameWidth, k.Direction*killer2FrameHeight+killer2FrameHeight))
	} else {
		if k.countKiller == k.RefreshRates*4 {
			k.countKiller = 0
			k.attackModle = false
		}
		pixCount := k.countKiller / k.RefreshRates
		img = killer2Image.SubImage(image.Rect(pixCount*killer2FrameWidth, k.Direction*killer2FrameHeight, pixCount*killer2FrameWidth+killer2FrameWidth, k.Direction*killer2FrameHeight+killer2FrameHeight))
		k.countKiller++
	}
	return img
}

// SetMove x & y 是killer当前的位置
func (k *Killer2) SetMove(x, y float64) {
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

// GetPosition 获得killer的位置
func (k *Killer2) getPosition() (float64, float64) {
	k.PosX = k.movX*k.Speed + k.PosX
	k.PosY = k.movY*k.Speed + k.PosY
	return k.PosX, k.PosY
}

func (k *Killer2) SelfUpdate(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(k.Scale, k.Scale)
	/*
		更新killer的位置
	*/
	x, y := k.getPosition()
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(k.getSubImage().(*ebiten.Image), op)
}
func (k *Killer2) HitArea(x, y float64) bool {
	rigth := k.PosX + 70*k.Scale // 80
	left := k.PosX - 70*k.Scale
	up := k.PosY - 40*k.Scale // 80
	down := k.PosY + 40*k.Scale
	return (x < rigth && x > left && y > up && y < down) && k.AttackModle()
}
