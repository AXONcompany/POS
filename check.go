package main
import (
"fmt"
"github.com/jackc/pgx/v5/pgtype"
)
func main() {
var n pgtype.Numeric
err := n.Scan(2.5)
fmt.Printf("Val: %+v, Err: %v\n", n, err)
}
