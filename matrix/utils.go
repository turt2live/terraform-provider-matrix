package matrix

import (
	"strings"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func nilIfEmptyString(val interface{}) interface{} {
	if val == "" {
		return nil
	}
	return val
}

func stripMxc(input string) (string, string, string, error) {
	if !strings.HasPrefix(input, "mxc://") {
		return "", "", "", fmt.Errorf("invalid mxc: missing protocol")
	}

	withoutProto := strings.TrimSpace(input[len("mxc://"):])
	withoutProto = strings.Split(withoutProto, "?")[0] // Strip query string
	withoutProto = strings.Split(withoutProto, "#")[0] // Strip fragment
	if len(withoutProto) == 0 {
		return "", "", "", fmt.Errorf("invalid mxc: no origin or media_id")
	}

	parts := strings.Split(withoutProto, "/")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid mxc: wrong number of segments. expected: %d  got: %d", 2, len(parts))
	}

	origin := parts[0]
	mediaId := parts[1]
	constructed := fmt.Sprintf("mxc://%s/%s", origin, mediaId)
	return constructed, origin, mediaId, nil
}

func setOfStrings(val *schema.Set) []string {
	res := make([]string, 0)

	if val != nil {
		for _, v := range val.List() {
			res = append(res, v.(string))
		}
	}

	return res
}
