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
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

var (
	AuthorizerType    = cel.OpaqueType("playground.k8s.Authorizer")
	PathCheckType     = cel.OpaqueType("playground.k8s.PathCheck")
	GroupCheckType    = cel.OpaqueType("playground.k8s.GroupCheck")
	ResourceCheckType = cel.OpaqueType("playground.k8s.ResourceCheck")
	DecisionType      = cel.OpaqueType("playground.k8s.Decision")
)

var _ traits.Receiver = &Authorizer{}

type Authorizer struct {
	receiverOnlyObjectVal
	Paths           map[string]*PathCheck             `yaml:"paths,omitempty"`
	Groups          map[string]*GroupCheck            `yaml:"groups,omitempty"`
	ServiceAccounts map[string]map[string]*Authorizer `yaml:"serviceAccounts,omitempty"`
}

func (a *Authorizer) Receive(function string, overload string, args []ref.Val) ref.Val {
	switch len(args) {
	case 1:
		switch function {
		case "path":
			if path, ok := getString(args[0].Value()); ok {
				if len(path) == 0 {
					return types.NewErr("path must not be empty")
				} else if a.Paths != nil {
					if pathCheck, ok := a.Paths[path]; ok {
						initReceiver(&pathCheck.receiverOnlyObjectVal, PathCheckType)
						return pathCheck
					}
				}
				pathCheck := &PathCheck{}
				initReceiver(&pathCheck.receiverOnlyObjectVal, PathCheckType)
				return pathCheck
			}
		case "group":
			if group, ok := getString(args[0].Value()); ok {
				if a.Groups != nil {
					if groupCheck, ok := a.Groups[group]; ok {
						initReceiver(&groupCheck.receiverOnlyObjectVal, GroupCheckType)
						return groupCheck
					}
				}
				groupCheck := &GroupCheck{}
				initReceiver(&groupCheck.receiverOnlyObjectVal, GroupCheckType)
				return groupCheck
			}
		}
	case 2:
		switch function {
		case "serviceAccount":
			// TODO check the namespace and name to see if they are valid
			if namespace, ok := getString(args[0].Value()); ok {
				if name, ok := getString(args[1].Value()); ok {
					if a.ServiceAccounts != nil {
						if namespacedServiceAccounts, ok := a.ServiceAccounts[namespace]; ok {
							if authorizer, ok := namespacedServiceAccounts[name]; ok {
								initReceiver(&authorizer.receiverOnlyObjectVal, AuthorizerType)
								return authorizer
							}
						}
					}
					authorizer := &Authorizer{}
					initReceiver(&authorizer.receiverOnlyObjectVal, AuthorizerType)
					return authorizer
				}
			}
		}
	}
	return types.NewErr("Error processing authorizer: %s, %s, %v", function, overload, args)

	// return types.NoSuchOverloadErr()
}

type PathCheck struct {
	receiverOnlyObjectVal
	Checks map[string]*Decision `yaml:"checks,omitempty"`
}

var _ traits.Receiver = &PathCheck{}

func (p *PathCheck) Receive(function string, overload string, args []ref.Val) ref.Val {
	if function == "check" && len(args) == 1 {
		if check, ok := getString(args[0].Value()); ok {
			if len(check) == 0 {
				return types.NewErr("must specify check")
			}
			if p.Checks != nil {
				if decision, ok := p.Checks[check]; ok {
					initReceiver(&decision.receiverOnlyObjectVal, DecisionType)
					return decision
				}
			}
			decision := &Decision{}
			initReceiver(&decision.receiverOnlyObjectVal, DecisionType)
			return decision
		}
		return types.NoSuchOverloadErr()

	}
	return types.NoSuchOverloadErr()
}

type GroupCheck struct {
	receiverOnlyObjectVal
	Resources map[string]*ResourceCheck `json:"resources,omitempty"`
}

var _ traits.Receiver = &GroupCheck{}

func (g *GroupCheck) Receive(function string, overload string, args []ref.Val) ref.Val {
	if function == "resource" && len(args) == 1 {
		if resource, ok := getString(args[0].Value()); ok {
			if len(resource) >= 0 && g.Resources != nil {
				if resourceCheck, ok := g.Resources[resource]; ok {
					initReceiver(&resourceCheck.receiverOnlyObjectVal, ResourceCheckType)
					return resourceCheck
				}
			}
			resourceCheck := &ResourceCheck{}
			initReceiver(&resourceCheck.receiverOnlyObjectVal, ResourceCheckType)
			return resourceCheck
		}
	}
	return types.NoSuchOverloadErr()
}

var _ traits.Receiver = &ResourceCheck{}

type ResourceCheck struct {
	receiverOnlyObjectVal
	namespace     *string
	name          *string
	noSubresource bool
	Subresources  map[string]*ResourceCheck                  `yaml:"subresources,omitempty"`
	Checks        map[string]map[string]map[string]*Decision `yaml:"checks,omitempty"`
}

