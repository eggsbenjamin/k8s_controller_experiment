// +build unit

package controllers

import (
	"testing"

	dbv1beta1 "github.com/eggsbenjamin/k8s_controller_experiment/api/v1beta1"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestCassandraClusterIdentifyAction(t *testing.T) {
	t.Run("add node", func(t *testing.T) {
		cassandraCluster := &dbv1beta1.CassandraCluster{
			Spec: dbv1beta1.CassandraClusterSpec{
				Nodes: []dbv1beta1.CassandraNodeSpec{
					{
						DiskSize: 1024,
					},
				},
			},
			Status: dbv1beta1.CassandraClusterStatus{
				Nodes: []dbv1beta1.CassandraNodeStatus{
					// 0 nodes in status
				},
			},
		}

		actionIdentifier := NewCassandraClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(cassandraCluster)
		require.NoError(t, err)
		require.IsType(t, &AddCassandraNode{}, action)
	})

	t.Run("remove node", func(t *testing.T) {
		cassandraCluster := &dbv1beta1.CassandraCluster{
			Spec: dbv1beta1.CassandraClusterSpec{
				Nodes: []dbv1beta1.CassandraNodeSpec{
					// zero nodes in spec
				},
			},
			Status: dbv1beta1.CassandraClusterStatus{
				Nodes: []dbv1beta1.CassandraNodeStatus{
					{
						IP:       "1.2.3.4",
						DiskSize: 1024,
					},
				},
			},
		}

		actionIdentifier := NewCassandraClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(cassandraCluster)
		require.NoError(t, err)
		require.IsType(t, &RemoveCassandraNode{}, action)
	})

	t.Run("update disk size", func(t *testing.T) {
		cassandraCluster := &dbv1beta1.CassandraCluster{
			Spec: dbv1beta1.CassandraClusterSpec{
				Nodes: []dbv1beta1.CassandraNodeSpec{
					{
						DiskSize: 1024,
					},
					{
						DiskSize: 512,
					},
				},
			},
			Status: dbv1beta1.CassandraClusterStatus{
				Nodes: []dbv1beta1.CassandraNodeStatus{
					{
						IP:       "1.2.3.4",
						DiskSize: 1024,
					},
					{
						IP:       "1.2.3.5",
						DiskSize: 1024,
					},
				},
			},
		}

		actionIdentifier := NewCassandraClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(cassandraCluster)
		require.NoError(t, err)
		require.IsType(t, &UpdateCassandraNodeDiskSize{}, action)
		require.Equal(t, 1, action.(*UpdateCassandraNodeDiskSize).nodeIndex)
	})

	t.Run("no action to take", func(t *testing.T) {
		cassandraCluster := &dbv1beta1.CassandraCluster{
			Spec: dbv1beta1.CassandraClusterSpec{
				Nodes: []dbv1beta1.CassandraNodeSpec{
					{
						DiskSize: 1024,
					},
				},
			},
			Status: dbv1beta1.CassandraClusterStatus{
				Nodes: []dbv1beta1.CassandraNodeStatus{
					{
						IP:       "1.2.3.4",
						DiskSize: 1024,
					},
				},
			},
		}

		actionIdentifier := NewCassandraClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(cassandraCluster)
		require.NoError(t, err)
		require.Nil(t, action)
	})
}
