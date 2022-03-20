package main

type position struct {
	prc    float32
	margin int
}

func (p position) getCoordinate(max int) int {
	// value = prc * MAX + abs
	return int(p.prc*float32(max)) - p.margin
}

type viewPosition struct {
	x0, y0, x1, y1 position
}

func (vp viewPosition) getCoordinates(maxX, maxY int) (int, int, int, int) {
	var x0 = vp.x0.getCoordinate(maxX)
	var y0 = vp.y0.getCoordinate(maxY)
	var x1 = vp.x1.getCoordinate(maxX)
	var y1 = vp.y1.getCoordinate(maxY)
	return x0, y0, x1, y1
}

const (
	outputView = "output"
	inputView  = "input"
)

var viewPositions = map[string]viewPosition{
	outputView: {
		position{0.0, 0},
		position{0.0, 0},
		position{1.0, 2},
		position{0.75, 2},
	},
	inputView: {
		position{0.0, 0},
		position{0.73, 0},
		position{1.0, 2},
		position{1.0, 2},
	},
}
