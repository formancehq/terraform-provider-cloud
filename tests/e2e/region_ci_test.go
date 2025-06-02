//go:build ci

package e2e_test

import "fmt"

func init() {
	RegionName = "staging"
	fmt.Println("RegionName set to:", RegionName)
}
