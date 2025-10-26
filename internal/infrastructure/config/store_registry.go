package config

import (
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
)

// StoreRegistry manages available App Store regions
type StoreRegistry struct {
	stores map[string]*valueobject.Store
}

// NewStoreRegistry creates a new store registry with all supported regions
func NewStoreRegistry() *StoreRegistry {
	registry := &StoreRegistry{
		stores: make(map[string]*valueobject.Store),
	}
	registry.initialize()
	return registry
}

// GetStore returns a store by region code
func (r *StoreRegistry) GetStore(region string) (*valueobject.Store, error) {
	store, exists := r.stores[region]
	if !exists {
		return nil, fmt.Errorf("unsupported region: %s", region)
	}
	return store, nil
}

// GetAllRegions returns all supported region codes
func (r *StoreRegistry) GetAllRegions() []string {
	regions := make([]string, 0, len(r.stores))
	for region := range r.stores {
		regions = append(regions, region)
	}
	return regions
}

// initialize populates the registry with all supported stores
func (r *StoreRegistry) initialize() {
	stores := []struct {
		region     string
		storeFront int
		hostPrefix int
	}{
		{"us", 143441, 36},
		{"ru", 143469, 45},
		{"gb", 143444, 71},
		{"ca", 143455, 71},
		{"fr", 143442, 71},
		{"hk", 143463, 71},
		{"br", 143503, 36},
		{"de", 143443, 36},
		{"jp", 143462, 36},
		{"id", 143476, 28},
		{"kr", 143466, 55},
		{"au", 143460, 55},
		{"in", 143467, 12},
		{"it", 143450, 12},
		{"my", 143473, 55},
		{"mx", 143468, 36},
		{"nl", 143452, 38},
		{"nz", 143461, 42},
		{"sg", 143464, 42},
		{"es", 143454, 40},
		{"za", 143472, 50},
		{"tw", 143470, 70},
		{"th", 143475, 36},
		{"ae", 143481, 36},
		{"vn", 143471, 18},
		{"cn", 143465, 33},
		{"pt", 143453, 39},
		{"tr", 143480, 39},
		{"ar", 143505, 11},
	}

	for _, s := range stores {
		store, _ := valueobject.NewStore(s.region, s.storeFront, s.hostPrefix)
		r.stores[s.region] = store
	}
}
