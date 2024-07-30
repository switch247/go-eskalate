package routinestest_test

import (
	"fmt"
	"time"
)

var dbData = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

func RunRoutine() {
	t0 := time.Now()
	for i := 0; i < len(dbData); i++ {
		continue
	}
	fmt.Printf("\n Total execution time %v", time.Since(t0))
}

func dbCall(i int) {

}
