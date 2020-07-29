package fq12over6over2

const Fq2Common = `

import (
	"github.com/consensys/gurvy/{{toLower .CurveName}}/fp"
)

// E2 is a degree two finite field extension of fp.Element
type E2 struct {
	A0, A1 fp.Element
}

// Equal returns true if z equals x, fasle otherwise
func (z *E2) Equal(x *E2) bool {
	return z.A0.Equal(&x.A0) && z.A1.Equal(&x.A1)
}

// SetString sets a E2 element from strings
func (z *E2) SetString(s1, s2 string) *E2 {
	z.A0.SetString(s1)
	z.A1.SetString(s2)
	return z
}

// SetZero sets an e2 elmt to zero
func (z *E2) SetZero() *E2 {
	z.A0.SetZero()
	z.A1.SetZero()
	return z
}

// Clone returns a copy of self
func (z *E2) Clone() *E2 {
	return &E2{
		A0: z.A0,
		A1: z.A1,
	}
}

// Set sets an E2 from x
func (z *E2) Set(x *E2) *E2 {
	z.A0.Set(&x.A0)
	z.A1.Set(&x.A1)
	return z
}

// SetOne sets z to 1 in Montgomery form and returns z
func (z *E2) SetOne() *E2 {
	z.A0.SetOne()
	z.A1.SetZero()
	return z
}

// SetRandom sets a0 and a1 to random values
func (z *E2) SetRandom() *E2 {
	z.A0.SetRandom()
	z.A1.SetRandom()
	return z
}

// IsZero returns true if the two elements are equal, fasle otherwise
func (z *E2) IsZero() bool {
	return z.A0.IsZero() && z.A1.IsZero()
}

// Neg negates an E2 element
func (z *E2) Neg(x *E2) *E2 {
	z.A0.Neg(&x.A0)
	z.A1.Neg(&x.A1)
	return z
}

// String implements Stringer interface for fancy printing
func (z *E2) String() string {
	return (z.A0.String() + "+" + z.A1.String() + "*u")
}

// ToMont converts to mont form
func (z *E2) ToMont() *E2 {
	z.A0.ToMont()
	z.A1.ToMont()
	return z
}

// FromMont converts from mont form
func (z *E2) FromMont() *E2 {
	z.A0.FromMont()
	z.A1.FromMont()
	return z
}

// Add adds two elements of E2
func (z *E2) Add(x, y *E2) *E2 {
	z.A0.Add(&x.A0, &y.A0)
	z.A1.Add(&x.A1, &y.A1)
	return z
}

// AddAssign adds x to z
func (z *E2) AddAssign(x *E2) *E2 {
	z.A0.AddAssign(&x.A0)
	z.A1.AddAssign(&x.A1)
	return z
}

// Sub two elements of E2
func (z *E2) Sub(x, y *E2) *E2 {
	z.A0.Sub(&x.A0, &y.A0)
	z.A1.Sub(&x.A1, &y.A1)
	return z
}

// SubAssign subs x from z
func (z *E2) SubAssign(x *E2) *E2 {
	z.A0.SubAssign(&x.A0)
	z.A1.SubAssign(&x.A1)
	return z
}

// Double doubles an E2 element
func (z *E2) Double(x *E2) *E2 {
	z.A0.Double(&x.A0)
	z.A1.Double(&x.A1)
	return z
}

// MulAssign sets z to the E2 product of z,x returns z
func (z *E2) MulAssign(x *E2) *E2 {
	z.Mul(z, x)
	return z
}

// MulByElement multiplies an element in E2 by an element in fp
func (z *E2) MulByElement(x *E2, y *fp.Element) *E2 {
	var yCopy fp.Element
	yCopy.Set(y)
	z.A0.Mul(&x.A0, &yCopy)
	z.A1.Mul(&x.A1, &yCopy)
	return z
}

// Conjugate conjugates an element in E2
func (z *E2) Conjugate(x *E2) *E2 {
	z.A0.Set(&x.A0)
	z.A1.Neg(&x.A1)
	return z
}

//-----------------------------------------------------------------------------
// Specific to {{.CurveName}}

{{- if eq .CurveName "bn256" }}
	{{- template "bn256" }}
{{- else if eq .CurveName "bls381" }}
	{{- template "bls381" }}
{{- else if eq .CurveName "bls377" }}
	{{- template "bls377" }}
{{- else }}
	// TODO implement Mul, Square, MulByNonResidue, MulByNonResidueInv, Inverse
{{- end }}
`

