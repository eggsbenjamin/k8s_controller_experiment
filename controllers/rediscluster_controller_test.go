// +build unit

package controllers

import (
	"testing"

	dbv1beta1 "github.com/eggsbenjamin/k8s_controller_experiment/api/v1beta1"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestRedisClusterIdentifyAction(t *testing.T) {
	t.Run("add node", func(t *testing.T) {
		redisCluster := &dbv1beta1.RedisCluster{
			Spec: dbv1beta1.RedisClusterSpec{
				Nodes: []dbv1beta1.RedisNodeSpec{
					{
						DiskSize: 1024,
					},
				},
			},
			Status: dbv1beta1.RedisClusterStatus{
				Nodes: []dbv1beta1.RedisNodeStatus{
					// 0 nodes in status
				},
			},
		}

		actionIdentifier := NewRedisClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(redisCluster)
		require.NoError(t, err)
		require.IsType(t, &AddRedisNode{}, action)
	})

	t.Run("remove node", func(t *testing.T) {
		redisCluster := &dbv1beta1.RedisCluster{
			Spec: dbv1beta1.RedisClusterSpec{
				Nodes: []dbv1beta1.RedisNodeSpec{
					// zero nodes in spec
				},
			},
			Status: dbv1beta1.RedisClusterStatus{
				Nodes: []dbv1beta1.RedisNodeStatus{
					{
						IP:       "1.2.3.4",
						DiskSize: 1024,
					},
				},
			},
		}

		actionIdentifier := NewRedisClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(redisCluster)
		require.NoError(t, err)
		require.IsType(t, &RemoveRedisNode{}, action)
	})

	t.Run("update disk size", func(t *testing.T) {
		redisCluster := &dbv1beta1.RedisCluster{
			Spec: dbv1beta1.RedisClusterSpec{
				Nodes: []dbv1beta1.RedisNodeSpec{
					{
						DiskSize: 1024,
					},
					{
						DiskSize: 512,
					},
				},
			},
			Status: dbv1beta1.RedisClusterStatus{
				Nodes: []dbv1beta1.RedisNodeStatus{
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

		actionIdentifier := NewRedisClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(redisCluster)
		require.NoError(t, err)
		require.IsType(t, &UpdateRedisNodeDiskSize{}, action)
		require.Equal(t, 1, action.(*UpdateRedisNodeDiskSize).nodeIndex)
	})

	t.Run("no action to take", func(t *testing.T) {
		redisCluster := &dbv1beta1.RedisCluster{
			Spec: dbv1beta1.RedisClusterSpec{
				Nodes: []dbv1beta1.RedisNodeSpec{
					{
						DiskSize: 1024,
					},
				},
			},
			Status: dbv1beta1.RedisClusterStatus{
				Nodes: []dbv1beta1.RedisNodeStatus{
					{
						IP:       "1.2.3.4",
						DiskSize: 1024,
					},
				},
			},
		}

		actionIdentifier := NewRedisClusterActionIdentifier(nil, zap.Logger(true))

		action, err := actionIdentifier.IdentifyAction(redisCluster)
		require.NoError(t, err)
		require.Nil(t, action)
	})
}
