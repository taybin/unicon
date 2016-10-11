package unicon

import (
	"regexp"
	"strconv"
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

func unmarshal(segment interface{}, path string, output map[string]interface{}) {
	switch segment := segment.(type) {
	case map[string]interface{}:
		path += "."
		unmarshalMap(segment, path, output)
	case []interface{}:
		unmarshalArray(segment, path, output)
	default:
		output[path] = segment
	}
}

func unmarshalMap(segment map[string]interface{}, segmentPath string, output map[string]interface{}) {
	for k, v := range segment {
		keyWithPath := segmentPath + k
		unmarshal(v, keyWithPath, output)
	}
}

func unmarshalArray(segment []interface{}, segmentPath string, output map[string]interface{}) {
	for i, v := range segment {
		keyWithPath := segmentPath + "[" + strconv.Itoa(i) + "]"
		unmarshal(v, keyWithPath, output)
	}
	output[segmentPath+".length"] = len(segment)
}
