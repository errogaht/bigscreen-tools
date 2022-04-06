package common

import (
	"fmt"
	"time"
)

func LogM(m string) {
	fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), m)
}
