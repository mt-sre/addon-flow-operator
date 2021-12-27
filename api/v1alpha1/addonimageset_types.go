/*
Copyright 2021.

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
	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddonImageSetSpec defines the desired state of AddonImageSet
type AddonImageSetSpec struct {
	// +kubebuilder:validation:Required
	// The name of the imageset along with the version.
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^quay\.io/osd-addons/[a-z-]+`
	// The url for the index image
	IndexImage string `json:"indexImage"`

	// +kubebuilder:validation:Required
	// A list of image urls of related operators
	RelatedImages []string `json:"relatedImages"`

	// +optional
	// OCM representation of an add-on parameter
	AddOnParameters *[]ocmv1.AddOnParameter `json:"addOnParameters"`

	// +optional
	// OCM representation of an addon-requirement
	AddOnRequirements *[]ocmv1.AddOnRequirement `json:"addOnRequirements"`

	// +optional
	// OCM representation of an add-on sub operator. A sub operator is an
	// operator who's life cycle is controlled by the add-on umbrella operator.
	SubOperators *[]ocmv1.AddOnSubOperator `json:"subOperators"`
}

// AddonImageSetStatus defines the observed state of AddonImageSet
type AddonImageSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// AddonImageSet is the Schema for the addonimagesets API
type AddonImageSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AddonImageSetSpec   `json:"spec,omitempty"`
	Status AddonImageSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// AddonImageSetList contains a list of AddonImageSet
type AddonImageSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AddonImageSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AddonImageSet{}, &AddonImageSetList{})
}
