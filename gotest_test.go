package main

import (
	"fmt"
	"testing"
	"os"
	"sync"
	"time"
)

var moduleFixtures = make([]*Pkg, 0)
var globalWg = &sync.WaitGroup{}

func tearDown() {
	for _, j := range moduleFixtures {
		<- j.channel
	}
	globalWg.Wait()
}

func PackageFixture(fixture func(Yield)) interface{} {
	y := &Pkg{channel: make(chan interface{}), wg: globalWg}
	go fixture(y)
	moduleFixtures = append(moduleFixtures, y)
	result, _ := <- y.channel
	fmt.Println(result)
	return result
}

func UnitFixture(fixture func(Yield)) func()(interface{}, func()) {
	return func() (value interface{}, cleanup func()) {
		y := &Unit{channel: make(chan interface{}), wg: &sync.WaitGroup{}}
		go fixture(y)
		value, _ = <- y.channel
		cleanup = func() {
			<- y.channel
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

type Unit struct {
	channel chan interface{}
	wg *sync.WaitGroup
}

func (u *Unit) Return(value interface{}) {
	u.wg.Add(1)
	u.channel <- value
	u.channel <- nil
}

func (u *Unit) Done() {
	u.wg.Done()
}

type Yield interface {
	Return(interface{})
	Done()
}

type Pkg struct {
	channel chan interface{}
	wg *sync.WaitGroup
}

func (j *Pkg) Return(value interface{}) {
	j.wg.Add(1)
	j.channel <- value
	j.channel <- nil
}

func (j *Pkg) Done() {
	j.wg.Done()
}

var fix = PackageFixture(func(y Yield) {
	defer y.Done()
	fmt.Println("setup")
	y.Return(5)
	time.Sleep(time.Second*1)
	fmt.Println("teardown")
})

func TestSimple(t *testing.T) {
	_, c := bazUnit()
	defer c()
	x := fix
	fmt.Println(x)
}

func TestMain(m *testing.M) {
	FixtureMain(m)
}

var bazUnit = UnitFixture(func(j Yield) {
	defer j.Done()
	fmt.Println("baz start")
	j.Return(nil)
	time.Sleep(time.Second * 3)
	fmt.Println("baz cleanup")
})