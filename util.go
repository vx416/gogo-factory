package factory

func getColName(options []string) string {
	if len(options) == 0 {
		return ""
	}
	return options[0]
}
