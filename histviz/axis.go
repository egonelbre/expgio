package main

type Axis struct {
	Data   Range
	Screen Range
}

type Range struct {
	Min, Max float64
}

func (axis *Axis) ToScreen(v float64) float64 {
	return lerp(invlerp(v, axis.Data.Min, axis.Data.Max), axis.Screen.Min, axis.Screen.Max)
}

func (axis *Axis) ToData(v float64) float64 {
	return lerp(invlerp(v, axis.Screen.Min, axis.Screen.Max), axis.Data.Min, axis.Data.Max)
}