const Fq2Specific = `

{{- define "bn256" }}

	// Mul sets z to the E2-product of x,y, returns z
	func (z *E2) Mul(x, y *E2) *E2 {
		var a, b, c fp.Element
		a.Add(&x.A0, &x.A1)
		b.Add(&y.A0, &y.A1)
		a.Mul(&a, &b)
		b.Mul(&x.A0, &y.A0)
		c.Mul(&x.A1, &y.A1)
		z.A1.Sub(&a, &b).Sub(&z.A1, &c)
		z.A0.Sub(&b, &c) //z.A0.MulByNonResidue(&c).Add(&z.A0, &b)
		return z
	}

	// Square sets z to the E2-product of x,x returns z
	func (z *E2) Square(x *E2) *E2 {
		// algo 22 https://eprint.iacr.org/2010/354.pdf
		var a, b fp.Element
		a.Add(&x.A0, &x.A1)
		b.Sub(&x.A0, &x.A1)
		a.Mul(&a, &b)
		b.Mul(&x.A0, &x.A1).Double(&b)
		z.A0.Set(&a)
		z.A1.Set(&b)
		return z
	}

	// MulByNonResidue multiplies a E2 by (9,1)
	func (z *E2) MulByNonResidue(x *E2) *E2 {
		var a, b fp.Element
		a.Double(&x.A0).Double(&a).Double(&a).Add(&a, &x.A0).Sub(&a, &x.A1)
		b.Double(&x.A1).Double(&b).Double(&b).Add(&b, &x.A1).Add(&b, &x.A0)
		z.A0.Set(&a)
		z.A1.Set(&b)
		return z
	}

	// MulByNonResidueInv multiplies a E2 by (9,1)^{-1}
	func (z *E2) MulByNonResidueInv(x *E2) *E2 {

		var nonresinv E2
		nonresinv.A0 = fp.Element{
			10477841894441615122,
			7327163185667482322,
			3635199979766503006,
			3215324977242306624,
		}
		nonresinv.A1 = fp.Element{
			7515750141297360845,
			14746352163864140223,
			11319968037783994424,
			30185921062296004,
		}
		z.Mul(x, &nonresinv)

		return z
	}

	// Inverse sets z to the E2-inverse of x, returns z
	func (z *E2) Inverse(x *E2) *E2 {
		// Algorithm 8 from https://eprint.iacr.org/2010/354.pdf
		var t0, t1 fp.Element
		t0.Square(&x.A0)
		t1.Square(&x.A1)
		t0.Add(&t0, &t1)
		t1.Inverse(&t0)
		z.A0.Mul(&x.A0, &t1)
		z.A1.Mul(&x.A1, &t1).Neg(&z.A1)

		return z
	}

{{- end }}

{{- define "bls381" }}

	// Mul sets z to the E2-product of x,y, returns z
	func (z *E2) Mul(x, y *E2) *E2 {
		var a, b, c fp.Element
		a.Add(&x.A0, &x.A1)
		b.Add(&y.A0, &y.A1)
		a.Mul(&a, &b)
		b.Mul(&x.A0, &y.A0)
		c.Mul(&x.A1, &y.A1)
		z.A1.Sub(&a, &b).Sub(&z.A1, &c)
		z.A0.Sub(&b, &c) //z.A0.MulByNonResidue(&c).Add(&z.A0, &b)
		return z
	}

	// Square sets z to the E2-product of x,x returns z
	func (z *E2) Square(x *E2) *E2 {
		// algo 22 https://eprint.iacr.org/2010/354.pdf
		var a, b fp.Element
		a.Add(&x.A0, &x.A1)
		b.Sub(&x.A0, &x.A1)
		a.Mul(&a, &b)
		b.Mul(&x.A0, &x.A1).Double(&b)
		z.A0.Set(&a)
		z.A1.Set(&b)
		return z
	}

	// MulByNonResidue multiplies a E2 by (1,1)
	func (z *E2) MulByNonResidue(x *E2) *E2 {
		var a fp.Element
		a.Sub(&x.A0, &x.A1)
		z.A1.Add(&x.A0, &x.A1)
		z.A0.Set(&a)
		return z
	}

	// MulByNonResidueInv multiplies a E2 by (1,1)^{-1}
	func (z *E2) MulByNonResidueInv(x *E2) *E2 {

		twoinv := fp.Element{
			1730508156817200468,
			9606178027640717313,
			7150789853162776431,
			7936136305760253186,
			15245073033536294050,
			1728177566264616342,
		}

		var tmp fp.Element
		tmp.Add(&x.A0, &x.A1)
		z.A1.Sub(&x.A1, &x.A0).Mul(&z.A1, &twoinv)
		z.A0.Set(&tmp).Mul(&z.A0, &twoinv)

		return z
	}

	// Inverse sets z to the E2-inverse of x, returns z
	func (z *E2) Inverse(x *E2) *E2 {
		// Algorithm 8 from https://eprint.iacr.org/2010/354.pdf
		var t0, t1 fp.Element
		t0.Square(&x.A0)
		t1.Square(&x.A1)
		t0.Add(&t0, &t1)
		t1.Inverse(&t0)
		z.A0.Mul(&x.A0, &t1)
		z.A1.Mul(&x.A1, &t1).Neg(&z.A1)

		return z
	}

{{- end }}

{{- define "bls377" }}

	// Mul sets z to the E2-product of x,y, returns z
	func (z *E2) Mul(x, y *E2) *E2 {
		var a, b, c fp.Element
		a.Add(&x.A0, &x.A1)
		b.Add(&y.A0, &y.A1)
		a.Mul(&a, &b)
		b.Mul(&x.A0, &y.A0)
		c.Mul(&x.A1, &y.A1)
		z.A1.Sub(&a, &b).Sub(&z.A1, &c)
		z.A0.Double(&c).Double(&z.A0).AddAssign(&c).Add(&z.A0, &b)
		return z
	}

	// Square sets z to the E2-product of x,x returns z
	func (z *E2) Square(x *E2) *E2 {
		//algo 22 https://eprint.iacr.org/2010/354.pdf
		var c0, c2 fp.Element
		c2.Double(&x.A1).Double(&c2).AddAssign(&x.A1).AddAssign(&x.A0)
		c0.Add(&x.A0, &x.A1)
		c0.Mul(&c0, &c2) // (x1+x2)*(x1+(u**2)x2)
		z.A1.Mul(&x.A0, &x.A1).Double(&z.A1)
		z.A0.Sub(&c0, &z.A1).SubAssign(&z.A1).SubAssign(&z.A1)

		return z
	}

	// MulByNonResidue multiplies a E2 by (0,1)
	func (z *E2) MulByNonResidue(x *E2) *E2 {
		a := x.A0
		b := x.A1 // fetching x.A1 in the function below is slower
		z.A0.Double(&b).Double(&z.A0).Add(&z.A0, &b)
		z.A1 = a
		return z
	}

	// MulByNonResidueInv multiplies a E2 by (0,1)^{-1}
	func (z *E2) MulByNonResidueInv(x *E2) *E2 {
		//z.A1.MulByNonResidueInv(&x.A0)
		a := x.A1
		fiveinv := fp.Element{
			330620507644336508,
			9878087358076053079,
			11461392860540703536,
			6973035786057818995,
			8846909097162646007,
			104838758629667239,
		}
		z.A1.Mul(&x.A0, &fiveinv)
		z.A0 = a
		return z
	}

	// Inverse sets z to the E2-inverse of x, returns z
	func (z *E2) Inverse(x *E2) *E2 {
		// Algorithm 8 from https://eprint.iacr.org/2010/354.pdf
		//var a, b, t0, t1, tmp fp.Element
		var t0, t1, tmp fp.Element
		a := &x.A0 // creating the buffers a, b is faster than querying &x.A0, &x.A1 in the functions call below
		b := &x.A1
		t0.Square(a)
		t1.Square(b)
		tmp.Double(&t1).Double(&tmp).Add(&tmp, &t1)
		t0.Sub(&t0, &tmp)
		t1.Inverse(&t0)
		z.A0.Mul(a, &t1)
		z.A1.Mul(b, &t1).Neg(&z.A1)

		return z
	}

{{- end }}

`
