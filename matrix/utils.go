package matrix

func nilIfEmptyString(val interface{}) interface{} {
	if val == "" {
		return nil
	}
	return val
}
