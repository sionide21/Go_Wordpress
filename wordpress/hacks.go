package wordpress

import (
	"xmlrpc"
	"strconv"
)

func getInt(i xmlrpc.StructValue, key string) (ret int) {
	ret, _ = strconv.Atoi(i.GetString(key))
	return
}
