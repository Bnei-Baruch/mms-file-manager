package logger_test

import (
	"testing"

	_ "github.com/Bnei-Baruch/mms-file-manager/services/logger"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoggerSpec(t *testing.T) {
	Convey("Writing logs", t, func() {
		Convey("writes to screen", func() {
		})
		Convey("writes to file", func() {
		})
		Convey("discards log", func() {
		})
		Convey("adds prefix to output", func() {
		})
	})
}
