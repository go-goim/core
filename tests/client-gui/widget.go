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
	name           string
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
	friendsView = "friends"
	outputView  = "output"
	inputView   = "input"
)

var (
	friend = viewPosition{
		friendsView,
		position{0.0, 0},
		position{0.0, 0},
		position{0.25, 1},
		position{1.0, 1},
	}

	outpu = viewPosition{
		outputView,
		position{0.25, 0},
		position{0.0, 0},
		position{1.0, 1},
		position{0.72, 1},
	}

	input = viewPosition{
		inputView,
		position{0.25, 0},
		position{0.73, 0},
		position{1.0, 1},
		position{1.0, 1},
	}
)
