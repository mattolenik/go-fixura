package fixura

import (
	"log"
	"os"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

var logger = log.New(os.Stdout, "", 0)

var fix = Fixture(func(yield Y) {
	yield.Fixture(5)
	logger.Println("teardown")
	time.Sleep(time.Second * 2)
}).(int)

func TestSimple(t *testing.T) {
	jig := &struct{ Baz int; FooFix string }{}
	defer LoadFixtures(jig)()
	logger.Println(jig.Baz)
	logger.Println(jig.FooFix)
	x := fix
	logger.Println(x)
}

func TestMain(m *testing.M) {
	GoTestMain(m)
}

var _ = UnitFixture("Baz", func(yield Y, args ...interface{}) {
	logger.Println("baz start")
	yield.Fixture(4)
	logger.Println("baz cleanup")
})

var _ = UnitFixture("FooFix", func(yield Y) {
	logger.Println("FooFix start")
	yield.Fixture("foof")
	logger.Println("FooFix stop")
})

func TestNew(t *testing.T) {
	jig := &struct{ FooFix string }{}
	defer LoadFixtures(jig)()
	assert.Equal(t, "foof", jig.FooFix)
}