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

package polynomial

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr/fft"
)

// Polynomial polynomial represented by coefficients bn254 fr field.
type Polynomial []fr.Element

// NewPolynomial
func NewPolynomial(size uint64) Polynomial {
	s := ecc.NextPowerOfTwo(size)
	res := make(Polynomial, size, s)
	return res
}

// Eval evaluates p at v
// returns a fr.Element
func (p *Polynomial) Eval(v *fr.Element) fr.Element {

	_p := *p
	res := _p[len(_p)-1]
	for i := len(_p) - 2; i >= 0; i-- {
		res.Mul(&res, v)
		res.Add(&res, &_p[i])
	}

	return res
}

// Copy returns a copy of the polynomial
func (p *Polynomial) Copy() *Polynomial {
	_p := make(Polynomial, len(*p))
	copy(_p, *p)
	return &_p
}

// sort returns pi, pj where len(pi)>=len(pj).
func sort(p1, p2 *Polynomial) (Polynomial, Polynomial) {
	s1 := len(*p1)
	s2 := len(*p2)
	if s2 > s1 {
		return *p2, *p1
	}
	return *p1, *p2
}

// resize adapt the size of p to max(len(p1),len(p2),len(p))
func (p *Polynomial) resize(p1, p2 *Polynomial) {
	if len(*p) > len(*p1) && len(*p) > len(*p2) {
		return
	}
	_p1, _ := sort(p1, p2)
	pad(p, len(_p1)-len(*p))
}

// Degree returns the degree of p
func (p *Polynomial) Degree() uint64 {
	_p := *p
	res := len(_p) - 1
	for i := len(_p) - 1; i >= 0; i-- {
		if !_p[i].IsZero() {
			return uint64(res)
		}
		res--
	}
	return uint64(res)
}

// Equal checks equality between two polynomials
func (p *Polynomial) Equal(p1 *Polynomial) bool {

	res := true
	_p1, _p2 := sort(p, p1)
	s1 := len(_p1)
	s2 := len(_p2)
	d := s1 - s2
	for i := 0; i < d; i++ {
		res = res && _p1[s1-1-i].IsZero()
	}
	for i := 0; i < s2; i++ {
		res = res && _p1[s2-1-i].Equal(&_p2[s2-1-i])
	}

	return res
}

// addSorted  sets p to p1+p2, assuming len(p)>=len(p1)>=len(p2)
func addSorted(p, p1, p2 *Polynomial) {
	i := 0
	for ; i < len(*p2); i++ {
		(*p)[i].Add(&(*p1)[i], &(*p2)[i])
	}
	for ; i < len(*p1); i++ {
		(*p)[i].Add(&(*p1)[i], &(*p2)[i])
	}
	for ; i < len(*p); i++ {
		(*p)[i].SetZero()
	}
}

// subSorted sets p to p1-p2, assuming len(p)>=len(p1)>=len(p2)
func subSorted(p, p1, p2 *Polynomial) {
	i := 0
	for ; i < len(*p2); i++ {
		(*p)[i].Sub(&(*p1)[i], &(*p2)[i])
	}
	for ; i < len(*p1); i++ {
		(*p)[i].Sub(&(*p1)[i], &(*p2)[i])
	}
	for ; i < len(*p); i++ {
		(*p)[i].SetZero()
	}
}

// Add adds p1 to p, and return p
func (p *Polynomial) Add(p1, p2 *Polynomial) *Polynomial {
	_p1, _p2 := sort(p1, p2)
	p.resize(p1, p2)
	addSorted(p, &_p1, &_p2)
	return p
}

// Sub adds p1 to p, and return p
func (p *Polynomial) Sub(p1, p2 *Polynomial) *Polynomial {
	_p1, _p2 := sort(p1, p2)
	p.resize(p1, p2)
	subSorted(p, &_p1, &_p2)
	return p
}

// Neg sets p to -p1
func (p *Polynomial) Neg(p1 *Polynomial) *Polynomial {
	p.resize(p1, p)
	*p = *p1
	for i := 0; i < len(*p); i++ {
		(*p)[i].Neg(&(*p)[i])
	}
	return p
}

// SetZero sets p=zero(p1)
func (p *Polynomial) SetZero() *Polynomial {
	for i := 0; i < len(*p); i++ {
		(*p)[i].SetZero()
	}
	return p
}

// pad pads p with i zeros
func pad(p *Polynomial, i int) {
	pad := make(Polynomial, i)
	*p = append(*p, pad...)
}

// mul sets p1=p2*p3, assuming all polynomials and the domain
// have the same size.
func mul(p1, p2, p3 *Polynomial) {
	p1.SetZero()
	d := fft.NewDomain(uint64(len(*p1)), 0, false)
	d.FFT(*p2, fft.DIF, 0)
	d.FFT(*p3, fft.DIF, 0)
	for i := 0; i < len(*p1); i++ {
		(*p1)[i].Mul(&(*p2)[i], &(*p3)[i])
	}
	d.FFTInverse(*p1, fft.DIT, 0)
}

// Mul sets p to p1*p2
func (p *Polynomial) Mul(p1, p2 *Polynomial) *Polynomial {

	_p1, _p2 := sort(p1, p2)
	s := len(_p1)
	n := ecc.NextPowerOfTwo(uint64(s))
	d1 := _p1.Degree()
	d2 := _p2.Degree()
	if d1+d2 < n {
		cp1 := make(Polynomial, n)
		cp2 := make(Polynomial, n)
		copy(cp1, _p1)
		copy(cp2, _p2)
		p.resize(&cp1, &cp2)
		mul(p, &cp1, &cp2)
		return p
	}
	cp1 := make(Polynomial, 2*n)
	cp2 := make(Polynomial, 2*n)
	copy(cp1, _p1)
	copy(cp2, _p2)
	p.resize(&cp1, &cp2)
	mul(p, &cp1, &cp2)
	return p
}
