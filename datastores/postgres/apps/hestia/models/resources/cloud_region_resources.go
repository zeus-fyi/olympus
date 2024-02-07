package hestia_compute_resources

import hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"

// Resources represents a collection of nodes.
type Resources struct {
	Nodes hestia_autogen_bases.NodesSlice `json:"nodes"`
}

// RegionResourcesMap maps region names to their corresponding Resources.
type RegionResourcesMap map[string]Resources

// CloudProviderRegionsResourcesMap maps cloud provider names to their RegionResourcesMap,
// allowing for a nested mapping of providers to regions to resources.
type CloudProviderRegionsResourcesMap map[string]RegionResourcesMap
