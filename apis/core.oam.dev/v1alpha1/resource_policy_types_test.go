/*
Copyright 2025 The KubeVela Authors.

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

package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/oam-dev/kubevela/pkg/oam"
)

func TestResourcePolicyRuleSelector_Match_ClusterNames(t *testing.T) {
	testCases := map[string]struct {
		selector ResourcePolicyRuleSelector
		input    *unstructured.Unstructured
		matched  bool
	}{
		"cluster name match": {
			selector: ResourcePolicyRuleSelector{
				ClusterNames:  []string{"cluster-prod"},
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						oam.LabelAppCluster: "cluster-prod",
					},
				},
			}},
			matched: true,
		},
		"cluster name mismatch": {
			selector: ResourcePolicyRuleSelector{
				ClusterNames:  []string{"cluster-prod"},
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						oam.LabelAppCluster: "cluster-staging",
					},
				},
			}},
			matched: false,
		},
		"cluster name match with multiple clusters": {
			selector: ResourcePolicyRuleSelector{
				ClusterNames:  []string{"cluster-prod", "cluster-staging"},
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						oam.LabelAppCluster: "cluster-staging",
					},
				},
			}},
			matched: true,
		},
		"empty cluster names selector matches any cluster": {
			selector: ResourcePolicyRuleSelector{
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						oam.LabelAppCluster: "any-cluster",
					},
				},
			}},
			matched: true,
		},
		"cluster name specified but resource has no cluster label": {
			selector: ResourcePolicyRuleSelector{
				ClusterNames:  []string{"cluster-prod"},
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{},
				},
			}},
			matched: false,
		},
		"cluster and component name both match": {
			selector: ResourcePolicyRuleSelector{
				ClusterNames:  []string{"cluster-prod"},
				CompNames:     []string{"my-component"},
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						oam.LabelAppCluster:   "cluster-prod",
						oam.LabelAppComponent: "my-component",
					},
				},
			}},
			matched: true,
		},
		"cluster match but component name mismatch": {
			selector: ResourcePolicyRuleSelector{
				ClusterNames:  []string{"cluster-prod"},
				CompNames:     []string{"my-component"},
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						oam.LabelAppCluster:   "cluster-prod",
						oam.LabelAppComponent: "other-component",
					},
				},
			}},
			matched: false,
		},
		"cluster match but resource type mismatch": {
			selector: ResourcePolicyRuleSelector{
				ClusterNames:  []string{"cluster-prod"},
				ResourceTypes: []string{"Deployment"},
			},
			input: &unstructured.Unstructured{Object: map[string]interface{}{
				"kind": "Service",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						oam.LabelAppCluster: "cluster-prod",
					},
				},
			}},
			matched: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := require.New(t)
			matched := tc.selector.Match(tc.input)
			r.Equal(tc.matched, matched, "Expected matched=%v but got matched=%v", tc.matched, matched)
		})
	}
}
