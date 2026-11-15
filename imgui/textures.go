package imgui

// #include "structs_accessor.h"
import "C"

import "unsafe"

// AppendTextures appends the textures the draw data refers to, each its own
// TextureData, to dst and returns it, so a backend can handle the requests
// of every one; passing the slice back each frame reuses it.
// Textures().Slice() gives only the first one right: its Vector points at
// one Go value, not at the C list of pointers (see #492).
func (self *DrawData) AppendTextures(dst []TextureData) []TextureData {
	selfArg, selfFin := self.Handle()
	defer selfFin()

	list := C.wrap_ImDrawData_GetTextures(selfArg)
	if list == nil || list.Size <= 0 || list.Data == nil {
		return dst
	}
	for _, tex := range unsafe.Slice(list.Data, list.Size) {
		dst = append(dst, TextureData{CData: tex})
	}
	return dst
}

// PixelBytes returns the texture's pixels in ImGui's memory, row after row,
// BytesPerPixel bytes a pixel, or nil when it has none. Unlike Pixels, it
// needs no conversion from uintptr, which go vet reports as a possible misuse
// of unsafe.Pointer. The bytes stay valid while ImGui keeps the pixels, until
// the texture is destroyed or created again.
func (self *TextureData) PixelBytes() []byte {
	selfArg, selfFin := self.Handle()
	defer selfFin()

	if selfArg.Pixels == nil {
		return nil
	}
	n := int(selfArg.Width) * int(selfArg.Height) * int(selfArg.BytesPerPixel)
	return unsafe.Slice((*byte)(unsafe.Pointer(selfArg.Pixels)), n)
}
