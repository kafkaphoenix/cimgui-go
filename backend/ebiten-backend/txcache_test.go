package ebitenbackend_test

import (
	"testing"

	ebitenbackend "github.com/AllenDang/cimgui-go/backend/ebiten-backend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestTextureCacheNextIdSkipsRemovedIDs(t *testing.T) {
	cache := ebitenbackend.NewCache()
	a := imgui.TextureID(cache.NextId())
	cache.SetTexture(a, ebiten.NewImage(1, 1))
	b := imgui.TextureID(cache.NextId())
	cache.SetTexture(b, ebiten.NewImage(1, 1))
	cache.RemoveTexture(a)
	if c := imgui.TextureID(cache.NextId()); c == a || c == b {
		t.Errorf("NextId gives %d after %d and %d, with %d removed", c, a, b, a)
	}
}

type game struct{}

func (game) Update() error              { return nil }
func (game) Draw(*ebiten.Image)         {}
func (game) Layout(w, h int) (int, int) { return w, h }

// A game texture removed from the cache, as DeleteTexture does, is no longer
// a game: ForEachGame, which Update and Draw run, must skip it.
func TestTextureCacheRemoveTextureDropsGame(t *testing.T) {
	cache := ebitenbackend.NewCache()
	id := imgui.TextureID(cache.NextId())
	cache.SetTexture(id, ebiten.NewImage(1, 1))
	cache.SetGameTexture(id, game{})
	cache.RemoveTexture(id)

	if _, ok := cache.GetGameTexture(id); ok {
		t.Fatalf("texture %d is still a game after its removal", id)
	}
	cache.ForEachGame(func(id imgui.TextureID, _ ebiten.Game, _ *ebiten.Image) {
		t.Errorf("ForEachGame runs removed texture %d", id)
	})
}
