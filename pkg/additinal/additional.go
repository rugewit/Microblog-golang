package additinal

import (
	"log"
	"time"
)

func Timer(name string) func() {
	start := time.Now()
	return func() {
		//fmt.Printf("%s took %v\n", name, time.Since(start))
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}
