package generator

type Config struct {
	Function Function
	Scale    Scale
}

var DefaultConfig = Config{
	Function: Sin,
	Scale:    Small,
}

type Function byte

const (
	Sin = Function(iota)
	Sawtooth
	Random
)

func (s Function) String() string {
	switch s {
	case Sin:
		return "SIN"
	case Sawtooth:
		return "SAW"
	case Random:
		return "RNG"
	default:
		return "???"
	}
}

func (Function) Options() []Function {
	return []Function{Sin, Sawtooth, Random}
}

type Scale byte

const (
	Small = Scale(iota)
	Medium
	Large
)

func (s Scale) String() string {
	switch s {
	case Small:
		return "SML"
	case Medium:
		return "MED"
	case Large:
		return "LRG"
	default:
		return "???"
	}
}

func (Scale) Options() []Scale {
	return []Scale{Small, Medium, Large}
}
