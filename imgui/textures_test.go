package imgui

import (
	"testing"
	"unsafe"
)

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

func TestTextureDataPixelBytes(t *testing.T) {
	ctx := CreateContext()
	defer DestroyContextV(ctx)
	io := CurrentIO()
	io.SetIniFilename("") // no imgui.ini
	io.SetBackendFlags(io.BackendFlags() | BackendFlagsRendererHasTextures)
	io.SetDisplaySize(Vec2{X: 640, Y: 360})
	NewFrame()
	Render()

	tex := io.Fonts().TexData()
	pix := tex.PixelBytes()
	if len(pix) != int(tex.SizeInBytes()) {
		t.Fatalf("got %d bytes, want %d", len(pix), tex.SizeInBytes())
	}
	if first := uintptr(unsafe.Pointer(&pix[0])); first != tex.Pixels() {
		t.Errorf("first byte at %#x, Pixels at %#x", first, tex.Pixels())
	}
	w, h, bpp := tex.Width(), tex.Height(), uintptr(tex.BytesPerPixel())
	if last := uintptr(unsafe.Pointer(&pix[len(pix)-1])); last != tex.PixelsAt(w-1, h-1)+bpp-1 {
		t.Errorf("last byte at %#x, want the last of the pixel at %#x", last, tex.PixelsAt(w-1, h-1))
	}

	freed := NewTextureData()
	defer freed.Destroy()
	freed.Create(TextureFormatRGBA32, 4, 4)
	freed.DestroyPixels() // keeps the size
	if pix := freed.PixelBytes(); pix != nil {
		t.Errorf("a texture whose pixels were freed gives %d bytes, want nil", len(pix))
	}
}
