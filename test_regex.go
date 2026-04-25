package main
import (
    "fmt"
    "regexp"
)
func main() {
    re := regexp.MustCompile(`^V\d+A\d+$`)
    fmt.Println(re.MatchString("V1A1\n"))
    fmt.Println(re.MatchString("V1A1\r\n"))
}
