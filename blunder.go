package passport

import "github.com/georgi-georgiev/blunder"

func NewBlunder() *blunder.Blunder {
	return blunder.NewRFC()
}
