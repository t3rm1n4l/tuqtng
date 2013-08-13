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

	"github.com/couchbaselabs/dparval"
)

type BinaryOperator struct {
	operator string
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

func (this *BinaryOperator) compare(item *dparval.Value) (*dparval.Value, error) {
	lv, err := this.Left.Evaluate(item)
	if err != nil {
		// this could either be real error, or MISSING
		// if either side is MISSING, the result is MISSING
		return nil, err
	}
	// if either side is NULL, the result is NULL
	if lv.Type() == dparval.NULL {
		return nil, nil
	}
	rv, err := this.Right.Evaluate(item)
	if err != nil {
		// this could either be real error, or MISSING
		// if either side is MISSING, the result is MISSING
		return nil, err
	}
	// if either side is NULL, the result is NULL
	if rv.Type() == dparval.NULL {
		return nil, nil
	}

	lvalue := lv.Value()
	rvalue := rv.Value()

	// if we got this far, we evaluated both sides
	// there were no errors, and neither side was NULL or MISSING
	// now check types (types must be the same)
	ltype := collationType(lvalue)
	rtype := collationType(rvalue)
	// ugly fixups for boolean (returns different values for true/false)
	if ltype == 2 {
		// fixup for boolean type
		ltype = 1
	}
	if rtype == 2 {
		rtype = 1
	}

	if ltype != rtype {
		return nil, &TypeMismatch{ltype, rtype}
	}

	return dparval.NewValue(float64(CollateJSON(lvalue, rvalue))), nil
}

func (this *BinaryOperator) Dependencies() ExpressionList {
	rv := ExpressionList{this.Left, this.Right}
	return rv
}

func (this *BinaryOperator) GetLeft() Expression {
	return this.Left
}

func (this *BinaryOperator) GetRight() Expression {
	return this.Right
}

func (this *BinaryOperator) SetLeft(left Expression) {
	this.Left = left
}

func (this *BinaryOperator) SetRight(right Expression) {
	this.Right = right
}

func (this *BinaryOperator) Operator() string {
	return this.operator
}

func (this *BinaryOperator) String() string {
	return fmt.Sprintf("%v %s %v", this.Left, this.operator, this.Right)
}

func (this *BinaryOperator) EquivalentTo(t Expression) bool {
	// another binary operator?
	that, ok := t.(BinaryOperatorExpression)
	if !ok {
		return false
	}

	// same type of operator?
	if this.operator != that.Operator() {
		return false
	}

	// same operands?
	if this.Left.EquivalentTo(that.GetLeft()) &&
		this.Right.EquivalentTo(that.GetRight()) {
		return true
	}

	return false
}