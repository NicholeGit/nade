package command

import (
	"testing"
)

func TestRunSuite(t *testing.T) {
	SetUp(t)
	defer TearDown()
	// Convey("command 测试", t, nil)
	runCase(t, appStartAndStop)
}

func runCase(t *testing.T, testCase func(*testing.T)) {
	Before()
	defer After()
	testCase(t)
}
