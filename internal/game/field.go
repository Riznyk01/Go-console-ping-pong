package game

import ping_pong "ping-pong"

type Field struct {
	Console [][]rune
}

func NewScreen(cfg *ping_pong.Config) *Field {
	scr := make([][]rune, cfg.WindowHeight)
	for i := 0; i < cfg.WindowHeight; i++ {
		lineRow := make([]rune, cfg.WindowWidth)
		lineRow[0] = cfg.RacketsLineFiller
		lineRow[len(lineRow)-1] = cfg.RacketsLineFiller
		for linePx := 1; linePx < len(lineRow)-1; linePx++ {
			lineRow[linePx] = cfg.EmptyFiller
		}
		scr[i] = lineRow
	}
	return &Field{
		Console: scr,
	}
}
