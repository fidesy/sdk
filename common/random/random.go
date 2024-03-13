package randomCommon

import (
	"fmt"
	"math/rand"
)

func RandomPort() string {
	return fmt.Sprintf("%d", rand.Intn(10000)+50000)
}
