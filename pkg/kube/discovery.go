package kube

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

// APIAvailability tracks which optional APIs are available in the cluster
type APIAvailability struct {
	RoutesAvailable        bool
	IngressRoutesAvailable bool
	HTTPRoutesAvailable    bool
}

// DiscoverAPIs checks which optional APIs are available in the cluster
func DiscoverAPIs(config *rest.Config) APIAvailability {
	client, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		logger.Warnf("Failed to create discovery client: %v", err)
		return APIAvailability{}
	}

	return APIAvailability{
		RoutesAvailable:        apiGroupExists(client, "route.openshift.io"),
		IngressRoutesAvailable: apiGroupExists(client, "traefik.io") || apiGroupExists(client, "traefik.containo.us"),
		HTTPRoutesAvailable:    apiGroupExists(client, "gateway.networking.k8s.io"),
	}
}

func apiGroupExists(client discovery.DiscoveryInterface, groupName string) bool {
	groups, err := client.ServerGroups()
	if err != nil {
		return false
	}

	for _, group := range groups.Groups {
		if group.Name == groupName {
			return true
		}
	}
	return false
}
