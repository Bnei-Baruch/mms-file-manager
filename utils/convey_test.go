package utils_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/Bnei-Baruch/mms-file-manager/utils"
	"time"
	"github.com/smartystreets/assertions"
)

func TestConveyEventuallySpec(t *testing.T) {
	Convey("Eventually", t, func() {
		Convey("When timeout is reached the assertion should fail", func() {
			res := utils.Eventually(func() interface{} {
				return false
			}, 1 * time.Second, assertions.ShouldBeTrue)
			So(res, ShouldNotBeEmpty)
		})

		Convey("When timeout is reached but the function still in process the assertion should fail", func() {
			res := utils.Eventually(func() interface{} {
				time.Sleep(2 * time.Second)
				return true
			}, 1 * time.Second, assertions.ShouldBeTrue)
			So(res, ShouldNotBeEmpty)
		})

		Convey("When timeout is not reached the assertion should not fail", func() {
			res := utils.Eventually(func() interface{} {
				return true
			}, 1 * time.Second, assertions.ShouldBeTrue)
			So(res, ShouldBeEmpty)
		})

		Convey("If assertion succeeds immediatlly then the function is called exactlly once and should suuceed", func() {
			var counter = 0
			res := utils.Eventually(func() interface{} {
				counter ++
				return true
			}, 1 * time.Second, assertions.ShouldBeTrue)
			So(res, ShouldBeEmpty)
			So(counter, ShouldEqual, 1)
		})

		Convey("assertion with parameters should work", func() {
			res := utils.Eventually(func() interface{} {
				return 10
			}, 1 * time.Second, assertions.ShouldEqual, 10)

			So(res, ShouldBeEmpty)
		})

		Convey("When timeout is shorter than polling duration it should fail", func() {
			res := utils.Eventually(func() interface{} {
				return true
			}, 5 * time.Millisecond, assertions.ShouldBeTrue)

			So(res, ShouldNotBeEmpty)
		})
	})
}
