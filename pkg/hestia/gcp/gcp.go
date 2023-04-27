package hestia_gcp

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

type GcpClient struct {
	*container.Service
}

func InitGcpClient(ctx context.Context, authJsonBytes []byte) (GcpClient, error) {
	client, err := container.NewService(ctx, option.WithCredentialsJSON(authJsonBytes), option.WithScopes(container.CloudPlatformScope))
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to create GKE API client")
		return GcpClient{}, err
	}
	return GcpClient{client}, nil
}

type GcpClusterInfo struct {
	ClusterName string
	ProjectID   string
	Zone        string
}

func (g *GcpClient) GetNodeSizes(ctx context.Context, ci GcpClusterInfo) ([]string, error) {
	nodePools, err := g.Projects.Zones.Clusters.NodePools.List(ci.ProjectID, ci.Zone, ci.ClusterName).Context(ctx).Do()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to retrieve node pools")
		return nil, err
	}

	// Print the name and status of each node pool
	for _, np := range nodePools.NodePools {
		fmt.Printf("%s: %s\n", np.Name, np.Status)
	}
	return nil, err
}
