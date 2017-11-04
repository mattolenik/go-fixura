package main

import (
	. "github.com/mattolenik/gotest/fixture"
	"log"
	"os"
	"testing"
	"time"
)

var logger = log.New(os.Stdout, "", 0)

var fix = Fixture(func(yield Y) {
	logger.Println("setup")
	yield.Fixture(5)
	logger.Println("teardown")
	time.Sleep(time.Second * 2)
})().(int)

func TestSimple(t *testing.T) {
	baz, bazDone := bazUnit()
	defer bazDone()
	logger.Println(baz)
	x := fix
	logger.Println(x)
}

func TestMain(m *testing.M) {
	GoTestMain(m)
}

var bazUnit = UnitFixture(func(yield Y) {
	logger.Println("baz start")
	yield.Fixture(4)
	logger.Println("baz cleanup")
})
