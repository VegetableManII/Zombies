package main

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/VegetableManII/mygame/Actors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 660
	screenHeight = 220

	// frameOX     = 0
	// frameOY     = 0
)

// var (
// 	runnerImage *ebiten.Image
// )
var killer Actors.Killer
var zombies []*Actors.Zombie
var generateChan chan *Actors.Zombie

var bg *ebiten.Image
var pressStart2pFont font.Face
var ctx *audio.Context
var hit []byte

type Game struct {
	// count int
	zombieNumbers int
	mu            sync.RWMutex
}

func init() {
	zombies = make([]*Actors.Zombie, 0, 16)
	generateChan = make(chan *Actors.Zombie, 0)
	fi, err := os.Stat("./Resources/music/kill.mp3")
	if err != nil {
		log.Fatalf("main.%s", err)
	}
	hit = make([]byte, fi.Size())
	// 加载资源
	f, err := os.Open("./Resources/music/kill.mp3")
	if err != nil {
		log.Fatalf("main.%s", err)
	}
	ctx = audio.NewContext(44100)
	s, err := mp3.Decode(ctx, f)
	if err != nil {
		log.Fatalf("main.%s", err)
	}
	hit, err = ioutil.ReadAll(s)
	if err != nil {
		log.Fatalf("main.%s", err)
	}

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 32
	pressStart2pFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    28,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := ebitenutil.NewImageFromFile("./Resources/background0.jpg")
	if err != nil {
		log.Fatalf("main.bacckground init err! %s", err)
	}
	bg = img
	// 随机位置间隔相同时间生成僵尸
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			rand.Seed(time.Now().Unix())
			x, y := rand.Intn(screenWidth), rand.Intn(screenHeight)
			z := &Actors.Zombie{PosX: x, PosY: y}
			generateChan <- z
		}
	}()
	fmt.Printf("游戏初始化...\n")
}

func (g *Game) Update() error {
	// g.count++
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
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		Actors.SetZombieSpeed(Actors.GetZombieSpeed() + 0.01)
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
	killer.SetMove(x, y)
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		killer.Attack()
	}

	/* 根据killer的位置更新僵尸的走向 */
	// TODO x,y := killer.Position()
	for i := range zombies {
		xrange := math.Abs(float64(killer.PosX + 10 - zombies[i].PosX))
		yrange := math.Abs(float64(killer.PosY + 5 - zombies[i].PosY))
		if xrange < 10.0 && yrange < 20.0 && killer.AttackModle() {
			zombies[i].Dead()
			p := audio.NewPlayerFromBytes(ctx, hit)
			p.Play()
		}
		zombies[i].SetMove(killer.PosX, killer.PosY)
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
			x, y := zombies[i].GetPosition()
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.5, 0.5)
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(zombies[i].GetSubImage().(*ebiten.Image), op)
		}
	}
	g.mu.RUnlock()
}

func (g *Game) updateKiller(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	/*
		更新killer的位置
	*/
	x, y := killer.GetPosition()
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(killer.GetSubImage().(*ebiten.Image), op)
}
func (g *Game) updateBackground(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screenWidth/1920.0), float64(screenHeight/640.0))
	screen.DrawImage(bg, op)
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.updateBackground(screen)
	fps := fmt.Sprintf("FPS:%0.2f Speed:%0.2fpx/tick\nWASD move J attack", ebiten.CurrentFPS(), Actors.GetZombieSpeed())
	text.Draw(screen, fps, pressStart2pFont, 0, 20, color.Black)
	g.updateZombies(screen)
	g.updateKiller(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	fmt.Printf("游戏开始...\n")

	// 放置killer
	killer.PosX = screenWidth / 2
	killer.PosY = screenHeight / 2
	killer.Speed = 2

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Zombies~~")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
