package somepacakge_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

// This snippet shows how to create and use golden-files. They helps us to test complex output without manually hardcoding it
// Very scalable way to test complex structures (write a String() method)
//
// = How to use golden files? =
// * generate golden files by running test with `-update-golden` flag (`go test -update-golden`)
// * check manually generated data in a files and if it's correct then commit it.

var updateGolden = flag.Bool("update-golden", false, "update golden files")

func TestAdd(t *testing.T) {
	// ... table (probably!)
	for _, tc := range сases {
		actual := doSomething(tc)
		golden := filepath.Join("test-fixtures", tc.Name+".golden")
		if *updateGolden {
			ioutil.WriteFile(golden, actual, 0644)
		}
		expected, _ := ioutil.ReadFile(golden)
		if !bytes.Equal(actual, expected) {
			// FAIL!
		}
	}
}
