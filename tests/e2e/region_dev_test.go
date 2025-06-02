//go:build dev

package e2e_test

import (
	"fmt"
	"os"
)

func init() {
	user := os.Getenv("USER")
	RegionName = fmt.Sprintf("https://%s.formance.dev", user)
	fmt.Println("RegionName set to:", RegionName)
}
