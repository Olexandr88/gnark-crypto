//go:build !amd64 && !arm64
// +build !amd64,!arm64

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

func add(z, x, y *Element) {
	_addGeneric(z, x, y)
}

func sub(z, x, y *Element) {
	_subGeneric(z, x, y)
}

func double(z, x *Element) {
	_doubleGeneric(z, x)
}

func neg(z, x *Element) {
	_negGeneric(z, x)
}

func mul(z, x, y *Element) {
	_mulGeneric(z, x, y)
}
