/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dbv1beta1 "github.com/eggsbenjamin/k8s_controller_experiment/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// RedisClusterReconciler reconciles a RedisCluster object
type RedisClusterReconciler struct {
	client.Client
	Log              logr.Logger
	actionIdentifier ActionIdentifier
}

/*
	2 steps
	- inspect resource to identify delta of highest priority
	- return meaningful identifier for this action corresponding to the delta
	- look up appropriate action, handler, using this identifier
*/

// +kubebuilder:rbac:groups=db.k8s.io,resources=redisclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=db.k8s.io,resources=redisclusters/status,verbs=get;update;patch
func (r *RedisClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("rediscluster", req.NamespacedName)

	redisCluster := &dbv1beta1.RedisCluster{}
	if err := r.Get(context.TODO(), types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, redisCluster); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil // deleted
		}

		return ctrl.Result{}, err
	}

	action, err := r.actionIdentifier.IdentifyAction(redisCluster)
	if err != nil {
		return ctrl.Result{}, err
	}

	if action != nil {
		return ctrl.Result{}, action.Execute()
	}

	return ctrl.Result{}, nil // no action to take
}

func (r *RedisClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1beta1.RedisCluster{}).
		Complete(r)
}

type AddRedisNode struct {
	redisCluster *dbv1beta1.RedisCluster
	k8sClient    client.Client
	log          logr.Logger
}

func (a *AddRedisNode) Execute() error {
	return nil
}

type RemoveRedisNode struct {
	redisCluster *dbv1beta1.RedisCluster
	k8sClient    client.Client
	log          logr.Logger
}

func (a *RemoveRedisNode) Execute() error {
	return nil
}

type UpdateRedisNodeDiskSize struct {
	redisCluster *dbv1beta1.RedisCluster
	nodeIndex    int
	k8sClient    client.Client
	log          logr.Logger
}

func (a *UpdateRedisNodeDiskSize) Execute() error {
	return nil
}

type RedisClusterActionIdentifier struct {
	k8sClient client.Client
	log       logr.Logger
}

func NewRedisClusterActionIdentifier(k8sClient client.Client, log logr.Logger) ActionIdentifier {
	return &RedisClusterActionIdentifier{
		k8sClient: k8sClient,
		log:       log,
	}
}

//  IdentifyAction inspects a RedisCluster resource to determine the delta of highest priority and returns an identifier for an appropriate action
func (c *RedisClusterActionIdentifier) IdentifyAction(obj runtime.Object) (Action, error) {
	redisCluster, ok := obj.(*dbv1beta1.RedisCluster)
	if !ok {
		return nil, fmt.Errorf("unexpected runtime object: %#v", obj)
	}

	if len(redisCluster.Spec.Nodes) > len(redisCluster.Status.Nodes) {
		return &AddRedisNode{
			redisCluster: redisCluster,
			k8sClient:    c.k8sClient,
			log:          c.log,
		}, nil
	}

	if len(redisCluster.Spec.Nodes) < len(redisCluster.Status.Nodes) {
		return &RemoveRedisNode{
			redisCluster: redisCluster,
			k8sClient:    c.k8sClient,
			log:          c.log,
		}, nil
	}

	for i := 0; i < len(redisCluster.Spec.Nodes); i++ {
		if redisCluster.Spec.Nodes[i].DiskSize != redisCluster.Status.Nodes[i].DiskSize {
			return &UpdateRedisNodeDiskSize{
				redisCluster: redisCluster,
				nodeIndex:    i,
				k8sClient:    c.k8sClient,
				log:          c.log,
			}, nil
		}
	}

	return nil, nil
}
