package actor

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	killerFrameWidth  = 70
	killerFrameHeight = 70
	killerFrameNum    = 4
)

type Killer struct {
	PosX, PosY         int
	Speed              int
	killerOX, killerOY int

	movX, movY int
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
func (k *Killer) SetMove(x, y int) {
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
func (k *Killer) getPosition() (int, int) {
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
