package utils

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var bg *ebiten.Image
var pressStart2pFont font.Face
var ctx *audio.Context
var hit []byte

func init() {
	fi, err := os.Stat("Resources/music/kill.mp3")
	if err != nil {
		log.Fatalf("utils.%s", err)
	}
	hit = make([]byte, fi.Size())
	// 加载资源
	f, err := os.Open("Resources/music/kill.mp3")
	if err != nil {
		log.Fatalf("utils.%s", err)
	}
	ctx = audio.NewContext(44100)
	s, err := mp3.DecodeWithSampleRate(44100, f)
	// s, err := mp3.Decode(ctx, f)
	if err != nil {
		log.Fatalf("utils.%s", err)
	}
	hit, err = ioutil.ReadAll(s)
	if err != nil {
		log.Fatalf("utils.%s", err)
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
	img, _, err := ebitenutil.NewImageFromFile("Resources/background0.jpg")
	if err != nil {
		log.Fatalf("utils.%s", err)
	}
	bg = img
}
func BackgroundUpdate(screen *ebiten.Image, screenWidth, screenHeight float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screenWidth/1920.0), float64(screenHeight/640.0))
	screen.DrawImage(bg, op)
}
func FrontUpdate(screen *ebiten.Image, speed float64) {
	fps := fmt.Sprintf("FPS:%0.2f Speed:%0.2fpx/tick\nWASD move J attack", ebiten.CurrentFPS(), speed)
	text.Draw(screen, fps, pressStart2pFont, 0, 20, color.Black)
}
func HitSound() {
	p := ctx.NewPlayerFromBytes(hit)
	p.Play()
}
