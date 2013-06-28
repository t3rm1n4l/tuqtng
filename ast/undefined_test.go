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
	"testing"
)

func TestUndefined(t *testing.T) {

	x := Undefined{"property"}
	err := x.Error()

	if err != "property is not defined" {
		t.Errorf("Expected property is not defined, got %v", err)
	}

	y := Undefined{}
	err = y.Error()

	if err != "not defined" {
		t.Errorf("Expected not defined, got %v", err)
	}

}
