package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"sync"
	"time"

	actors "github.com/VegetableManII/Zombies/actor"
	"github.com/VegetableManII/Zombies/utils"
)

const (
	screenWidth  = 660
	screenHeight = 220
)

var killer actors.Killer
var killer2 actors.Killer2
var zombies []*actors.Zombie
var generateChan chan *actors.Zombie

type Game struct {
	// count int
	zombieNumbers int
	mu            sync.RWMutex
}

func init() {
	zombies = make([]*actors.Zombie, 0, 16)
	generateChan = make(chan *actors.Zombie, 0)
	// 随机位置间隔相同时间生成僵尸
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			rand.Seed(time.Now().Unix())
			x, y := rand.Intn(screenWidth), rand.Intn(screenHeight)
			z := &actors.Zombie{PosX: float64(x), PosY: float64(y), Target: (x + y) % 2} // 随机锁定一个敌人
			generateChan <- z
		}
	}()
	fmt.Printf("游戏初始化...\n")
}

func (g *Game) Update() error {
	// Update函数和Draw 是串行执行，不加default会阻塞
	// Update和Draw中对 Actors 对象实体的操作不需要加锁
	select {
	case zomb := <-generateChan:
		g.mu.Lock()
		zombies = append(zombies, zomb)
		g.mu.Unlock()
	default:
	}
	/*
		读取键盘输入
	*/
	// P1 输入
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		actors.SetZombieSpeed(actors.GetZombieSpeed() + 0.01)
	}
	x, y := 0, 0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		x = 1
	}
	killer.SetMove(float64(x), float64(y))
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		killer.Attack()
	}
	// P2 输入
	x, y = 0, 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		x = 1
	}
	killer2.SetMove(float64(x), float64(y))
	if ebiten.IsKeyPressed(ebiten.KeyNumpad0) {
		killer2.Attack()
	}

	/* 根据killer的位置更新僵尸的走向 */
	// TODO x,y := killer.Position()
	for i := range zombies {
		if killer.HitArea(zombies[i].PosX, zombies[i].PosY) {
			zombies[i].Dead()
		}
		if killer2.HitArea(zombies[i].PosX, zombies[i].PosY) {
			zombies[i].Dead()
		}
		zombies[i].SetMove(killer.PosX, killer.PosY, killer2.PosX, killer2.PosY)
	}
	return nil
}
func (g *Game) updateZombies(screen *ebiten.Image) {
	/*
		更新每一只僵尸的位置
	*/
	g.mu.RLock()
	if num := len(zombies); num != 0 {
		for i := range zombies {
			zombies[i].SelfUpdate(screen)
		}
	}
	g.mu.RUnlock()
}

func (g *Game) Draw(screen *ebiten.Image) {
	utils.BackgroundUpdate(screen, float64(screenWidth), float64(screenHeight))
	utils.FrontUpdate(screen, actors.GetZombieSpeed())
	g.updateZombies(screen)
	killer.SelfUpdate(screen)
	killer2.SelfUpdate(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	fmt.Printf("游戏开始...\n")

	// 放置killer
	killer.PosX = screenWidth/2 - 4
	killer.PosY = screenHeight / 2
	killer.Speed = 2
	killer.Scale = 0.5

	killer2.PosX = screenWidth/2 + 4
	killer2.PosY = screenHeight / 2
	killer2.Speed = 2
	killer2.RefreshRates = 4
	killer2.Scale = 0.4

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Zombies~~")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
