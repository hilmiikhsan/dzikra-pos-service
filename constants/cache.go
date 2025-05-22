package constants

import "time"

const (
	CacheTTL             = 10 * time.Minute
	CacheShippingCostTTL = 30 * time.Minute
	CacheKeyProvinces    = "address:provinces"
	CacheKeyCitys        = "address:cities"
	CacheKeySubDistricts = "address:subdistricts"
	CacheKeyShippingCost = "shipping:cost"
)
