package imgui

import "testing"

func TestDrawDataAppendTexturesListsEveryTexture(t *testing.T) {
	ctx := CreateContext()
	defer DestroyContextV(ctx)
	io := CurrentIO()
	io.SetIniFilename("") // no imgui.ini
	io.SetBackendFlags(io.BackendFlags() | BackendFlagsRendererHasTextures)
	io.SetDisplaySize(Vec2{X: 640, Y: 360})

	// frame draws text at size, when not 0, and completes the textures'
	// requests as a backend would.
	frame := func(size float32) []TextureData {
		NewFrame()
		if size != 0 {
			Begin("text")
			PushFont(CurrentFont(), size)
			Text("The quick brown fox jumps over the lazy dog")
			PopFont()
			End()
		}
		Render()
		textures := CurrentDrawData().AppendTextures(nil)
		for _, tex := range textures {
			if tex.Status() == TextureStatusWantCreate {
				tex.SetTexID(TextureID(tex.UniqueID() + 1))
			}
			tex.SetStatus(TextureStatusOK)
		}
		return textures
	}

	first := frame(40)
	if len(first) != 1 || first[0].CData != io.Fonts().TexData().CData {
		t.Fatalf("first frame: got %d textures, want the atlas alone", len(first))
	}
	frame(0)
	frame(0)
	io.Fonts().CompactCache() // drops the unused size, moving the atlas to a new texture
	got := frame(0)
	want := []TextureData{first[0], *io.Fonts().TexData()}
	if len(got) != len(want) {
		t.Fatalf("after the atlas moved: got %d textures, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].CData != want[i].CData {
			t.Errorf("texture %d: got #%d, want #%d", i, got[i].UniqueID(), want[i].UniqueID())
		}
	}
	if n := CurrentPlatformIO().Textures().Size; int(n) != len(got) {
		t.Errorf("got %d textures, ImGui lists %d", len(got), n)
	}
}
