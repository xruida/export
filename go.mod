module github.com/xruida/export

require (
	baliance.com/gooxml v1.0.1
	github.com/issue9/assert v1.0.0
	github.com/issue9/is v1.1.1 // indirect
	github.com/issue9/logs v1.0.0
	github.com/issue9/version v1.0.0
	github.com/issue9/web v0.23.0
	gopkg.in/yaml.v2 v2.2.1
)

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180905080454-ebe1bf3edb33
	golang.org/x/text => github.com/golang/text v0.3.0
)
