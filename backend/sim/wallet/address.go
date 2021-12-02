// Copyright 2019 - See NOTICE file for copyright holders.
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

package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"math/big"

	"perun.network/go-perun/log"
	"perun.network/go-perun/wallet"
)

// Address represents a simulated address.
type Address ecdsa.PublicKey

const (
	// ElemLen is the length of the binary representation of a single element
	// of the address in bytes.
	ElemLen = 32

	// AddressBinaryLength is the length of the binary representation of
	// Address in bytes.
	AddressBinaryLen = 2 * ElemLen
)

// compile time check that we implement the perun Address interface.
var _ wallet.Address = (*Address)(nil)

// NewRandomAddress creates a new address using the randomness
// provided by rng.
func NewRandomAddress(rng io.Reader) *Address {
	privateKey, err := ecdsa.GenerateKey(curve, rng)
	if err != nil {
		log.Panicf("Creation of account failed with error", err)
	}

	return &Address{
		Curve: privateKey.Curve,
		X:     privateKey.X,
		Y:     privateKey.Y,
	}
}

// Bytes converts this address to bytes.
func (a *Address) Bytes() []byte {
	data := a.byteArray()
	return data[:]
}

// ByteArray converts an address into a 64-byte array. The returned array
// consists of two 32-byte chunks representing the public key's X and Y values.
func (a *Address) byteArray() (data [AddressBinaryLen]byte) {
	xb := a.X.Bytes()
	yb := a.Y.Bytes()

	// Left-pad with 0 bytes.
	copy(data[ElemLen-len(xb):ElemLen], xb)
	copy(data[AddressBinaryLen-len(yb):AddressBinaryLen], yb)

	return data
}

// String converts this address to a human-readable string.
func (a *Address) String() string {
	return fmt.Sprintf("0x%x", a.byteArray())
}

// Equal checks the equality of two addresses. The implementation must be
// equivalent to checking `Address.Cmp(Address) == 0`.
// Pancis if the passed address is of the wrong type.
func (a *Address) Equal(addr wallet.Address) bool {
	b, ok := addr.(*Address)
	if !ok {
		log.Panic("wrong address type")
	}
	return (a.X.Cmp(b.X) == 0) && (a.Y.Cmp(b.Y) == 0)
}

// Cmp checks the ordering of two addresses according to following definition:
//   -1 if (a.X <  addr.X) || ((a.X == addr.X) && (a.Y < addr.Y))
//    0 if (a.X == addr.X) && (a.Y == addr.Y)
//   +1 if (a.X >  addr.X) || ((a.X == addr.X) && (a.Y > addr.Y))
// So the X coordinate is weighted higher.
// Pancis if the passed address is of the wrong type.
func (a *Address) Cmp(addr wallet.Address) int {
	b, ok := addr.(*Address)
	if !ok {
		log.Panic("wrong address type")
	}
	const EQ = 0
	xCmp, yCmp := a.X.Cmp(b.X), a.Y.Cmp(b.Y)
	if xCmp != EQ {
		return xCmp
	}
	return yCmp
}

// MarshalBinary marshals the address into its binary representation.
// Error will always be nil, it is for implementing BinaryMarshaler.
func (a *Address) MarshalBinary() ([]byte, error) {
	data := a.byteArray()
	return data[:], nil
}

// UnmarshalBinary unmarshals the address from its binary representation.
func (a *Address) UnmarshalBinary(data []byte) error {
	if len(data) != AddressBinaryLen {
		return fmt.Errorf("unexpected address length %d, want %d", len(data), AddressBinaryLen) //nolint: goerr113
	}
	a.X = new(big.Int).SetBytes(data[:ElemLen])
	a.Y = new(big.Int).SetBytes(data[ElemLen:])
	a.Curve = curve

	return nil
}
