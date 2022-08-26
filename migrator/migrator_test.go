package migrator_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/anton7r/mgx/migrator"
)

func TestTimePrintParse(t *testing.T) {
	for i := 0; i < 99; i++ {
		t.Error(migrator.PrintTime(time.Now()))
	}
}

func niceFormat(stamp time.Time) string {
	return strconv.FormatInt(stamp.UnixMilli(), 36)
}

func niceParse(id string) (int64, error) {
	return strconv.ParseInt(id, 36, 0)
}

//Is a larger than b
func niceIsLargerThanParse(a string, b string) bool {
	aLen := len(a)
	bLen := len(b)

	if aLen > bLen {
		return true
	} else if aLen < bLen {
		return false
	}

	return a > b
}

func TestNiceFormatCompare(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour * 3)

	if (now.UnixMilli() > later.UnixMilli()) != (niceIsLargerThanParse(niceFormat(now), niceFormat(later))) {
		t.Error("Did not sort in correct order")
	}

}

func TestPrintingTimeDifferently(t *testing.T) {
	t.Error(niceFormat(time.Now()))
}

func BenchmarkTimePrint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		migrator.PrintTime(time.Now())
	}
}

func BenchmarkTimePrint2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now().UnixMilli()
	}
}
func BenchmarkTimeParse2(b *testing.B) {
	stam := strconv.FormatInt(time.Now().UnixMilli(), 10)

	for i := 0; i < b.N; i++ {
		strconv.Atoi(stam)
	}
}

func BenchmarkTimePrint3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		niceFormat(time.Now())
	}
}

var unformatted = strconv.FormatInt(time.Now().UnixMilli(), 10)
var unformatted2 = strconv.FormatInt(time.Now().Add(time.Hour).Add(time.Second).Add(time.Minute).UnixMilli(), 10)

var formatted = niceFormat(time.Now())
var formatted2 = niceFormat(time.Now().Add(time.Hour).Add(time.Second).Add(time.Minute))

func BenchmarkTimeParse3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		niceParse(formatted)
	}
}

func BenchmarkCompare2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		niceIsLargerThanParse(unformatted, unformatted2)
	}
}

//This way we do not need to parse the string as the value is there already
func BenchmarkCompare3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		niceIsLargerThanParse(formatted, formatted2)
	}
}

//we can do fast unparsing if we sort them alphabetically but before that we check if the length match and give the longer length one better
// This shall be explored further as BenchmarkCompare3 had good results indicating that it could be used as
