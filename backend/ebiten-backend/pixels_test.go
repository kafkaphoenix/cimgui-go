package ebitenbackend

import (
	"testing"
	"unsafe"

	"github.com/AllenDang/cimgui-go/imgui"
)

// newTextureData wraps a C pointer in a Go value, so it is called through a
// function value too (see render_test.go).
var newTextureData = imgui.NewTextureData

func TestPixelBytes(t *testing.T) {
	ctx := createContext()
	defer imgui.DestroyContextV(ctx)
	io := currentIO()
	io.SetIniFilename("") // no imgui.ini
	io.SetBackendFlags(io.BackendFlags() | imgui.BackendFlagsRendererHasTextures)
	io.SetDisplaySize(imgui.Vec2{X: 640, Y: 360})
	imgui.NewFrame()
	imgui.Render()

	tex := io.Fonts().TexData()
	pix := pixelBytes(*tex)
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

	freed := newTextureData()
	defer freed.Destroy()
	freed.Create(imgui.TextureFormatRGBA32, 4, 4)
	freed.DestroyPixels() // keeps the size
	if pix := pixelBytes(*freed); pix != nil {
		t.Errorf("a texture whose pixels were freed gives %d bytes, want nil", len(pix))
	}
}
