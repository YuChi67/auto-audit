package common

type (
	NotifyEnum       uint32
)

const (
	BUG        NotifyEnum = iota + 1
	VULNERABILITY
	COVERAGE
)

var Notify = map[string]NotifyEnum{
	"bug":             BUG,
	"vulnerability":   VULNERABILITY,
	"coverage":        COVERAGE,
}


var NotifyStr = map[NotifyEnum]string{
	BUG: 		    "bug",
	VULNERABILITY: "vulnerability",
	COVERAGE: 	   "coverage",
}