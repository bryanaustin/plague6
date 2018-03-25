package distributor

import (
	"math/rand"
	"time"
)

var rander *rand.Rand

func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Distributor struct {
	Options []Option
	total   uint
}

type Option struct {
	Item    interface{}
	Target  float32
	current float32
	total   uint
}

func (d *Distributor) Pick() (item interface{}) {
	d.Normalize()
	s := rander.Float32()
	total := float32(0.0)
	for _, o := range d.Options {
		item = o.Item
		total += o.current
		if s <= total {
			d.total++
			o.total++
			return
		}
	}
	return
}

func (d *Distributor) Normalize() {
	if d.total == uint(0) {
		for _, o := range d.Options {
			o.current = o.Target
		}
		return
	}

	gt := float64(d.total)
	on := float64(1.0) / gt
	for _, o := range d.Options {
		actual := float32(float64(o.total) * on)
		diff := o.Target - actual
		o.current = o.Target + (diff * float32(-1.0))
	}
}
