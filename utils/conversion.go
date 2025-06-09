// Copyright 2025 Undistro Authors
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

package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"google.golang.org/protobuf/types/known/structpb"
)

type conversionTraits interface {
	traits.Iterable
	traits.Indexer
}

var stringType = reflect.TypeOf("")

func ConvertValToNative(val ref.Val) (any, error) {
	valType := val.Type()
	switch valType {
	case types.ListType:
		if iterable, ok := val.(conversionTraits); !ok {
			return nil, errors.New("type conversion error from list to iterable")
		} else {
			values := []any{}
			iter := iterable.Iterator()
			for iter.HasNext() == types.True {
				if value, err := ConvertValToNative(iter.Next()); err != nil {
					return nil, err
				} else {
					values = append(values, value)
				}
			}
			return values, nil
		}
	case types.MapType:
		if iterable, ok := val.(conversionTraits); !ok {
			return nil, errors.New("type conversion error from map to iterable")
		} else {
			values := map[string]any{}
			iter := iterable.Iterator()
			for iter.HasNext() == types.True {
				keyVal := iter.Next()
				if key, err := keyVal.ConvertToNative(stringType); err != nil {
					return nil, fmt.Errorf("unexpected map key type: %v", keyVal.Type())
				} else if value, err := ConvertValToNative(iterable.Get(keyVal)); err != nil {
					return nil, err
				} else {
					values[key.(string)] = value
				}
			}
			return values, nil
		}
	case types.OptionalType:
		opt, ok := val.(*types.Optional)
		if !ok {
			return nil, errors.New("type conversion error for optional")
		} else if !opt.HasValue() {
			return nil, nil
		}
		val = opt.GetValue()
		fallthrough
	default:
		if value, err := val.ConvertToNative(reflect.TypeOf(&structpb.Value{})); err != nil {
			return nil, err
		} else {
			return value, nil
		}
	}
}
