package bdns

var caaServfailException map[string]bool

func init() {
	const exceptions = `
foo.com
`
	for _, v := range strings.Split(exceptions, "\n") {
		if len(v) > 0 {
			caaServfailException[v] = true
		}
	}
}
