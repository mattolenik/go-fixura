package main

import (
	"log"
	"os"
	"testing"
	"time"
	. "github.com/mattolenik/gotest/fixture"
)

var logger = log.New(os.Stdout, "", 0)

var fix = PackageFixture(func(yield Y) {
	logger.Println("setup")
	yield.Fixture(5)
	logger.Println("teardown")
	time.Sleep(time.Second*2)
})

func TestSimple(t *testing.T) {
	baz, bazDone := bazUnit()
	defer bazDone()
	logger.Println(baz)
	x := fix
	logger.Println(x)
}

func TestMain(m *testing.M) {
	FixtureMain(m)
}

type Y = *Yield

var bazUnit = UnitFixture(func(yield Y) {
	logger.Println("baz start")
	yield.Fixture(4)
	logger.Println("baz cleanup")
})
