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
	killerFrameNum    = 4
)

type Killer struct {
	PosX, PosY         float64
	Speed              float64
	Scale              float64
	killerOX, killerOY int

	movX, movY float64
}

var killerImage *ebiten.Image
var attackModle bool
var countKiller int

func init() {
	img, _, err := ebitenutil.NewImageFromFile("./Resources/killer.png")
	if err != nil {
		log.Fatalf("Actor.%s", err)
	}
	killerImage = img
	attackModle = false
	countKiller = 0
}

// 应该不用上锁
func (k *Killer) Attack() {
	attackModle = true
}
func (k *Killer) AttackModle() bool {
	return attackModle
}

// GetSubImage 获killer的图像
func (k *Killer) getSubImage() image.Image {
	var img image.Image
	// 感觉不用加锁也行？？
	if !attackModle {
		countKiller = 0
		img = killerImage.SubImage(image.Rect(0, k.killerOY*killerFrameHeight, 0+killerFrameWidth, k.killerOY*killerFrameHeight+killerFrameHeight))
	} else {
		if countKiller == 4 {
			countKiller = 0
			attackModle = false

		}
		img = killerImage.SubImage(image.Rect(countKiller*killerFrameWidth, k.killerOY*killerFrameHeight, countKiller*killerFrameWidth+killerFrameWidth, k.killerOY*killerFrameHeight+killerFrameHeight))
		countKiller++
	}

	return img
}

// SetMove x & y 是killer当前的位置
func (k *Killer) SetMove(x, y float64) {
	k.movX, k.movY = x, y
	if x == -1 {
		k.killerOY = 1
	}
	if x == 1 {
		k.killerOY = 2
	}
	if y == -1 {
		k.killerOY = 3
	}
	if y == 1 {
		k.killerOY = 0
	}

}

// GetPosition 获得killer的位置
func (k *Killer) getPosition() (float64, float64) {
	k.PosX = k.movX*k.Speed + k.PosX
	k.PosY = k.movY*k.Speed + k.PosY
	return k.PosX, k.PosY
}
func (k *Killer) SelfUpdate(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	/*
		更新killer的位置
	*/
	x, y := k.getPosition()
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(k.getSubImage().(*ebiten.Image), op)
}

func (k *Killer) HitArea(x, y float64) bool {
	rigth := k.PosX + 20*k.Scale // 35
	left := k.PosX - 20*k.Scale
	up := k.PosY - 10*k.Scale // 35
	down := k.PosY + 10*k.Scale
	return (x < rigth && x > left && y > up && y < down) && k.AttackModle()
}
