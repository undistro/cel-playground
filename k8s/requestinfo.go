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

import "gopkg.in/yaml.v3"

// gvkType := apiservercel.NewObjectType("kubernetes.GroupVersionKind", fields(
//
//	field("group", apiservercel.StringType, true),
//	field("version", apiservercel.StringType, true),
//	field("kind", apiservercel.StringType, true),
//
// ))
type GVKType struct {
	Group   string `yaml:"group"`
	Version string `yaml:"version"`
	Kind    string `yaml:"kind"`
}

// gvrType := apiservercel.NewObjectType("kubernetes.GroupVersionResource", fields(
//
//	field("group", apiservercel.StringType, true),
//	field("version", apiservercel.StringType, true),
//	field("resource", apiservercel.StringType, true),
//
// ))
type GVRType struct {
	Group    string `yaml:"group"`
	Version  string `yaml:"version"`
	Resource string `yaml:"resource"`
}

// userInfoType := apiservercel.NewObjectType("kubernetes.UserInfo", fields(
//
//	field("username", apiservercel.StringType, false),
//	field("uid", apiservercel.StringType, false),
//	field("groups", apiservercel.NewListType(apiservercel.StringType, -1), false),
//	field("extra", apiservercel.NewMapType(apiservercel.StringType, apiservercel.NewListType(apiservercel.StringType, -1), -1), false),
//
// ))
type UserInfo struct {
	Username string              `yaml:"username,omitempty"`
	UID      string              `yaml:"uid,omitempty"`
	Groups   []string            `yaml:"groups,omitempty"`
	Extra    map[string][]string `yaml:"extra,omitempty"`
}

// return apiservercel.NewObjectType("kubernetes.AdmissionRequest", fields(
//
//	field("kind", gvkType, true),
//	field("resource", gvrType, true),
//	field("subResource", apiservercel.StringType, false),
//	field("requestKind", gvkType, true),
//	field("requestResource", gvrType, true),
//	field("requestSubResource", apiservercel.StringType, false),
//	field("name", apiservercel.StringType, true),
//	field("namespace", apiservercel.StringType, false),
//	field("operation", apiservercel.StringType, true),
//	field("userInfo", userInfoType, true),
//	field("dryRun", apiservercel.BoolType, false),
//	field("options", apiservercel.DynType, false),
//
// ))
type AdmissionRequest struct {
	Kind               GVKType  `yaml:"kind"`
	Resource           GVRType  `yaml:"resource"`
	SubResource        string   `yaml:"subResource,omitempty"`
	RequestKind        *GVKType `yaml:"requestKind,omitempty"`
	RequestResource    *GVRType `yaml:"requestResource"`
	RequestSubResource string   `yaml:"requestSubResource,omitempty"`
	Name               string   `yaml:"name"`
	Namespace          string   `yaml:"namespace,omitempty"`
	Operation          string   `yaml:"operation"`
	UserInfo           UserInfo `yaml:"userInfo"`
	DryRun             *bool    `yaml:"dryRun,omitempty"`
	// TODO patch/create options
	// Options            any      `yaml:"options,omitempty"`
}

func deserializeRequest(requestData []byte) (map[string]any, error) {
	if requestData == nil {
		return nil, nil
	}
	admissionRequest := AdmissionRequest{}
	if err := yaml.Unmarshal(requestData, &admissionRequest); err != nil {
		return nil, err
	} else {
		return convertToMap(&admissionRequest)
	}
}
