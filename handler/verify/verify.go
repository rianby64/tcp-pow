package verify

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/pkg/errors"
)

const (
	alphabet     = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	alphabetSize = len(alphabet)
)

func GenerateKey(size int) []byte {
	key := make([]byte, size)
	alphabetInBytes := []byte(alphabet)

	for index := 0; index < size; index++ {
		randomIndex := rand.Intn(alphabetSize)
		key[index] = alphabetInBytes[randomIndex]
	}

	return key
}

type Validator interface {
	Check(stamp string) bool
}

type Handler struct {
	validator Validator
}

func (handler *Handler) Handler(conn net.Conn) error {
	puzzle := GenerateKey(64)

	if n, err := conn.Write(puzzle); err != nil {
		return errors.Wrap(err, "conn.Write(key)")
	} else if n != 64 {
		// I decided to trigger this error as if the client doesn't receive the whole puzzle
		// then chances are low the client will send something correct
		return fmt.Errorf("n != 64")
	}

	resolvedPuzzle := make([]byte, 104)
	if n, err := conn.Read(resolvedPuzzle); err != nil {
		return errors.Wrap(err, "conn.Read(resolvedPuzzle)")
	} else {
		resolvedPuzzle = resolvedPuzzle[:n]
	}

	if len(resolvedPuzzle) == 0 {
		return fmt.Errorf("len(resolvedPuzzle) == 0")
	}

	hasSolvedPuzzleOurPuzzle := bytes.Contains(resolvedPuzzle, puzzle)
	isSolvedPuzzleCorrect := handler.validator.Check(string(resolvedPuzzle))

	if hasSolvedPuzzleOurPuzzle && isSolvedPuzzleCorrect {
		return nil
	}

	return fmt.Errorf("incorrect puzzle")
}

func New(validator Validator) *Handler {
	rand.Seed(time.Now().Unix())

	return &Handler{
		validator: validator,
	}
}
