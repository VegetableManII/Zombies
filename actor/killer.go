package actor

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type KillerBase struct {
	PosX, PosY   float64
	Speed        float64
	Scale        float64
	Direction    int
	RefreshRates int
	movX, movY   float64
	attackModle  bool
	countKiller  int
}

// 应该不用上锁
func (k *KillerBase) Attack() {
	k.attackModle = true
}
func (k *KillerBase) AttackModle() bool {
	return k.attackModle
}

// GetPosition 获得killer的位置
func (k *KillerBase) getPosition() (float64, float64) {
	k.PosX = k.movX*k.Speed + k.PosX
	k.PosY = k.movY*k.Speed + k.PosY
	return k.PosX, k.PosY
}
func (k *KillerBase) getUpdateDrawImageOptions() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(k.Scale, k.Scale)
	/*
		更新killer的位置
	*/
	x, y := k.getPosition()
	op.GeoM.Translate(float64(x), float64(y))
	return op
	// screen.DrawImage(k.getSubImage().(*ebiten.Image), op)
}
