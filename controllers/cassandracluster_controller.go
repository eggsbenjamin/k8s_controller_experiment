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

// CassandraClusterReconciler reconciles a CassandraCluster object
type CassandraClusterReconciler struct {
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

// +kubebuilder:rbac:groups=db.k8s.io,resources=cassandraclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=db.k8s.io,resources=cassandraclusters/status,verbs=get;update;patch
func (r *CassandraClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("cassandracluster", req.NamespacedName)

	cassandraCluster := &dbv1beta1.CassandraCluster{}
	if err := r.Get(context.TODO(), types.NamespacedName{}, cassandraCluster); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil // deleted
		}

		return ctrl.Result{}, err
	}

	action, err := r.actionIdentifier.IdentifyAction(cassandraCluster)
	if err != nil {
		return ctrl.Result{}, err
	}

	if action != nil {
		return ctrl.Result{}, action.Execute()
	}

	return ctrl.Result{}, nil // no action to take
}

func (r *CassandraClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1beta1.CassandraCluster{}).
		Complete(r)
}

type AddCassandraNode struct {
	cassandraCluster *dbv1beta1.CassandraCluster
	k8sClient        client.Client
	log              logr.Logger
}

func (a *AddCassandraNode) Execute() error {
	return nil
}

type RemoveCassandraNode struct {
	cassandraCluster *dbv1beta1.CassandraCluster
	k8sClient        client.Client
	log              logr.Logger
}

func (a *RemoveCassandraNode) Execute() error {
	return nil
}

type UpdateCassandraNodeDiskSize struct {
	cassandraCluster *dbv1beta1.CassandraCluster
	nodeIndex        int
	k8sClient        client.Client
	log              logr.Logger
}

func (a *UpdateCassandraNodeDiskSize) Execute() error {
	return nil
}

type CassandraClusterActionIdentifier struct {
	k8sClient client.Client
	log       logr.Logger
}

func NewCassandraClusterActionIdentifier(k8sClient client.Client, log logr.Logger) ActionIdentifier {
	return &CassandraClusterActionIdentifier{
		k8sClient: k8sClient,
		log:       log,
	}
}

//  IdentifyAction inspects a CassandraCluster resource to determine the delta of highest priority and returns an identifier for an appropriate action
func (c *CassandraClusterActionIdentifier) IdentifyAction(obj runtime.Object) (Action, error) {
	cassandraCluster, ok := obj.(*dbv1beta1.CassandraCluster)
	if !ok {
		return nil, fmt.Errorf("unexpected runtime object: %#v", obj)
	}

	if len(cassandraCluster.Spec.Nodes) > len(cassandraCluster.Status.Nodes) {
		return &AddCassandraNode{
			cassandraCluster: cassandraCluster,
			k8sClient:        c.k8sClient,
			log:              c.log,
		}, nil
	}

	if len(cassandraCluster.Spec.Nodes) < len(cassandraCluster.Status.Nodes) {
		return &RemoveCassandraNode{
			cassandraCluster: cassandraCluster,
			k8sClient:        c.k8sClient,
			log:              c.log,
		}, nil
	}

	for i := 0; i < len(cassandraCluster.Spec.Nodes); i++ {
		if cassandraCluster.Spec.Nodes[i].DiskSize != cassandraCluster.Status.Nodes[i].DiskSize {
			return &UpdateCassandraNodeDiskSize{
				cassandraCluster: cassandraCluster,
				nodeIndex:        i,
				k8sClient:        c.k8sClient,
				log:              c.log,
			}, nil
		}
	}

	return nil, nil
}
