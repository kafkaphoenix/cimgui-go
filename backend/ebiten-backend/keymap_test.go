package ebitenbackend

import (
	"testing"

	"github.com/AllenDang/cimgui-go/imgui"
)

// Sending a frame's keys to ImGui allocates nothing: the backend does it
// every frame.
func TestSendInputAllocatesNothing(t *testing.T) {
	ctx := createContext()
	defer imgui.DestroyContextV(ctx)
	io := currentIO()
	io.SetIniFilename("") // no imgui.ini
	var chars []rune
	if n := testing.AllocsPerRun(100, func() { chars = sendInput(io, chars) }); n != 0 {
		t.Errorf("sending a frame's input allocated %v times, want 0", n)
	}
}
