// Copyright 2020 Consensys Software Inc.
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

package iop

import (
	"encoding/binary"
	"io"
	"math/big"
	"math/bits"

	"github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
	"github.com/consensys/gnark-crypto/ecc/bw6-761/fr/fft"
)

// Basis indicates the basis in which a polynomial is represented.
type Basis uint32

const (
	Canonical Basis = 1 << iota
	Lagrange
	LagrangeCoset
)

// Layout indicates if a polynomial has a BitReverse or a Regular layout
type Layout uint32

const (
	Regular Layout = 8 << iota
	BitReverse
)

// Form describes the form of a polynomial.
// TODO should be a regular enum?
type Form struct {
	Basis  Basis
	Layout Layout
}

// enum of the possible Form values for type-safe switches
// in this package
var (
	canonicalRegular        = Form{Canonical, Regular}
	canonicalBitReverse     = Form{Canonical, BitReverse}
	lagrangeRegular         = Form{Lagrange, Regular}
	lagrangeBitReverse      = Form{Lagrange, BitReverse}
	lagrangeCosetRegular    = Form{LagrangeCoset, Regular}
	lagrangeCosetBitReverse = Form{LagrangeCoset, BitReverse}
)

// Polynomial wraps a polynomial so that it is
// interpreted as P'(X)=P(\omega^{s}X).
// Dimension is the dimension of the smallest subspace to which the
// polynomial belongs, that is degree(p)+1.
// For instance if len(P)=32 but P.Size=8, it means that P has been
// extended (e.g. it is evaluated on a larger set) but P is a polynomial
// of degree 7.
// blindedSize is the size of the polynomial when it is blinded. By
// default blindedSize=Size, until the polynomial is blinded.
type Polynomial struct {
	*polynomial
	shift       int
	dimension   int
	blindedSize int
}

// NewPolynomial returned a Polynomial from the provided coefficients in the given form.
// A Polynomial can be seen as a "shared pointer" on a list of coefficients.
// It is the responsibility of the user to call the Clone method if the coefficients
// shouldn't be mutated.
func NewPolynomial(coeffs *[]fr.Element, form Form) *Polynomial {
	return &Polynomial{
		polynomial:  newPolynomial(coeffs, form),
		dimension:   len(*coeffs),
		blindedSize: len(*coeffs),
	}
}

// Shift the wrapped polynomial; it doesn't modify the underlying data structure,
// but flag the Polynomial such that it will be interpreted as p(\omega^shift X)
func (p *Polynomial) Shift(shift int) *Polynomial {
	p.shift = shift
	return p
}

// BlindedSize returns the the size of the polynomial when it is blinded. By
// default blindedSize=Size, until the polynomial is blinded.
func (p *Polynomial) BlindedSize() int {
	return p.blindedSize
}

// Size returns the dimension of the smallest subspace to which
// p belongs, that is deg(p)+1
func (p *Polynomial) Dimension() int {
	return p.dimension
}

// Blind blinds a polynomial q by adding Q(X)*(X^{n}-1),
// where deg Q = blindingOrder and Q is random, and n is the
// size of q. Sets the result to p and returns it, as well as the
// blinding factor Q(X).
//
// blindingOrder is the degree of Q, where the blinding is Q(X)*(X^{n}-1)
// where n is the size of p. The size of p is modified since the underlying
// polynomial is of bigger degree now. The new size is p.Size+1+blindingOrder.
//
// /!\ The code panics if wq is not in canonical, regular layout
func (p *Polynomial) Blind(blindingOrder int) (*Polynomial, *Polynomial) {
	// check that p is in canonical basis
	if p.Form != canonicalRegular {
		panic("the input must be in canonical basis, regular layout")
	}

	// we add Q*(x^{n}-1) so the new size is deg(Q)+n+1
	// where n is the size of wq.
	newSize := p.dimension + blindingOrder + 1

	// Resize p. The size of wq might has already been increased
	// (e.g. when the polynomial is evaluated on a larger domain),
	// if that's the case we don't resize the polynomial.
	p.grow(newSize)

	q := make([]fr.Element, blindingOrder+1)

	// blinding: we add Q(X)(X^{n}-1) to P, where deg(Q)=blindingOrder
	var r fr.Element

	for i := 0; i <= blindingOrder; i++ {
		r.SetRandom()
		q[i].Set(&r)
		(*p.coefficients)[i].Sub(&(*p.coefficients)[i], &r)
		(*p.coefficients)[i+p.dimension].Add(&(*p.coefficients)[i+p.dimension], &r)
	}
	p.blindedSize = newSize
	b := NewPolynomial(&q, Form{Basis: Canonical, Layout: Regular})

	return p, b
}

