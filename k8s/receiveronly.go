/*
Copyright 2023 The Kubernetes Authors.

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

package k8s

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// receiverOnlyObjectVal provides an implementation of ref.Val for
// any object type that has receiver functions but does not expose any fields to
// CEL.
type receiverOnlyObjectVal struct {
	typeValue *types.Type
}

// receiverOnlyVal returns a receiverOnlyObjectVal for the given type.
func receiverOnlyVal(objectType *cel.Type) receiverOnlyObjectVal {
	return receiverOnlyObjectVal{typeValue: types.NewTypeValue(objectType.String(), traits.ReceiverType)}
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (a receiverOnlyObjectVal) ConvertToNative(typeDesc reflect.Type) (any, error) {
	return nil, fmt.Errorf("type conversion error from '%s' to '%v'", a.typeValue.String(), typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (a receiverOnlyObjectVal) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case a.typeValue:
		return a
	case types.TypeType:
		return a.typeValue
	}
	return types.NewErr("type conversion error from '%s' to '%s'", a.typeValue, typeVal)
}

// Equal implements ref.Val.Equal.
func (a receiverOnlyObjectVal) Equal(other ref.Val) ref.Val {
	o, ok := other.(receiverOnlyObjectVal)
	if !ok {
		return types.MaybeNoSuchOverloadErr(other)
	}
	return types.Bool(a == o)
}

// Type implements ref.Val.Type.
func (a receiverOnlyObjectVal) Type() ref.Type {
	return a.typeValue
}

// Value implements ref.Val.Value.
func (a receiverOnlyObjectVal) Value() any {
	return types.NoSuchOverloadErr()
}
