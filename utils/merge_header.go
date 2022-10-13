package utils

import "net/http"

func MergeHeader(target http.Header, source http.Header) http.Header {
	merged := target
	if merged == nil {
		merged = make(http.Header)
	}
	for key, values := range source {
		if merged.Get(key) == "" {
			for _, value := range values {
				merged.Set(key, value)
			}
		} else {
			for _, value := range values {
				merged.Add(key, value)
			}
		}
	}
	return merged
}
