//go:build !amd64
// +build !amd64

// Copyright 2020 ConsenSys Software Inc.
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

// Code generated by consensys/gnark-crypto DO NOT EDIT

package little

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

// MulBy3 x *= 3
func MulBy3(x *Element) {
	mulByConstant(x, 3)
}

// MulBy5 x *= 5
func MulBy5(x *Element) {
	mulByConstant(x, 5)
}

// MulBy13 x *= 13
func MulBy13(x *Element) {
	mulByConstant(x, 13)
}

// Butterfly sets
// a = a + b
// b = a - b
func Butterfly(a, b *Element) {
	_butterflyGeneric(a, b)
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func fromMont(z *Element) {
	_fromMontGeneric(z)
}

/*func add(z,  x, y *Element) {
	_addGeneric(z,x,y)
}

func double(z,  x *Element) {
	_doubleGeneric(z,x)
}


func sub(z,  x, y *Element) {
	_subGeneric(z,x,y)
}

func neg(z,  x *Element) {
	_negGeneric(z,x)
}

func mul(z, x, y *Element) {
	_mulGeneric(z, x, y)
}

*/

func reduce(z *Element) {
	_reduceGeneric(z)
}
