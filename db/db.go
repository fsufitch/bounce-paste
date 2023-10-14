package db

import "fmt"

var Something = "this comes from the db"

func GetSomeText() string {
	return fmt.Sprintf("foo %s", "bar")
}
