package main

type Axis struct {
	Data   Range
	Screen Range

	Reverse bool
}

type Range struct {
	Min, Max float64
}

func (axis *Axis) ToScreen(v float64) float64 {
	return v
}

func (axis *Axis) ToData(v float64) float64 {
	return v
}