func (r *ResourceCheck) Receive(function string, overload string, args []ref.Val) ref.Val {
	if len(args) == 1 {
		switch function {
		case "subresource":
			if r.noSubresource {
				return types.NewErr("subresource already invoked")
			}
			if subresource, ok := getString(args[0].Value()); ok {
				if len(subresource) == 0 {
					return r
				} else if r.Subresources != nil {
					if resourceCheck, ok := r.Subresources[subresource]; ok {
						initResourceReceiver(resourceCheck, r.namespace, r.name, true)
						return resourceCheck
					}
				}
				resourceCheck := &ResourceCheck{}
				initResourceReceiver(resourceCheck, nil, nil, true)
				return resourceCheck
			}
		case "namespace":
			if r.namespace != nil {
				return types.NewErr("namespace already invoked")
			}
			if namespace, ok := getString(args[0].Value()); ok {
				resourceCheck := &*r
				initResourceReceiver(resourceCheck, &namespace, resourceCheck.name, r.noSubresource)
				return resourceCheck
			}
		case "name":
			if r.name != nil {
				return types.NewErr("name already invoked")
			}
			if name, ok := getString(args[0].Value()); ok {
				if len(name) == 0 {
					return r
				}
				resourceCheck := &*r
				initResourceReceiver(resourceCheck, resourceCheck.namespace, &name, r.noSubresource)
				return resourceCheck
			}
		case "check":
			return getDecision(args[0].Value(), r.Checks, r.namespace, r.name)
		}
	}
	return types.NoSuchOverloadErr()
}

func getDecision(checkVal any, checks map[string]map[string]map[string]*Decision, namespace *string, name *string) ref.Val {
	if check, ok := getString(checkVal); ok {
		if len(check) == 0 {
			return types.NewErr("must specify check")
		}
		if checks != nil {
			checkNamespace := getValOrEmpty(namespace)
			checkName := getValOrEmpty(name)
			if namespacedChecks, ok := checks[checkNamespace]; ok {
				if namedChecks, ok := namespacedChecks[checkName]; ok {
					if decision, ok := namedChecks[check]; ok {
						initReceiver(&decision.receiverOnlyObjectVal, DecisionType)
						return decision
					}
				}
			}
		}
		decision := &Decision{}
		initReceiver(&decision.receiverOnlyObjectVal, DecisionType)
		return decision
	}
	return types.NoSuchOverloadErr()
}

type Decision struct {
	receiverOnlyObjectVal
	Error    string `yaml:"error,omitempty"`
	Decision string `yaml:"decision,omitempty"`
	Reason   string `yaml:"reason,omitempty"`
}

var _ traits.Receiver = &Decision{}

func (d *Decision) Receive(function string, overload string, args []ref.Val) ref.Val {
	if len(args) == 0 {
		switch function {
		case "errored":
			return types.Bool(d.Error != "")
		case "error":
			return types.String(d.Error)
		case "allowed":
			return types.Bool(d.Decision == "allow")
		case "reason":
			return types.String(d.Reason)
		}
	}
	return types.NoSuchOverloadErr()
}

func initResourceReceiver(resourceCheck *ResourceCheck, namespace *string, name *string, noSubresource bool) {
	resourceCheck.namespace = namespace
	resourceCheck.name = name
	resourceCheck.noSubresource = noSubresource
	initReceiver(&resourceCheck.receiverOnlyObjectVal, ResourceCheckType)
}

func initReceiver(receiver *receiverOnlyObjectVal, varType *types.Type) {
	if receiver.typeValue == nil {
		*receiver = receiverOnlyVal(varType)
	}
}

func getString(val any) (string, bool) {
	if strptr, ok := val.(*string); ok {
		if strptr == nil {
			return "", false
		} else {
			return strings.TrimSpace(*strptr), ok
		}
	} else if str, ok := val.(string); ok {
		return strings.TrimSpace(str), ok
	} else {
		return "", false
	}
}

func getValOrEmpty(val any) string {
	if str, ok := getString(val); ok {
		return str
	} else {
		return ""
	}
}

func getAuthorizerRequestResource(authorizer *Authorizer, request map[string]any) (*ResourceCheck, error) {
	if authorizer == nil || request == nil {
		return nil, nil
	}
	name := getValOrEmpty(request["name"])
	namespace := getValOrEmpty(request["namespace"])
	resourceMap := request["resource"].(map[string]any)
	group := getValOrEmpty(resourceMap["group"])
	resource := getValOrEmpty(resourceMap["resource"])

	receivers := [][2]string{
		{"group", group},
		{"resource", resource},
		{"namespace", namespace},
		{"name", name},
	}

	var receiver traits.Receiver = authorizer
	for _, receiverFunction := range receivers {
		val := receiver.Receive(receiverFunction[0], "", []ref.Val{types.String(receiverFunction[1])})
		if err, ok := val.(*types.Err); ok {
			return nil, err
		}
		receiver = val.(traits.Receiver)
	}
	return receiver.(*ResourceCheck), nil
}
