package iptable

import (
	"ava/core"
)

func Allow(dst string) bool {
	if !stringInSlice(dst, core.Sites) {
		return false
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
