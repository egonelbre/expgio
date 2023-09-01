package generator

import (
	"context"
	"math"
	"math/rand"
	"time"
)

const (
	UpdateDebounce = 500 * time.Millisecond
)

type Client struct {
	Status chan Status
	Data   chan Data

	control chan Config
	active  Config

	initial Config
}

type Status string

type Data struct {
	Status   string
	Values   []Point
	Min, Max Point
}

func NewClient(initial Config) *Client {
	return &Client{
		Status: make(chan Status, 1),
		Data:   make(chan Data, 1),

		control: make(chan Config, 1),

		initial: initial,
	}
}

func (client *Client) InitialConfig() Config {	return client.initial }

func (client *Client) Run(ctx context.Context) {
	client.reconfigure(client.initial)

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return

		case <-tick.C:
			const N = 100

			data := Data{
				Status: "OK",
				Values: make([]Point, N),
			}

			for i := range data.Values {
				var p Point

				now := float64(time.Now().Unix()) / 5
				t := float64(i)/N + now
				switch client.active.Function {
				case Sin:
					p = Point{X: float32(t), Y: float32(math.Sin(t))}
				case Sawtooth:
					p = Point{X: float32(t), Y: float32(math.Mod(t, 1))}
				case Random:
					p = Point{X: rand.Float32(), Y: rand.Float32()}
				}

				switch client.active.Scale {
				case Small:
					p.X *= 0.1
					p.Y *= 0.1
				case Medium:
				case Large:
					p.X *= 10
					p.Y *= 10
				}

				data.Values[i] = p
			}

			sendAndDropOld(client.Data, data)

		case next := <-client.control:

			sendAndDropOld(client.Status, "Waiting for more control")

			// wait for more updates, if there are any
			wait:
			for {
				tick := time.NewTimer(UpdateDebounce)
				select {
				case <-tick.C:
					break wait
				case next = <-client.control:
					tick.Stop()
					continue wait
				}
			}

			// finally update
			client.reconfigure(next)
		}
	}
}

func (client *Client) reconfigure(next Config) {
	// TODO: this could also check what has changed between p.active and next
	// to only update the necessary things
	client.active = next

	sendAndDropOld(client.Status, "Reconfiguring")

	// simulate something slow
	time.Sleep(2 * time.Second)

	sendAndDropOld(client.Status, "Ready")
}

func (client *Client) Update(config Config) {
	sendAndDropOld(client.control, config)
}

type Point struct {
	X, Y float32
}

func sendAndDropOld[C chan T, T any](ch C, v T) {
	// remove the old value from the channel
	select {
	case <-ch:
	default:
	}

	// send a new value
	ch <- v
}
