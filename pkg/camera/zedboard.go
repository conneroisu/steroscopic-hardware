package camera

import "image"

// ZedBoardCamera represents a ZedBoard camera
type ZedBoardCamera struct {
}

// NewZedBoardCamera creates a new ZedBoard camera
func NewZedBoardCamera(comPort string) *ZedBoardCamera {
	return &ZedBoardCamera{}
}

var _ Camer = (*ZedBoardCamera)(nil)

// :GoImpl z *ZedBoardCamera camera.Camer

// Stream streams the camera
func (z *ZedBoardCamera) Stream(_ chan *image.Gray) {
	panic("not implemented") // TODO: Implement
}

// Close closes the camera
func (z *ZedBoardCamera) Close() error {
	panic("not implemented") // TODO: Implement
}
