package main

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Board struct {
	Size image.Point // number of columns and rows

	// ...
}

type BoardStyle struct {
	CellSize unit.Dp
	Gap      unit.Dp
	Color    color.NRGBA

	*Board
}

func (board BoardStyle) Layout(gtx layout.Context) layout.Dimensions {
	cellSize := gtx.Metric.Dp(board.CellSize)
	gap := gtx.Metric.Dp(board.Gap)

	size := image.Point{
		X: board.Size.X*cellSize + (board.Size.X-1)*gap,
		Y: board.Size.Y*cellSize + (board.Size.Y-1)*gap,
	}

	cell := image.Point{
		X: cellSize,
		Y: cellSize,
	}

	for cx := 0; cx < board.Size.X; cx++ {
		for cy := 0; cy < board.Size.Y; cy++ {
			corner := image.Point{
				X: cx * (cellSize + gap),
				Y: cy * (cellSize + gap),
			}

			r := clip.UniformRRect(image.Rectangle{
				Min: corner,
				Max: corner.Add(cell),
			}, gap)

			paint.FillShape(gtx.Ops, board.Color, r.Op(gtx.Ops))
		}
	}

	return layout.Dimensions{
		Size: size,
	}
}
