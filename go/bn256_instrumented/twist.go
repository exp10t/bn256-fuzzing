package main

import fuzz_helper "github.com/guidovranken/go-coverage-instrumentation/helper"

import (
	"math/big"
)

// twistPoint implements the elliptic curve y²=x³+3/ξ over GF(p²). Points are
// kept in Jacobian form and t=z² when valid. The group G₂ is the set of
// n-torsion points of this curve over GF(p²) (where n = Order)
type twistPoint struct {
	x, y, z, t *gfP2
}

var twistB = &gfP2{
	bigFromBase10("266929791119991161246907387137283842545076965332900288569378510910307636690"),
	bigFromBase10("19485874751759354771024239261021720505790618469301721065564631296452457478373"),
}

// twistGen is the generator of group G₂.
var twistGen = &twistPoint{
	&gfP2{
		bigFromBase10("11559732032986387107991004021392285783925812861821192530917403151452391805634"),
		bigFromBase10("10857046999023057135944570762232829481370756359578518086990519993285655852781"),
	},
	&gfP2{
		bigFromBase10("4082367875863433681332203403145435568316851327593401208105741076214120093531"),
		bigFromBase10("8495653923123431417604973247489272438418190587263600148770280649306958101930"),
	},
	&gfP2{
		bigFromBase10("0"),
		bigFromBase10("1"),
	},
	&gfP2{
		bigFromBase10("0"),
		bigFromBase10("1"),
	},
}

func newTwistPoint(pool *bnPool) *twistPoint {
	fuzz_helper.AddCoverage(22588)
	return &twistPoint{
		newGFp2(pool),
		newGFp2(pool),
		newGFp2(pool),
		newGFp2(pool),
	}
}

func (c *twistPoint) String() string {
	fuzz_helper.AddCoverage(44810)
	return "(" + c.x.String() + ", " + c.y.String() + ", " + c.z.String() + ")"
}

func (c *twistPoint) Put(pool *bnPool) {
	fuzz_helper.AddCoverage(5262)
	c.x.Put(pool)
	c.y.Put(pool)
	c.z.Put(pool)
	c.t.Put(pool)
}

func (c *twistPoint) Set(a *twistPoint) {
	fuzz_helper.AddCoverage(17878)
	c.x.Set(a.x)
	c.y.Set(a.y)
	c.z.Set(a.z)
	c.t.Set(a.t)
}

// IsOnCurve returns true iff c is on the curve where c must be in affine form.
func (c *twistPoint) IsOnCurve() bool {
	fuzz_helper.AddCoverage(45021)
	pool := new(bnPool)
	yy := newGFp2(pool).Square(c.y, pool)
	xxx := newGFp2(pool).Square(c.x, pool)
	xxx.Mul(xxx, c.x, pool)
	yy.Sub(yy, xxx)
	yy.Sub(yy, twistB)
	yy.Minimal()
	return yy.x.Sign() == 0 && yy.y.Sign() == 0
}

func (c *twistPoint) SetInfinity() {
	fuzz_helper.AddCoverage(39040)
	c.z.SetZero()
}

func (c *twistPoint) IsInfinity() bool {
	fuzz_helper.AddCoverage(2095)
	return c.z.IsZero()
}

