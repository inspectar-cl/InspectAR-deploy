package pkg

func JoinPath(a, b string) string {
	switch {
	case a == "" || a == "/":
		if b == "" {
			return "/"
		}
		if b[0] != '/' {
			return "/" + b
		}
		return b
	default:
		if a[len(a)-1] == '/' && b != "" && b[0] == '/' {
			return a + b[1:]
		}
		if a[len(a)-1] != '/' && (b == "" || b[0] != '/') {
			return a + "/" + b
		}
		return a + b
	}
}
