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

package bw6756

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bw6-756/fp"

	"github.com/consensys/gnark-crypto/ecc/bw6-756/fr"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func TestG2AffineEndomorphism(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	properties.Property("[BW6-756] check that phi(P) = lambdaGLV * P", prop.ForAll(
		func(a fp.Element) bool {
			var p, res1, res2 G2Jac
			g := MapToCurveG2Svdw(a)
			p.FromAffine(&g)
			res1.phi(&p)
			res2.mulWindowed(&p, &lambdaGLV)

			return p.IsInSubGroup() && res1.Equal(&res2)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] check that phi^2(P) + phi(P) + P = 0", prop.ForAll(
		func(a fp.Element) bool {
			var p, res, tmp G2Jac
			g := MapToCurveG2Svdw(a)
			p.FromAffine(&g)
			tmp.phi(&p)
			res.phi(&tmp).
				AddAssign(&tmp).
				AddAssign(&p)

			return res.Z.IsZero()
		},
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestMapToCurveG2(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	properties.Property("[G2] Svsw mapping should output point on the curve", prop.ForAll(
		func(a fp.Element) bool {
			g := MapToCurveG2Svdw(a)
			return g.IsInSubGroup()
		},
		GenFp(),
	))

	properties.Property("[G2] Svsw mapping should be deterministic", prop.ForAll(
		func(a fp.Element) bool {
			g1 := MapToCurveG2Svdw(a)
			g2 := MapToCurveG2Svdw(a)
			return g1.Equal(&g2)
		},
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG2AffineIsOnCurve(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	properties.Property("[BW6-756] g2Gen (affine) should be on the curve", prop.ForAll(
		func(a fp.Element) bool {
			var op1, op2 G2Affine
			op1.FromJacobian(&g2Gen)
			op2.FromJacobian(&g2Gen)
			op2.Y.Mul(&op2.Y, &a)
			return op1.IsOnCurve() && !op2.IsOnCurve()
		},
		GenFp(),
	))

	properties.Property("[BW6-756] g2Gen (Jacobian) should be on the curve", prop.ForAll(
		func(a fp.Element) bool {
			var op1, op2, op3 G2Jac
			op1.Set(&g2Gen)
			op3.Set(&g2Gen)

			op2 = fuzzJacobianG2Affine(&g2Gen, a)
			op3.Y.Mul(&op3.Y, &a)
			return op1.IsOnCurve() && op2.IsOnCurve() && !op3.IsOnCurve()
		},
		GenFp(),
	))

	properties.Property("[BW6-756] IsInSubGroup and MulBy subgroup order should be the same", prop.ForAll(
		func(a fp.Element) bool {
			var op1, op2 G2Jac
			op1 = fuzzJacobianG2Affine(&g2Gen, a)
			_r := fr.Modulus()
			op2.ScalarMultiplication(&op1, _r)
			return op1.IsInSubGroup() && op2.Z.IsZero()
		},
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG2AffineConversions(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	properties.Property("[BW6-756] Affine representation should be independent of the Jacobian representative", prop.ForAll(
		func(a fp.Element) bool {
			g := fuzzJacobianG2Affine(&g2Gen, a)
			var op1 G2Affine
			op1.FromJacobian(&g)
			return op1.X.Equal(&g2Gen.X) && op1.Y.Equal(&g2Gen.Y)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] Affine representation should be independent of a Extended Jacobian representative", prop.ForAll(
		func(a fp.Element) bool {
			var g g2JacExtended
			g.X.Set(&g2Gen.X)
			g.Y.Set(&g2Gen.Y)
			g.ZZ.Set(&g2Gen.Z)
			g.ZZZ.Set(&g2Gen.Z)
			gfuzz := fuzzExtendedJacobianG2Affine(&g, a)

			var op1 G2Affine
			op1.fromJacExtended(&gfuzz)
			return op1.X.Equal(&g2Gen.X) && op1.Y.Equal(&g2Gen.Y)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] Jacobian representation should be the same as the affine representative", prop.ForAll(
		func(a fp.Element) bool {
			var g G2Jac
			var op1 G2Affine
			op1.X.Set(&g2Gen.X)
			op1.Y.Set(&g2Gen.Y)

			var one fp.Element
			one.SetOne()

			g.FromAffine(&op1)

			return g.X.Equal(&g2Gen.X) && g.Y.Equal(&g2Gen.Y) && g.Z.Equal(&one)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] Converting affine symbol for infinity to Jacobian should output correct infinity in Jacobian", prop.ForAll(
		func() bool {
			var g G2Affine
			g.X.SetZero()
			g.Y.SetZero()
			var op1 G2Jac
			op1.FromAffine(&g)
			var one, zero fp.Element
			one.SetOne()
			return op1.X.Equal(&one) && op1.Y.Equal(&one) && op1.Z.Equal(&zero)
		},
	))

	properties.Property("[BW6-756] Converting infinity in extended Jacobian to affine should output infinity symbol in Affine", prop.ForAll(
		func() bool {
			var g G2Affine
			var op1 g2JacExtended
			var zero fp.Element
			op1.X.Set(&g2Gen.X)
			op1.Y.Set(&g2Gen.Y)
			g.fromJacExtended(&op1)
			return g.X.Equal(&zero) && g.Y.Equal(&zero)
		},
	))

	properties.Property("[BW6-756] Converting infinity in extended Jacobian to Jacobian should output infinity in Jacobian", prop.ForAll(
		func() bool {
			var g G2Jac
			var op1 g2JacExtended
			var zero, one fp.Element
			one.SetOne()
			op1.X.Set(&g2Gen.X)
			op1.Y.Set(&g2Gen.Y)
			g.fromJacExtended(&op1)
			return g.X.Equal(&one) && g.Y.Equal(&one) && g.Z.Equal(&zero)
		},
	))

	properties.Property("[BW6-756] [Jacobian] Two representatives of the same class should be equal", prop.ForAll(
		func(a, b fp.Element) bool {
			op1 := fuzzJacobianG2Affine(&g2Gen, a)
			op2 := fuzzJacobianG2Affine(&g2Gen, b)
			return op1.Equal(&op2)
		},
		GenFp(),
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG2AffineOps(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	genScalar := GenFr()

	properties.Property("[BW6-756] [Jacobian] Add should call double when having adding the same point", prop.ForAll(
		func(a, b fp.Element) bool {
			fop1 := fuzzJacobianG2Affine(&g2Gen, a)
			fop2 := fuzzJacobianG2Affine(&g2Gen, b)
			var op1, op2 G2Jac
			op1.Set(&fop1).AddAssign(&fop2)
			op2.Double(&fop2)
			return op1.Equal(&op2)
		},
		GenFp(),
		GenFp(),
	))

	properties.Property("[BW6-756] [Jacobian] Adding the opposite of a point to itself should output inf", prop.ForAll(
		func(a, b fp.Element) bool {
			fop1 := fuzzJacobianG2Affine(&g2Gen, a)
			fop2 := fuzzJacobianG2Affine(&g2Gen, b)
			fop2.Neg(&fop2)
			fop1.AddAssign(&fop2)
			return fop1.Equal(&g2Infinity)
		},
		GenFp(),
		GenFp(),
	))

	properties.Property("[BW6-756] [Jacobian] Adding the inf to a point should not modify the point", prop.ForAll(
		func(a fp.Element) bool {
			fop1 := fuzzJacobianG2Affine(&g2Gen, a)
			fop1.AddAssign(&g2Infinity)
			var op2 G2Jac
			op2.Set(&g2Infinity)
			op2.AddAssign(&g2Gen)
			return fop1.Equal(&g2Gen) && op2.Equal(&g2Gen)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] [Jacobian Extended] addMixed (-G) should equal subMixed(G)", prop.ForAll(
		func(a fp.Element) bool {
			fop1 := fuzzJacobianG2Affine(&g2Gen, a)
			var p1, p1Neg G2Affine
			p1.FromJacobian(&fop1)
			p1Neg = p1
			p1Neg.Y.Neg(&p1Neg.Y)
			var o1, o2 g2JacExtended
			o1.addMixed(&p1Neg)
			o2.subMixed(&p1)

			return o1.X.Equal(&o2.X) &&
				o1.Y.Equal(&o2.Y) &&
				o1.ZZ.Equal(&o2.ZZ) &&
				o1.ZZZ.Equal(&o2.ZZZ)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] [Jacobian Extended] doubleMixed (-G) should equal doubleNegMixed(G)", prop.ForAll(
		func(a fp.Element) bool {
			fop1 := fuzzJacobianG2Affine(&g2Gen, a)
			var p1, p1Neg G2Affine
			p1.FromJacobian(&fop1)
			p1Neg = p1
			p1Neg.Y.Neg(&p1Neg.Y)
			var o1, o2 g2JacExtended
			o1.doubleMixed(&p1Neg)
			o2.doubleNegMixed(&p1)

			return o1.X.Equal(&o2.X) &&
				o1.Y.Equal(&o2.Y) &&
				o1.ZZ.Equal(&o2.ZZ) &&
				o1.ZZZ.Equal(&o2.ZZZ)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] [Jacobian] Addmix the negation to itself should output 0", prop.ForAll(
		func(a fp.Element) bool {
			fop1 := fuzzJacobianG2Affine(&g2Gen, a)
			fop1.Neg(&fop1)
			var op2 G2Affine
			op2.FromJacobian(&g2Gen)
			fop1.AddMixed(&op2)
			return fop1.Equal(&g2Infinity)
		},
		GenFp(),
	))

	properties.Property("[BW6-756] scalar multiplication (double and add) should depend only on the scalar mod r", prop.ForAll(
		func(s fr.Element) bool {

			r := fr.Modulus()
			var g G2Jac
			g.mulGLV(&g2Gen, r)

			var scalar, blindedScalar, rminusone big.Int
			var op1, op2, op3, gneg G2Jac
			rminusone.SetUint64(1).Sub(r, &rminusone)
			op3.mulWindowed(&g2Gen, &rminusone)
			gneg.Neg(&g2Gen)
			s.ToBigIntRegular(&scalar)
			blindedScalar.Mul(&scalar, r).Add(&blindedScalar, &scalar)
			op1.mulWindowed(&g2Gen, &scalar)
			op2.mulWindowed(&g2Gen, &blindedScalar)

			return op1.Equal(&op2) && g.Equal(&g2Infinity) && !op1.Equal(&g2Infinity) && gneg.Equal(&op3)

		},
		genScalar,
	))

	properties.Property("[BW6-756] scalar multiplication (GLV) should depend only on the scalar mod r", prop.ForAll(
		func(s fr.Element) bool {

			r := fr.Modulus()
			var g G2Jac
			g.mulGLV(&g2Gen, r)

			var scalar, blindedScalar, rminusone big.Int
			var op1, op2, op3, gneg G2Jac
			rminusone.SetUint64(1).Sub(r, &rminusone)
			op3.ScalarMultiplication(&g2Gen, &rminusone)
			gneg.Neg(&g2Gen)
			s.ToBigIntRegular(&scalar)
			blindedScalar.Mul(&scalar, r).Add(&blindedScalar, &scalar)
			op1.ScalarMultiplication(&g2Gen, &scalar)
			op2.ScalarMultiplication(&g2Gen, &blindedScalar)

			return op1.Equal(&op2) && g.Equal(&g2Infinity) && !op1.Equal(&g2Infinity) && gneg.Equal(&op3)

		},
		genScalar,
	))

	properties.Property("[BW6-756] GLV and Double and Add should output the same result", prop.ForAll(
		func(s fr.Element) bool {

			var r big.Int
			var op1, op2 G2Jac
			s.ToBigIntRegular(&r)
			op1.mulWindowed(&g2Gen, &r)
			op2.mulGLV(&g2Gen, &r)
			return op1.Equal(&op2) && !op1.Equal(&g2Infinity)

		},
		genScalar,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG2AffineCofactorCleaning(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	properties.Property("[BW6-756] Clearing the cofactor of a random point should set it in the r-torsion", prop.ForAll(
		func() bool {
			var a, x, b fp.Element
			a.SetRandom()

			x.Square(&a).Mul(&x, &a).Add(&x, &bTwistCurveCoeff)

			for x.Legendre() != 1 {
				a.SetRandom()

				x.Square(&a).Mul(&x, &a).Add(&x, &bTwistCurveCoeff)

			}

			b.Sqrt(&x)
			var point, pointCleared, infinity G2Jac
			point.X.Set(&a)
			point.Y.Set(&b)
			point.Z.SetOne()
			pointCleared.ClearCofactor(&point)
			infinity.Set(&g2Infinity)
			return point.IsOnCurve() && pointCleared.IsInSubGroup() && !pointCleared.Equal(&infinity)
		},
	))
	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

func TestG2AffineBatchScalarMultiplication(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	genScalar := GenFr()

	// size of the multiExps
	const nbSamples = 10

	properties.Property("[BW6-756] BatchScalarMultiplication should be consistant with individual scalar multiplications", prop.ForAll(
		func(mixer fr.Element) bool {
			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer).
					FromMont()
			}

			result := BatchScalarMultiplicationG2(&g2GenAff, sampleScalars[:])

			if len(result) != len(sampleScalars) {
				return false
			}

			for i := 0; i < len(result); i++ {
				var expectedJac G2Jac
				var expected G2Affine
				var b big.Int
				expectedJac.mulGLV(&g2Gen, sampleScalars[i].ToBigInt(&b))
				expected.FromJacobian(&expectedJac)
				if !result[i].Equal(&expected) {
					return false
				}
			}
			return true
		},
		genScalar,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// ------------------------------------------------------------
// benches

func BenchmarkG2JacIsInSubGroup(b *testing.B) {
	var a G2Jac
	a.Set(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.IsInSubGroup()
	}

}

func BenchmarkG2AffineBatchScalarMul(b *testing.B) {
	// ensure every words of the scalars are filled
	var mixer fr.Element
	mixer.SetString("7716837800905789770901243404444209691916730933998574719964609384059111546487")

	const pow = 15
	const nbSamples = 1 << pow

	var sampleScalars [nbSamples]fr.Element

	for i := 1; i <= nbSamples; i++ {
		sampleScalars[i-1].SetUint64(uint64(i)).
			Mul(&sampleScalars[i-1], &mixer).
			FromMont()
	}

	for i := 5; i <= pow; i++ {
		using := 1 << i

		b.Run(fmt.Sprintf("%d points", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				_ = BatchScalarMultiplicationG2(&g2GenAff, sampleScalars[:using])
			}
		})
	}
}

func BenchmarkG2JacScalarMul(b *testing.B) {

	var scalar big.Int
	r := fr.Modulus()
	scalar.SetString("5243587517512619047944770508185965837690552500527637822603658699938581184513", 10)
	scalar.Add(&scalar, r)

	var doubleAndAdd G2Jac

	b.Run("double and add", func(b *testing.B) {
		b.ResetTimer()
		for j := 0; j < b.N; j++ {
			doubleAndAdd.mulWindowed(&g2Gen, &scalar)
		}
	})

	var glv G2Jac
	b.Run("GLV", func(b *testing.B) {
		b.ResetTimer()
		for j := 0; j < b.N; j++ {
			glv.mulGLV(&g2Gen, &scalar)
		}
	})

}

func BenchmarkG2AffineCofactorClearing(b *testing.B) {
	var a G2Jac
	a.Set(&g2Gen)
	for i := 0; i < b.N; i++ {
		a.ClearCofactor(&a)
	}
}

func BenchmarkG2JacAdd(b *testing.B) {
	var a G2Jac
	a.Double(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.AddAssign(&g2Gen)
	}
}

func BenchmarkG2JacAddMixed(b *testing.B) {
	var a G2Jac
	a.Double(&g2Gen)

	var c G2Affine
	c.FromJacobian(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.AddMixed(&c)
	}

}

func BenchmarkG2JacDouble(b *testing.B) {
	var a G2Jac
	a.Set(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.DoubleAssign()
	}

}

func BenchmarkG2JacExtAddMixed(b *testing.B) {
	var a g2JacExtended
	a.doubleMixed(&g2GenAff)

	var c G2Affine
	c.FromJacobian(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.addMixed(&c)
	}
}

func BenchmarkG2JacExtSubMixed(b *testing.B) {
	var a g2JacExtended
	a.doubleMixed(&g2GenAff)

	var c G2Affine
	c.FromJacobian(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.subMixed(&c)
	}
}

func BenchmarkG2JacExtDoubleMixed(b *testing.B) {
	var a g2JacExtended
	a.doubleMixed(&g2GenAff)

	var c G2Affine
	c.FromJacobian(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.doubleMixed(&c)
	}
}

func BenchmarkG2JacExtDoubleNegMixed(b *testing.B) {
	var a g2JacExtended
	a.doubleMixed(&g2GenAff)

	var c G2Affine
	c.FromJacobian(&g2Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.doubleNegMixed(&c)
	}
}

func BenchmarkG2JacExtAdd(b *testing.B) {
	var a, c g2JacExtended
	a.doubleMixed(&g2GenAff)
	c.double(&a)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.add(&c)
	}
}

func BenchmarkG2JacExtDouble(b *testing.B) {
	var a g2JacExtended
	a.doubleMixed(&g2GenAff)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.double(&a)
	}
}

func fuzzJacobianG2Affine(p *G2Jac, f fp.Element) G2Jac {
	var res G2Jac
	res.X.Mul(&p.X, &f).Mul(&res.X, &f)
	res.Y.Mul(&p.Y, &f).Mul(&res.Y, &f).Mul(&res.Y, &f)
	res.Z.Mul(&p.Z, &f)
	return res
}

func fuzzExtendedJacobianG2Affine(p *g2JacExtended, f fp.Element) g2JacExtended {
	var res g2JacExtended
	var ff, fff fp.Element
	ff.Square(&f)
	fff.Mul(&ff, &f)
	res.X.Mul(&p.X, &ff)
	res.Y.Mul(&p.Y, &fff)
	res.ZZ.Mul(&p.ZZ, &ff)
	res.ZZZ.Mul(&p.ZZZ, &fff)
	return res
}
