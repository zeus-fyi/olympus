package hestia_gcp

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

type GcpClusterInfo struct {
	ClusterName string
	ProjectID   string
	Zone        string
}

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

//func (g *GcpClient) ListNodeTypes(ctx context.Context, client *container.Service, projectID, zone, clusterID string) error {
//	c, err := compute.NewNodeTypesRESTClient(ctx)
//	if err != nil {
//		log.Ctx(ctx).Err(err).Msg("Failed to create compute client")
//		return err
//	}
//	return nil
//}

func (g *GcpClient) ListNodes(ctx context.Context, ci GcpClusterInfo) ([]*container.NodePool, error) {
	nodePools, err := g.Projects.Zones.Clusters.NodePools.List(ci.ProjectID, ci.Zone, ci.ClusterName).Context(ctx).Do()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to retrieve node pools")
		return nil, err
	}
	return nodePools.NodePools, err
}
