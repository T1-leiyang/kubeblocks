/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CollectorDataSourceSpec defines the desired state of CollectorDataSource
type OTeldCollectorTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image is the image of the oteld
	Image string `json:"image,omitempty"`

	// LogsLevel is the level of the oteld collector logs
	LogsLevel string `json:"logLevel,omitempty"`

	// MetricsPort is the port of the oteld metrics
	MetricsPort int `json:"metricsPort,omitempty"`

	// Resources is the resource requirements for the oteld
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Mode represents how the OTeld should be deployed (deployment, daemonset, statefulset or sidecar)
	Mode Mode `json:"mode,omitempty"`

	// UseConfigMap indicates whether to use configmap to store oteld config
	UseConfigMap bool `json:"useConfigMap"`

	// NodeSelector to schedule OpenTelemetry Collector pods.
	// This is only relevant to daemonset, statefulset, and deployment mode
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// ENV vars to set on the OpenTelemetry Collector's Pods. These can then in certain cases be
	// consumed in the config file for the Collector.
	// +optional
	Env []v1.EnvVar `json:"env,omitempty"`

	// PodAnnotations is the set of annotations that will be attached to
	// Collector and Target Allocator pods.
	// +optional
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`

	// Volumes represents which volumes to use in the underlying collector deployment(s).
	// +optional
	// +listType=atomic
	Volumes []v1.Volume `json:"volumes,omitempty"`

	// VolumeMounts represents the mount points to use in the underlying collector deployment(s)
	// +optional
	// +listType=atomic
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`

	// Ports allows a set of ports to be exposed by the underlying v1.Service. By default, the operator
	// will attempt to infer the required ports by parsing the .Spec.Config property but this property can be
	// used to open additional ports that can't be inferred by the operator, like for custom receivers.
	// +optional
	// +listType=atomic
	Ports []v1.ServicePort `json:"ports,omitempty"`

	// SecurityContext configures the container security context for
	// the opentelemetry-collector container.
	//
	// In deployment, daemonset, or statefulset mode, this controls
	// the security context settings for the primary application
	// container.
	//
	// In sidecar mode, this controls the security context for the
	// injected sidecar container.
	//
	// +optional
	SecurityContext v1.SecurityContext `json:"securityContext,omitempty"`

	// PodSecurityContext configures the pod security context for the
	// opentelemetry-collector pod, when running as a deployment, daemonset,
	// or statefulset.
	//
	// In sidecar mode, the opentelemetry-operator will ignore this setting.
	//
	// +optional
	PodSecurityContext v1.PodSecurityContext `json:"podSecurityContext,omitempty"`
}

// CollectorDataSourceStatus defines the observed state of CollectorDataSource
type OTeldCollectorTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CollectorDataSource is the Schema for the collectordatasources API
type OTeldCollectorTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OTeldCollectorTemplateSpec   `json:"spec,omitempty"`
	Status OTeldCollectorTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CollectorDataSourceList contains a list of CollectorDataSource
type OTeldCollectorTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OTeldCollectorTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OTeldCollectorTemplate{}, &OTeldCollectorTemplateList{})
}