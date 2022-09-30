package actor

import (
	"github.com/VegetableManII/Zombies/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	_ "image/png"
	"log"
	"math"
)

const (
	zombie0FrameWidth  = 32
	zombie0FrameHeight = 48

	zombie1FrameWidth  = 64
	zombie1FrameHeight = 64
)

type Zombie struct {
	PosX, PosY            float64
	countZombie, zombie0Y int
	movX, movY            float64
	dead                  bool
}

var zombieImage *ebiten.Image
var zombieDiedImage *ebiten.Image
var step float64
var zomRefreshRate int = 20

// SetZombieSpeed 设置僵尸的移动速度
func SetZombieSpeed(speed float64) {
	step = speed
}
func GetZombieSpeed() float64 {
	return step
}
func init() {
	img, _, err := ebitenutil.NewImageFromFile("./Resources/zombie.png")
	if err != nil {
		log.Fatalf("Actor.%s", err)
	}
	zombieImage = img
	img, _, err = ebitenutil.NewImageFromFile("./Resources/zombie_died.png")
	if err != nil {
		log.Fatalf("Actor.%s", err)
	}
	zombieDiedImage = img
	SetZombieSpeed(0.5)
}
func (z *Zombie) Dead() {
	z.dead = true
	utils.HitSound()
}

// SetMove x & y 是killer当前的位置
func (z *Zombie) SetMove(x, y float64) {
	if z.dead {
		z.movX = 0
		z.movY = 0
		return
	}
	// x,y 为左上角
	// x,y 设置为距离 killer 的中心点
	x, y = x+10, y+5
	px := step
	if y < z.PosY {
		z.movY = -px
		z.zombie0Y = 3
	} else {
		z.movY = px
		z.zombie0Y = 0
	}
	// 优化僵尸与猎手处于同一个Y轴的时候的僵尸朝向
	if math.Abs(x-z.PosX) < 5.0 {
		return
	} else {
		if x < z.PosX {
			z.movX = -px
			z.zombie0Y = 1
		} else {
			z.movX = px
			z.zombie0Y = 2
		}
	}
	return
}

// GetPosition 获得僵尸此次的运动方向
func (z *Zombie) getPosition() (float64, float64) {
	z.PosX = z.movX + z.PosX
	z.PosY = z.movY + z.PosY
	return z.PosX, z.PosY
}
func (z *Zombie) getSubImage() image.Image {
	// TODO 根据不同的运动方向显示不同的图片
	var img image.Image
	if z.dead {
		z.countZombie = 0
		img = zombieDiedImage.SubImage(image.Rect(z.countZombie*zombie1FrameWidth, z.zombie0Y*zombie1FrameHeight,
			z.countZombie*zombie1FrameWidth+zombie1FrameWidth, z.zombie0Y*zombie1FrameHeight+zombie1FrameHeight))
	} else {
		if z.countZombie == zomRefreshRate*4 {
			z.countZombie = 0
		}
		pixCount := int(z.countZombie / zomRefreshRate)
		img = zombieImage.SubImage(image.Rect(pixCount*zombie0FrameWidth, z.zombie0Y*zombie0FrameHeight,
			pixCount*zombie0FrameWidth+zombie0FrameWidth, z.zombie0Y*zombie0FrameHeight+zombie0FrameHeight))
		z.countZombie++
	}
	return img
}
func (z *Zombie) SelfUpdate(screen *ebiten.Image) {
	x, y := z.getPosition()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(z.getSubImage().(*ebiten.Image), op)
}
