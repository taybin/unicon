package unicon

import (
	"regexp"
	"strings"
)

var nsSep = regexp.MustCompile("[-_:]")

func namespaceKey(key string, namespaces []string) string {
	namespaced := nsSep.ReplaceAllString(key, ".")

	for _, ns := range namespaces {
		nsWith := strings.Join([]string{ns, "."}, "")
		if strings.HasPrefix(strings.ToLower(namespaced), nsWith) {
			return namespaced
		}
	}
	return key
}

func nsSlice(namespaces []string) (lowered []string) {
	for _, ns := range namespaces {
		// put in lowercase
		lowered = append(lowered, strings.ToLower(ns))
	}

	return
}
