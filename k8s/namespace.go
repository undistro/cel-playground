// Copyright 2024 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"time"

	"gopkg.in/yaml.v3"
)

// specType := apiservercel.NewObjectType("kubernetes.NamespaceSpec", fields(
//
//	field("finalizers", apiservercel.NewListType(apiservercel.StringType, -1), true),
//
// ))
type NamespaceSpecType struct {
	Finalizers []string `yaml:"finalizers"`
}

// conditionType := apiservercel.NewObjectType("kubernetes.NamespaceCondition", fields(
//
//	field("status", apiservercel.StringType, true),
//	field("type", apiservercel.StringType, true),
//	field("lastTransitionTime", apiservercel.TimestampType, true),
//	field("message", apiservercel.StringType, true),
//	field("reason", apiservercel.StringType, true),
//
// ))
type NamespaceConditionType struct {
	Status             string    `yaml:"status"`
	Type               string    `yaml:"type"`
	LastTransitionTime time.Time `yaml:"lastTransitionTime"`
	Message            string    `yaml:"message"`
	Reason             string    `yaml:"reason"`
}

// statusType := apiservercel.NewObjectType("kubernetes.NamespaceStatus", fields(
//
//	field("conditions", apiservercel.NewListType(conditionType, -1), true),
//	field("phase", apiservercel.StringType, true),
//
// ))
type NamespaceStatusType struct {
	Conditions []NamespaceConditionType `yaml:"conditions"`
	Phase      string                   `yaml:"phase"`
}

// metadataType := apiservercel.NewObjectType("kubernetes.NamespaceMetadata", fields(
//
//	field("name", apiservercel.StringType, true),
//	field("generateName", apiservercel.StringType, true),
//	field("namespace", apiservercel.StringType, true),
//	field("labels", apiservercel.NewMapType(apiservercel.StringType, apiservercel.StringType, -1), true),
//	field("annotations", apiservercel.NewMapType(apiservercel.StringType, apiservercel.StringType, -1), true),
//	field("UID", apiservercel.StringType, true),
//	field("creationTimestamp", apiservercel.TimestampType, true),
//	field("deletionGracePeriodSeconds", apiservercel.IntType, true),
//	field("deletionTimestamp", apiservercel.TimestampType, true),
//	field("generation", apiservercel.IntType, true),
//	field("resourceVersion", apiservercel.StringType, true),
//	field("finalizers", apiservercel.NewListType(apiservercel.StringType, -1), true),
//
// ))
type NamespaceMetadataType struct {
	Name                       string            `yaml:"name"`
	GenerateName               string            `yaml:"generateName"`
	Namespace                  string            `yaml:"namespace"`
	Labels                     map[string]string `yaml:"labels"`
	Annotations                map[string]string `yaml:"annotations"`
	UID                        string            `yaml:"UID"`
	CreationTimestamp          time.Time         `yaml:"creationTimestamp"`
	DeletionGracePeriodSeconds int64             `yaml:"deletionGracePeriodSeconds"`
	DeletionTimestamp          time.Time         `yaml:"deletionTimestamp"`
	Generation                 int64             `yaml:"generation"`
	ResourceVersion            string            `yaml:"resourceVersion"`
	Finalizers                 []string          `yaml:"finalizers"`
}

// return apiservercel.NewObjectType("kubernetes.Namespace", fields(
//
//	field("metadata", metadataType, true),
//	field("spec", specType, true),
//	field("status", statusType, true),
//
// ))
type NamespaceType struct {
	Metadata NamespaceMetadataType `yaml:"metadata"`
	Spec     NamespaceSpecType     `yaml:"spec"`
	Status   NamespaceStatusType   `yaml:"status"`
}

func deserializeNamespace(namespaceData []byte) (map[string]any, error) {
	if namespaceData == nil {
		return nil, nil
	}
	namespace := NamespaceType{}
	if err := yaml.Unmarshal(namespaceData, &namespace); err != nil {
		return nil, err
	} else {
		return convertToMap(&namespace)
	}
}
