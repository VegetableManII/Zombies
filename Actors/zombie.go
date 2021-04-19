package Actors

import (
	"image"
	_ "image/png"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	zombie0FrameWidth  = 32
	zombie0FrameHeight = 48

	zombie1FrameWidth  = 64
	zombie1FrameHeight = 64
)

type Zombie struct {
	PosX, PosY            int
	countZombie, zombie0Y int
	movX, movY            int
	mu                    sync.RWMutex
	dead                  bool
}

var zombieImage *ebiten.Image
var zombieDiedImage *ebiten.Image
var step float64

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
	SetZombieSpeed(1)
}
func (z *Zombie) Dead() {
	z.dead = true
}
func (z *Zombie) IsDead() bool {
	// 应该不用上锁
	return z.dead
}

// SetMove x & y 是killer当前的位置
func (z *Zombie) SetMove(x, y int) {
	if z.dead {
		z.movX = 0
		z.movY = 0
		return
	}
	// x,y 为左上角
	// x,y 设置为距离 killer 的中心点
	x, y = x+10, y+5
	px := int(step)
	z.mu.Lock()
	defer z.mu.Unlock()
	if y < z.PosY {
		z.movY = -px
		z.zombie0Y = 3
	} else {
		z.movY = px
		z.zombie0Y = 0
	}
	if x < z.PosX {
		z.movX = -px
		z.zombie0Y = 1
	} else {
		z.movX = px
		z.zombie0Y = 2
	}
	return
}

// GetPosition 获得僵尸此次的运动方向
func (z *Zombie) GetPosition() (int, int) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	z.PosX = z.movX + z.PosX
	z.PosY = z.movY + z.PosY
	return z.PosX, z.PosY
}
func (z *Zombie) GetSubImage() image.Image {
	// TODO 根据不同的运动方向显示不同的图片
	var img image.Image
	if z.dead {
		z.countZombie = 0
		img = zombieDiedImage.SubImage(image.Rect(z.countZombie*zombie1FrameWidth, z.zombie0Y*zombie1FrameHeight,
			z.countZombie*zombie1FrameWidth+zombie1FrameWidth, z.zombie0Y*zombie1FrameHeight+zombie1FrameHeight))
		z.countZombie++
	} else {
		if z.countZombie == 4 {
			z.countZombie = 0
		}
		img = zombieImage.SubImage(image.Rect(z.countZombie*zombie0FrameWidth, z.zombie0Y*zombie0FrameHeight,
			z.countZombie*zombie0FrameWidth+zombie0FrameWidth, z.zombie0Y*zombie0FrameHeight+zombie0FrameHeight))
		z.countZombie++
	}
	return img
}
