package pow

import (
	"github.com/catalinc/hashcash"
	"github.com/pkg/errors"
)

type Validator interface {
	Check(stamp string) bool
}

type POW struct {
	proof Validator
}

func (pow *POW) Check(stamp string) bool {
	return pow.proof.Check(stamp)
}

func New(proof Validator) *POW {
	return &POW{
		proof: proof,
	}
}

func Generate(data string) (string, error) {
	leadingBits := uint(20)
	saltLength := uint(8)
	extra := ""

	h := hashcash.New(leadingBits, saltLength, extra)

	// Mint a new stamp
	stamp, err := h.Mint(data)
	if err != nil {
		return "", errors.Wrap(err, "stamp, err := h.Mint(string(data))")
	}

	return stamp, nil
}
