package tencent

import (
	"regexp"
)

type Domain struct {
	domain    string
	subDomain string
}

func getDomain(str string) *Domain {
	domainRegexp := regexp.MustCompile(`(.*)\.(.*\..*)\.`)
	params := domainRegexp.FindStringSubmatch(str)
	return &Domain{domain: params[len(params)-1], subDomain: params[len(params)-2]}
}
