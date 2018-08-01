package store

import (
	"strconv"
)

type (
	GeoStore map[string]*ItunesStore

	ItunesStore struct {
		XAppleStoreFront string
		StoreFront       TypeStoreFront
		HostPrefix       int
	}

	TypeStoreFront int
)

const (
	DefaultCategory = "36"

	//POPULARS iPhone
	TopFree     = "27"
	TopPaid     = "30"
	TopGrossing = "38"

	//POPULARS iPad
	IPadTopFree     = "44"
	IPadTopPaid     = "47"
	IPadTopGrossing = "46"

	IPhoneDeviceCode         = 29
	XAppleStoreFrontTemplate = "%d,%d"
)

func (s TypeStoreFront) String() string {
	return strconv.Itoa(int(s))
}

func (s TypeStoreFront) Int() int {
	return int(s)
}

func GetStoreList() GeoStore {
	return GeoStore{
		"ru": &ItunesStore{
			XAppleStoreFront: "143469-16,32",
			StoreFront:       143469,
			HostPrefix:       45,
		},
		"us": &ItunesStore{
			XAppleStoreFront: "143441-1,32",
			StoreFront:       143441,
			HostPrefix:       36,
		},
		"gb": &ItunesStore{
			XAppleStoreFront: "143444,32",
			StoreFront:       143444,
			HostPrefix:       71,
		},
		"ca": &ItunesStore{
			XAppleStoreFront: "143455-6,32",
			StoreFront:       143455,
			HostPrefix:       71,
		},
		"fr": &ItunesStore{
			XAppleStoreFront: "143442-2,32",
			StoreFront:       143442,
			HostPrefix:       71,
		},
		"hk": &ItunesStore{
			XAppleStoreFront: "143463-45,32",
			StoreFront:       143463,
			HostPrefix:       71,
		},
		"br": &ItunesStore{
			XAppleStoreFront: "143503-15,32",
			StoreFront:       143503,
			HostPrefix:       36,
		},
		"de": &ItunesStore{
			XAppleStoreFront: "143443-2,32",
			StoreFront:       143443,
			HostPrefix:       36,
		},
		"jp": &ItunesStore{
			XAppleStoreFront: "143462-2,32",
			StoreFront:       143462,
			HostPrefix:       36,
		},
		"id": &ItunesStore{
			XAppleStoreFront: "143476-2,32",
			StoreFront:       143476,
			HostPrefix:       28,
		},
		"kr": &ItunesStore{
			XAppleStoreFront: "143466-13,32",
			StoreFront:       143466,
			HostPrefix:       55,
		},
		"au": &ItunesStore{
			XAppleStoreFront: "143460,32",
			StoreFront:       143460,
			HostPrefix:       55,
		},
		"in": &ItunesStore{
			XAppleStoreFront: "143467,32",
			StoreFront:       143467,
			HostPrefix:       12,
		},
		"it": &ItunesStore{
			XAppleStoreFront: "143450-7,32",
			StoreFront:       143450,
			HostPrefix:       12,
		},
		"my": &ItunesStore{
			XAppleStoreFront: "143473-2,32",
			StoreFront:       143473,
			HostPrefix:       55,
		},
		"mx": &ItunesStore{
			XAppleStoreFront: "143468-2,32",
			StoreFront:       143468,
			HostPrefix:       36,
		},
		"nl": &ItunesStore{
			XAppleStoreFront: "143452-2,32",
			StoreFront:       143452,
			HostPrefix:       38,
		},
		"nz": &ItunesStore{
			XAppleStoreFront: "143461,32",
			StoreFront:       143461,
			HostPrefix:       42,
		},
		"sg": &ItunesStore{
			XAppleStoreFront: "143464-2,32",
			StoreFront:       143464,
			HostPrefix:       42,
		},
		"es": &ItunesStore{
			XAppleStoreFront: "143454-2,32",
			StoreFront:       143454,
			HostPrefix:       40,
		},
		"za": &ItunesStore{
			XAppleStoreFront: "143472,32",
			StoreFront:       143472,
			HostPrefix:       50,
		},
		"tw": &ItunesStore{
			XAppleStoreFront: "143470-2,32",
			StoreFront:       143470,
			HostPrefix:       70,
		},
		"th": &ItunesStore{
			XAppleStoreFront: "143475-2,32",
			StoreFront:       143475,
			HostPrefix:       36,
		},
		"ae": &ItunesStore{
			XAppleStoreFront: "143481,32",
			StoreFront:       143481,
			HostPrefix:       36,
		},
		"vn": &ItunesStore{
			XAppleStoreFront: "143471-2,32",
			StoreFront:       143471,
			HostPrefix:       18,
		},
		"cn": &ItunesStore{
			XAppleStoreFront: "143465-2,32",
			StoreFront:       143465,
			HostPrefix:       33,
		},
		"pt": &ItunesStore{
			XAppleStoreFront: "143453-24,32",
			StoreFront:       143453,
			HostPrefix:       39,
		},
		"tr": &ItunesStore{
			XAppleStoreFront: "143480-2,32",
			StoreFront:       143480,
			HostPrefix:       39,
		},
		"ar": &ItunesStore{
			XAppleStoreFront: "143505-2,32",
			StoreFront:       143505,
			HostPrefix:       11,
		},
	}
}

func PopIdExists(popId string) bool {
	if popId != TopFree && popId != TopPaid && popId != TopGrossing {
		return false
	}

	return true
}
