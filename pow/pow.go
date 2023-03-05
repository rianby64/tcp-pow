package pow

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