// Evaluate evaluates p at x.
// The code panics if the function is not in canonical form.
func (p *Polynomial) Evaluate(x fr.Element) fr.Element {

	if p.shift == 0 {
		return p.polynomial.evaluate(x)
	}

	var g fr.Element
	if p.shift <= 5 {
		gen, err := fft.Generator(uint64(p.dimension))
		if err != nil {
			panic(err)
		}
		g = smallExp(gen, p.shift)
		x.Mul(&x, &g)
		return p.polynomial.evaluate(x)
	}

	bs := big.NewInt(int64(p.shift))
	g = *g.Exp(g, bs)
	x.Mul(&x, &g)
	return p.polynomial.evaluate(x)
}

// Clone returns a deep copy of p. The underlying polynomial is cloned;
// see also ShallowClone to perform a ShallowClone on the underlying polynomial.
// If capacity is provided, the new coefficient slice capacity will be set accordingly.
func (p *Polynomial) Clone(capacity ...int) *Polynomial {
	res := p.ShallowClone()
	res.polynomial = p.polynomial.clone(capacity...)
	return res
}

// ShallowClone returns a shallow copy of p. The underlying polynomial coefficient
// is NOT cloned and both objects will point to the same coefficient vector.
func (p *Polynomial) ShallowClone() *Polynomial {
	res := *p
	return &res
}

// GetCoeff returns the i-th entry of p, taking the layout in account.
func (p *Polynomial) GetCoeff(i int) fr.Element {

	n := p.coefficients.Len()
	rho := n / p.dimension
	if p.polynomial.Form.Layout == Regular {
		return (*p.coefficients)[(i+rho*p.shift)%n]
	} else {
		nn := uint64(64 - bits.TrailingZeros(uint(n)))
		iRev := bits.Reverse64(uint64((i+rho*p.shift)%n)) >> nn
		return (*p.coefficients)[iRev]
	}

}

// polynomial represents a polynomial, the vector of coefficients
// along with the basis and the layout.
type polynomial struct {
	coefficients *fr.Vector
	Form
}

// Coefficients returns a slice on the underlying data structure.
func (p *polynomial) Coefficients() []fr.Element {
	return (*p.coefficients)
}

// newPolynomial creates a new polynomial. The slice coeff NOT copied
// but directly assigned to the new polynomial.
func newPolynomial(coeffs *[]fr.Element, form Form) *polynomial {
	return &polynomial{coefficients: (*fr.Vector)(coeffs), Form: form}
}

// clone returns a deep copy of the underlying data structure.
func (p *polynomial) clone(capacity ...int) *polynomial {
	c := p.coefficients.Len()
	if len(capacity) == 1 && capacity[0] > c {
		c = capacity[0]
	}
	newCoeffs := make(fr.Vector, p.coefficients.Len(), c)
	r := &polynomial{
		coefficients: &newCoeffs,
		Form:         p.Form,
	}
	copy((*r.coefficients), (*p.coefficients))
	return r
}

// evaluate evaluates p at x.
// The code panics if the function is not in canonical form.
func (p *polynomial) evaluate(x fr.Element) fr.Element {

	var r fr.Element
	if p.Basis != Canonical {
		panic("p must be in canonical basis")
	}

	if p.Layout == Regular {
		for i := p.coefficients.Len() - 1; i >= 0; i-- {
			r.Mul(&r, &x).Add(&r, &(*p.coefficients)[i])
		}
	} else {
		nn := uint64(64 - bits.TrailingZeros(uint(p.coefficients.Len())))
		for i := p.coefficients.Len() - 1; i >= 0; i-- {
			iRev := bits.Reverse64(uint64(i)) >> nn
			r.Mul(&r, &x).Add(&r, &(*p.coefficients)[iRev])
		}
	}

	return r

}

// ToRegular changes the layout of p to Regular.
// Leaves p unchanged if p's layout was already Regular.
func (p *Polynomial) ToRegular() *Polynomial {
	if p.Layout == Regular {
		return p
	}
	fft.BitReverse((*p.coefficients))
	p.Layout = Regular
	return p
}

// ToBitReverse changes the layout of p to BitReverse.
// Leaves p unchanged if p's layout was already BitReverse.
func (p *Polynomial) ToBitReverse() *Polynomial {
	if p.Layout == BitReverse {
		return p
	}
	fft.BitReverse((*p.coefficients))
	p.Layout = BitReverse
	return p
}

