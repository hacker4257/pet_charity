package utils

import (
	"fmt"
	"time"
)

func GenerateOrderNo(prefix string) string {
	return fmt.Sprintf("%s%s%s",
		prefix,
		time.Now().Format("20060102150405"),
		RandomCode(6),
	)
}
