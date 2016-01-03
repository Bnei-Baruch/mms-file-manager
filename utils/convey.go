package utils

import (
	"time"
)

//Eventually wraps an actual value allowing assertions to be made on it.
//The assertion is tried periodically until it passes or a timeout occurs.
//
//The timeout interval is configurable as the first argument:
//
//Both intervals can either be specified as time.Duration, parsable duration strings or as floats/integers.  In the
//last case they are interpreted as seconds.
//
//If Eventually is passed an actual that is a function taking no arguments and returning at least one value,
//then Eventually will call the function periodically and try the matcher against the function's first return value.
//
//Example:
//
//    Eventually(func() int {
//        return thingImPolling.Count()
//    }).Should(BeNumerically(">=", 17))
//
//Note that this example could be rewritten:
//
//    Eventually(thingImPolling.Count).Should(BeNumerically(">=", 17))
//
//If the function returns more than one value, then Eventually will pass the first value to the matcher and
//assert that all other values are nil/zero.
//This allows you to pass Eventually a function that returns a value and an error - a common pattern in Go.
//
//For example, consider a method that returns a value and an error:
//    func FetchFromDB() (string, error)
//
//Then
//    Eventually(FetchFromDB).Should(Equal("hasselhoff"))
//
//Will pass only if the the returned error is nil and the returned string passes the matcher.
//
//Eventually's default timeout is 1 second, and its default polling interval is 10ms
func Eventually(actual interface{}, expected ...interface{}) (res string) {
	assertionFunc := expected[1].(func(interface{}, ...interface{}) string)
	expectedParams := make([]interface{}, 0)
	if len(expected) > 2 {
		expectedParams = expected[2:]
	}

	actualFunc := actual.(func() interface{})

	timeoutInterval := expected[0].(time.Duration)
	pollingInterval := time.Millisecond * 10

	if pollingInterval >= timeoutInterval {
		return "Timeout interval is too short"
	}

	ticker := time.NewTicker(pollingInterval)
	done := make(chan bool)
	defer func() {
		ticker.Stop()
		close(done)
	}()

	go func() {
		time.Sleep(timeoutInterval)
		res = "Timeout reached without successful assertion"
		done <- true
	}()

	go func() {
		for {
			actualVal := actualFunc()
			if res = assertionFunc(actualVal, expectedParams...); res == "" {
				done <- true
			}
			select {
			case <-done:
				return
			case <-ticker.C:
			}
		}
	}()

	<-done
	return
}

