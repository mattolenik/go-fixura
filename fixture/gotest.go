package fixture

import (
	"sync"
	"os"
	"testing"
)

var packageFixtures = make([]*Yield, 0)
var packageWg = &sync.WaitGroup{}

func PackageFixture(fixture func(*Yield)) interface{} {
	y := &Yield{channel: make(chan interface{}), wg: packageWg}
	go func() {
		defer y.Done()
		fixture(y)
	}()
	packageFixtures = append(packageFixtures, y)
	result, _ := <-y.channel
	return result
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

func FixtureMain(m *testing.M) {
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func tearDown() {
	for _, j := range packageFixtures {
		<-j.channel
	}
	packageWg.Wait()
}

type Yield struct {
	channel chan interface{}
	wg      *sync.WaitGroup
}

func (u *Yield) Fixture(value interface{}) {
	u.wg.Add(1)
	u.channel <- value
	u.channel <- nil
}

func (u *Yield) Done() {
	u.wg.Done()
}