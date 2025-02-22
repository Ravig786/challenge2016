package models

import "sync"

type DistributorPermissions struct {
	Includes map[string]struct{}
	Excludes map[string]struct{}
	Parent   string
}

type DistributorRegistry struct {
	sync.RWMutex
	Distributors map[string]*DistributorPermissions
}

var Registry *DistributorRegistry

func InitDistributorRegistry() {
	Registry = &DistributorRegistry{
		Distributors: make(map[string]*DistributorPermissions),
	}
}