func (c *twistPoint) Add(a, b *twistPoint, pool *bnPool) {
	fuzz_helper.AddCoverage(21668)

	if a.IsInfinity() {
		fuzz_helper.AddCoverage(42483)
		c.Set(b)
		return
	} else {
		fuzz_helper.AddCoverage(6577)
	}
	fuzz_helper.AddCoverage(45213)
	if b.IsInfinity() {
		fuzz_helper.AddCoverage(17393)
		c.Set(a)
		return
	} else {
		fuzz_helper.AddCoverage(64174)
	}
	fuzz_helper.AddCoverage(16619)

	z1z1 := newGFp2(pool).Square(a.z, pool)
	z2z2 := newGFp2(pool).Square(b.z, pool)
	u1 := newGFp2(pool).Mul(a.x, z2z2, pool)
	u2 := newGFp2(pool).Mul(b.x, z1z1, pool)

	t := newGFp2(pool).Mul(b.z, z2z2, pool)
	s1 := newGFp2(pool).Mul(a.y, t, pool)

	t.Mul(a.z, z1z1, pool)
	s2 := newGFp2(pool).Mul(b.y, t, pool)

	h := newGFp2(pool).Sub(u2, u1)
	xEqual := h.IsZero()

	t.Add(h, h)
	i := newGFp2(pool).Square(t, pool)
	j := newGFp2(pool).Mul(h, i, pool)

	t.Sub(s2, s1)
	yEqual := t.IsZero()
	if xEqual && yEqual {
		fuzz_helper.AddCoverage(38740)
		c.Double(a, pool)
		return
	} else {
		fuzz_helper.AddCoverage(35657)
	}
	fuzz_helper.AddCoverage(12692)
	r := newGFp2(pool).Add(t, t)

	v := newGFp2(pool).Mul(u1, i, pool)

	t4 := newGFp2(pool).Square(r, pool)
	t.Add(v, v)
	t6 := newGFp2(pool).Sub(t4, j)
	c.x.Sub(t6, t)

	t.Sub(v, c.x)
	t4.Mul(s1, j, pool)
	t6.Add(t4, t4)
	t4.Mul(r, t, pool)
	c.y.Sub(t4, t6)

	t.Add(a.z, b.z)
	t4.Square(t, pool)
	t.Sub(t4, z1z1)
	t4.Sub(t, z2z2)
	c.z.Mul(t4, h, pool)

	z1z1.Put(pool)
	z2z2.Put(pool)
	u1.Put(pool)
	u2.Put(pool)
	t.Put(pool)
	s1.Put(pool)
	s2.Put(pool)
	h.Put(pool)
	i.Put(pool)
	j.Put(pool)
	r.Put(pool)
	v.Put(pool)
	t4.Put(pool)
	t6.Put(pool)
}

func (c *twistPoint) Double(a *twistPoint, pool *bnPool) {
	fuzz_helper.AddCoverage(30358)

	A := newGFp2(pool).Square(a.x, pool)
	B := newGFp2(pool).Square(a.y, pool)
	C_ := newGFp2(pool).Square(B, pool)

	t := newGFp2(pool).Add(a.x, B)
	t2 := newGFp2(pool).Square(t, pool)
	t.Sub(t2, A)
	t2.Sub(t, C_)
	d := newGFp2(pool).Add(t2, t2)
	t.Add(A, A)
	e := newGFp2(pool).Add(t, A)
	f := newGFp2(pool).Square(e, pool)

	t.Add(d, d)
	c.x.Sub(f, t)

	t.Add(C_, C_)
	t2.Add(t, t)
	t.Add(t2, t2)
	c.y.Sub(d, c.x)
	t2.Mul(e, c.y, pool)
	c.y.Sub(t2, t)

	t.Mul(a.y, a.z, pool)
	c.z.Add(t, t)

	A.Put(pool)
	B.Put(pool)
	C_.Put(pool)
	t.Put(pool)
	t2.Put(pool)
	d.Put(pool)
	e.Put(pool)
	f.Put(pool)
}

func (c *twistPoint) Mul(a *twistPoint, scalar *big.Int, pool *bnPool) *twistPoint {
	fuzz_helper.AddCoverage(23294)
	sum := newTwistPoint(pool)
	sum.SetInfinity()
	t := newTwistPoint(pool)

	for i := scalar.BitLen(); i >= 0; i-- {
		fuzz_helper.AddCoverage(11162)
		t.Double(sum, pool)
		if scalar.Bit(i) != 0 {
			fuzz_helper.AddCoverage(49217)
			sum.Add(t, a, pool)
		} else {
			fuzz_helper.AddCoverage(34511)
			sum.Set(t)
		}
	}
	fuzz_helper.AddCoverage(61639)

	c.Set(sum)
	sum.Put(pool)
	t.Put(pool)
	return c
}

func (c *twistPoint) MakeAffine(pool *bnPool) *twistPoint {
	fuzz_helper.AddCoverage(64074)
	if c.z.IsOne() {
		fuzz_helper.AddCoverage(39226)
		return c
	} else {
		fuzz_helper.AddCoverage(2297)
	}
	fuzz_helper.AddCoverage(28614)

	zInv := newGFp2(pool).Invert(c.z, pool)
	t := newGFp2(pool).Mul(c.y, zInv, pool)
	zInv2 := newGFp2(pool).Square(zInv, pool)
	c.y.Mul(t, zInv2, pool)
	t.Mul(c.x, zInv2, pool)
	c.x.Set(t)
	c.z.SetOne()
	c.t.SetOne()

	zInv.Put(pool)
	t.Put(pool)
	zInv2.Put(pool)

	return c
}

func (c *twistPoint) Negative(a *twistPoint, pool *bnPool) {
	fuzz_helper.AddCoverage(40870)
	c.x.Set(a.x)
	c.y.SetZero()
	c.y.Sub(c.y, a.y)
	c.z.Set(a.z)
	c.t.SetZero()
}
