package fixura

import (
	"os"
	"sync"
	"testing"
	"reflect"
)

type Y = *Yield

var fixtures = make([]Y, 0)
var fixtureWg = &sync.WaitGroup{}

func Fixture(fixture func(Y)) interface{} {
	y := &Yield{channel: make(chan interface{}), wg: fixtureWg}
	go func() {
		defer y.Done()
		fixture(y)
	}()
	fixtures = append(fixtures, y)
	result, _ := <-y.channel
	return result
}

var fixturesByName = make(map[string]func(...interface{}) (interface{}, func()))

func UnitFixture(name string, fixtureFunc interface{}) interface{} {
	fix := func(args ...interface{}) (result interface{}, cleanup func()) {
		y := &Yield{channel: make(chan interface{}), wg: &sync.WaitGroup{}}
		go func() {
			defer y.Done()
			switch f := fixtureFunc.(type) {
			case func(Y):
				f(y)
			case func(Y, ...interface{}):
				f(y, args)
			}
		}()
		fixtures = append(fixtures, y)
		result, _ = <-y.channel
		cleanup = func() {
			<-y.channel
			y.wg.Wait()
		}
		return
	}
	fixturesByName[name] = fix
	return nil
}

func LoadFixtures(jig interface{}) func() {
	val := reflect.ValueOf(jig).Elem()
	var cleanupFuncs []func()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		fixture, ok := fixturesByName[typeField.Name]
		if !ok {
			panic("No fixture registered with name " + typeField.Name)
		}
		value, cleanup := fixture()
		cleanupFuncs = append(cleanupFuncs, cleanup)
		valueField.Set(reflect.ValueOf(value))
	}
	return func() {
		for _, f := range cleanupFuncs {
			f()
		}
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
