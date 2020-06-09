// Code generated by internal/pairing DO NOT EDIT
package bw6_761

import (
	"testing"

	"github.com/consensys/gurvy/bw6_761/fp"
	"github.com/consensys/gurvy/bw6_761/fr"
)

func TestPairingLineEval(t *testing.T) {
	t.Skip()
	G := G2Jac{}
	G.X.SetString("number")
	G.Y.SetString("number")
	G.Z.SetString("1")

	H := G2Jac{}
	H.X.SetString("number")
	H.Y.SetString("number")
	H.Z.SetString("1")

	var a, b, c fp.Element
	a.SetString("2903903751748121992039561169443592957526674295618607189912579965473824470836812596944859552608502931201741951820932")
	b.SetString("1774816561618860752500414710493623341591800940439525170025341193002058157919964682036775870087093646544275868761902")
	c.SetString("1")
	P := G1Jac{}
	P.X = a
	P.Y = b
	P.Z = c

	var Paff G1Affine
	P.ToAffineFromJac(&Paff)

	lRes := &lineEvalRes{}
	lineEvalJac(G, H, &Paff, lRes)

	var expectedA, expectedB, expectedC G2CoordType
	expectedA.SetString("number")
	expectedB.SetString("number")
	expectedC.SetString("number")

	if !lRes.r1.Equal(&expectedA) {
		t.Fatal("Error A coeff")
	}
	if !lRes.r0.Equal(&expectedB) {
		t.Fatal("Error A coeff")
	}
	if !lRes.r2.Equal(&expectedC) {
		t.Fatal("Error A coeff")
	}
}

func TestMagicPairing(t *testing.T) {

	var r1, r2 PairingResult

	r1.SetRandom()
	r2.SetRandom()

	t.Log(r1)
	t.Log(r2)

	curve := BW6_761()

	res1 := curve.FinalExponentiation(&r1)
	res2 := curve.FinalExponentiation(&r2)

	if res1.Equal(&res2) {
		t.Fatal("TestMagicPairing failed")
	}
}

func TestComputePairing(t *testing.T) {

	curve := BW6_761()

	G := curve.g2Gen.Clone()
	P := curve.g1Gen.Clone()
	sG := G.Clone()
	sP := P.Clone()

	var Gaff, sGaff G2Affine
	var Paff, sPaff G1Affine

	// checking bilinearity

	// check 1
	scalar := fr.Element{123}
	sG.ScalarMul(curve, sG, scalar)
	sP.ScalarMul(curve, sP, scalar)

	var mRes1, mRes2, mRes3 PairingResult

	P.ToAffineFromJac(&Paff)
	sP.ToAffineFromJac(&sPaff)
	G.ToAffineFromJac(&Gaff)
	sG.ToAffineFromJac(&sGaff)

	res1 := curve.FinalExponentiation(curve.MillerLoop(Paff, sGaff, &mRes1))
	res2 := curve.FinalExponentiation(curve.MillerLoop(sPaff, Gaff, &mRes2))

	if !res1.Equal(&res2) {
		t.Fatal("pairing failed")
	}

	// check 2
	s1G := G.Clone()
	s2G := G.Clone()
	s3G := G.Clone()
	s1 := fr.Element{29372983}
	s2 := fr.Element{209302420904}
	var s3 fr.Element
	s3.Add(&s1, &s2)
	s1G.ScalarMul(curve, s1G, s1)
	s2G.ScalarMul(curve, s2G, s2)
	s3G.ScalarMul(curve, s3G, s3)

	var s1Gaff, s2Gaff, s3Gaff G2Affine
	s1G.ToAffineFromJac(&s1Gaff)
	s2G.ToAffineFromJac(&s2Gaff)
	s3G.ToAffineFromJac(&s3Gaff)

	rs1 := curve.FinalExponentiation(curve.MillerLoop(Paff, s1Gaff, &mRes1))
	rs2 := curve.FinalExponentiation(curve.MillerLoop(Paff, s2Gaff, &mRes2))
	rs3 := curve.FinalExponentiation(curve.MillerLoop(Paff, s3Gaff, &mRes3))
	rs1.Mul(&rs2, &rs1)
	if !rs3.Equal(&rs1) {
		t.Fatal("pairing failed2")
	}

}

//--------------------//
//     benches		  //
//--------------------//

func BenchmarkLineEval(b *testing.B) {

	curve := BW6_761()

	H := G2Jac{}
	H.ScalarMul(curve, &curve.g2Gen, fr.Element{1213})

	lRes := &lineEvalRes{}
	var g1GenAff G1Affine
	curve.g1Gen.ToAffineFromJac(&g1GenAff)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lineEvalJac(curve.g2Gen, H, &g1GenAff, lRes)
	}

}

func BenchmarkPairing(b *testing.B) {

	curve := BW6_761()

	var mRes PairingResult

	var g1GenAff G1Affine
	var g2GenAff G2Affine

	curve.g1Gen.ToAffineFromJac(&g1GenAff)
	curve.g2Gen.ToAffineFromJac(&g2GenAff)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		curve.FinalExponentiation(curve.MillerLoop(g1GenAff, g2GenAff, &mRes))
	}
}

func BenchmarkFinalExponentiation(b *testing.B) {

	var a PairingResult

	curve := BW6_761()

	a.SetString(
		"1382424129690940106527336948935335363935127549146605398842626667204683483408227749",
		"0121296909401065273369489353353639351275491466053988426266672046834834082277499690",
		"7336948129690940106527336948935335363935127549146605398842626667204683483408227749",
		"6393512129690940106527336948935335363935127549146605398842626667204683483408227749",
		"2581296909401065273369489353353639351275491466053988426266672046834834082277496644",
		"5331296909401065273369489353353639351275491466053988426266672046834834082277495363")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		curve.FinalExponentiation(&a)
	}

}
