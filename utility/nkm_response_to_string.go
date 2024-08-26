package utility

import (
	"encoding/json"

	"github.com/samber/lo"
)

func Response(v interface{}) string {
	return string(lo.Must1(json.Marshal(v)))
}
