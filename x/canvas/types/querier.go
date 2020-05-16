package types

const (
	QueryCanvas = "canvas"
	QueryPoints = "points"
)

type QueryCanvasParams struct {
	Id string `json:"id"`
}

type QueryPointsParams struct {
	Id string `json:"id"`
}