// ToLagrange converts p to Lagrange form.
// Leaves p unchanged if p was already in Lagrange form.
func (p *Polynomial) ToLagrange(d *fft.Domain) *Polynomial {
	id := p.Form
	p.grow(int(d.Cardinality))
	switch id {
	case canonicalRegular:
		p.Layout = BitReverse
		d.FFT((*p.coefficients), fft.DIF)
	case canonicalBitReverse:
		p.Layout = Regular
		d.FFT((*p.coefficients), fft.DIT)
	case lagrangeRegular, lagrangeBitReverse:
		return p
	case lagrangeCosetRegular:
		p.Layout = Regular
		d.FFTInverse((*p.coefficients), fft.DIF, fft.OnCoset())
		d.FFT((*p.coefficients), fft.DIT)
	case lagrangeCosetBitReverse:
		p.Layout = BitReverse
		d.FFTInverse((*p.coefficients), fft.DIT, fft.OnCoset())
		d.FFT((*p.coefficients), fft.DIF)
	default:
		panic("unknown ID")
	}
	p.Basis = Lagrange
	return p
}

// ToCanonical converts p to canonical form.
// Leaves p unchanged if p was already in Canonical form.
func (p *Polynomial) ToCanonical(d *fft.Domain) *Polynomial {
	id := p.Form
	p.grow(int(d.Cardinality))
	switch id {
	case canonicalRegular, canonicalBitReverse:
		return p
	case lagrangeRegular:
		p.Layout = BitReverse
		d.FFTInverse((*p.coefficients), fft.DIF)
	case lagrangeBitReverse:
		p.Layout = Regular
		d.FFTInverse((*p.coefficients), fft.DIT)
	case lagrangeCosetRegular:
		p.Layout = BitReverse
		d.FFTInverse((*p.coefficients), fft.DIF, fft.OnCoset())
	case lagrangeCosetBitReverse:
		p.Layout = Regular
		d.FFTInverse((*p.coefficients), fft.DIT, fft.OnCoset())
	default:
		panic("unknown ID")
	}
	p.Basis = Canonical
	return p
}

func (p *polynomial) grow(newSize int) {
	offset := newSize - p.coefficients.Len()
	if offset > 0 {
		(*p.coefficients) = append((*p.coefficients), make(fr.Vector, offset)...)
	}
}

// ToLagrangeCoset Sets p to q, in LagrangeCoset form and returns it.
func (p *Polynomial) ToLagrangeCoset(d *fft.Domain) *Polynomial {
	id := p.Form
	p.grow(int(d.Cardinality))
	switch id {
	case canonicalRegular:
		p.Layout = BitReverse
		d.FFT((*p.coefficients), fft.DIF, fft.OnCoset())
	case canonicalBitReverse:
		p.Layout = Regular
		d.FFT((*p.coefficients), fft.DIT, fft.OnCoset())
	case lagrangeRegular:
		p.Layout = Regular
		d.FFTInverse((*p.coefficients), fft.DIF)
		d.FFT((*p.coefficients), fft.DIT, fft.OnCoset())
	case lagrangeBitReverse:
		p.Layout = BitReverse
		d.FFTInverse((*p.coefficients), fft.DIT)
		d.FFT((*p.coefficients), fft.DIF, fft.OnCoset())
	case lagrangeCosetRegular, lagrangeCosetBitReverse:
		return p
	default:
		panic("unknown ID")
	}

	p.Basis = LagrangeCoset
	return p
}

// WriteTo implements io.WriterTo
func (p *Polynomial) WriteTo(w io.Writer) (int64, error) {
	// encode coefficients
	n, err := p.polynomial.coefficients.WriteTo(w)
	if err != nil {
		return n, err
	}

	// encode Form.Basis, Form.Layout, shift, size & blindedSize as uint32
	var data = []uint32{
		uint32(p.Basis),
		uint32(p.Layout),
		uint32(p.shift),
		uint32(p.dimension),
		uint32(p.blindedSize),
	}
	for _, v := range data {
		err = binary.Write(w, binary.BigEndian, v)
		if err != nil {
			return n, err
		}
		n += 4
	}
	return n, nil
}

// ReadFrom implements io.ReaderFrom
func (p *Polynomial) ReadFrom(r io.Reader) (int64, error) {
	// decode coefficients
	if p.polynomial == nil {
		p.polynomial = new(polynomial)
	}
	if p.polynomial.coefficients == nil {
		v := make(fr.Vector, 0)
		p.polynomial.coefficients = &v
	}
	n, err := p.polynomial.coefficients.ReadFrom(r)
	if err != nil {
		return n, err
	}

	// decode Form.Basis, Form.Layout, shift, size & blindedSize as uint32
	var data [5]uint32
	var buf [4]byte
	for i := range data {
		read, err := io.ReadFull(r, buf[:4])
		n += int64(read)
		if err != nil {
			return n, err
		}
		data[i] = binary.BigEndian.Uint32(buf[:4])
	}

	p.Basis = Basis(data[0])
	p.Layout = Layout(data[1])
	p.shift = int(data[2])
	p.dimension = int(data[3])
	p.blindedSize = int(data[4])

	return n, nil
}
