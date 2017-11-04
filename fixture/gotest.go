package fixture

import (
	"os"
	"sync"
	"testing"
)

type Y = *Yield

var fixtures = make([]*Yield, 0)
var fixtureWg = &sync.WaitGroup{}

func Fixture(fixture func(*Yield)) func() interface{} {
	return func() interface{} {
		y := &Yield{channel: make(chan interface{}), wg: fixtureWg}
		go func() {
			defer y.Done()
			fixture(y)
		}()
		fixtures = append(fixtures, y)
		result, _ := <-y.channel
		return result
	}
}

func UnitFixture(fixture func(*Yield)) func() (value interface{}, cleanup func()) {
	return func() (value interface{}, cleanup func()) {
		y := &Yield{channel: make(chan interface{}), wg: &sync.WaitGroup{}}
		go func() {
			defer y.Done()
			fixture(y)
		}()
		value, _ = <-y.channel
		cleanup = func() {
			<-y.channel
			y.wg.Wait()
		}
		return
	}
}

func GoTestMain(m *testing.M) {
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func tearDown() {
	for _, j := range fixtures {
		<-j.channel
	}
	fixtureWg.Wait()
}

type Yield struct {
	channel    chan interface{}
	wg         *sync.WaitGroup
	hasYielded bool
}

func (y *Yield) Fixture(value interface{}) {
	if y.hasYielded {
		panic("Yield.Fixture can only be called once")
	}
	y.hasYielded = true
	y.wg.Add(1)
	y.channel <- value
	y.channel <- nil
}

func (y *Yield) Done() {
	close(y.channel)
	y.wg.Done()
}
