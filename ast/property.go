//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package ast

import (
	"fmt"
)

type Property struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func NewProperty(path string) *Property {
	return &Property{
		Type: "property",
		Path: path,
	}
}

func (this *Property) Evaluate(item Item) (Value, error) {
	if item == nil {
		return nil, &Undefined{this.Path}
	}
	rv, err := item.GetPath(this.Path)
	if err != nil {
		return nil, err
	}
	return rv, nil
}

func (this *Property) String() string {
	return fmt.Sprintf("%v", this.Path)
}