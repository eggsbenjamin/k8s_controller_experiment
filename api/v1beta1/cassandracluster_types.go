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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CassandraNodeSpec struct {
	DiskSize int `json:"diskSize,omitempty"`
}

type CassandraNodeStatus struct {
	IP       string `json:"ip,omitempty"`
	DiskSize int    `json:"diskSize,omitempty"`
}

// CassandraClusterSpec defines the desired state of CassandraCluster
type CassandraClusterSpec struct {
	Nodes []CassandraNodeSpec `json:"nodes,omitempty"`
}

// CassandraClusterStatus defines the observed state of CassandraCluster
type CassandraClusterStatus struct {
	Nodes []CassandraNodeStatus `json:"nodes,omitempty"`
}

// +kubebuilder:object:root=true

// CassandraCluster is the Schema for the cassandraclusters API
type CassandraCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CassandraClusterSpec   `json:"spec,omitempty"`
	Status CassandraClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CassandraClusterList contains a list of CassandraCluster
type CassandraClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CassandraCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CassandraCluster{}, &CassandraClusterList{})
}
