package verify

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/rianby64/tcp-pow/models"
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
	validator       Validator
	keySize         int
	maxLengthKeyMsg int
}

func (handler *Handler) Handler(conn net.Conn) error {
	puzzle := GenerateKey(handler.keySize)

	if n, err := conn.Write(puzzle); err != nil {
		return errors.Wrap(err, "conn.Write(key)")
	} else if n != len(puzzle) {
		// I decided to trigger this error because if the client doesn't receive the whole puzzle
		// then chances are low the client will send something correct
		return errors.Wrapf(models.ErrIncorrect, "n(%d) != %d", n, len(puzzle))
	}

	resolvedPuzzle := make([]byte, handler.maxLengthKeyMsg)
	if n, err := conn.Read(resolvedPuzzle); err != nil {
		return errors.Wrap(err, "conn.Read(resolvedPuzzle)")
	} else {
		resolvedPuzzle = resolvedPuzzle[:n]
	}

	if len(resolvedPuzzle) == 0 {
		return errors.Wrap(models.ErrEmpty, "len(resolvedPuzzle) == 0")
	}

	hasSolvedPuzzleOurPuzzle := bytes.Contains(resolvedPuzzle, puzzle)
	isSolvedPuzzleCorrect := handler.validator.Check(string(resolvedPuzzle))

	if hasSolvedPuzzleOurPuzzle && isSolvedPuzzleCorrect {
		return nil
	}

	return errors.Wrapf(models.ErrIncorrect,
		"hasSolvedPuzzleOurPuzzle(%t) && isSolvedPuzzleCorrect(%t)",
		hasSolvedPuzzleOurPuzzle, isSolvedPuzzleCorrect,
	)
}

func New(validator Validator, keySize int, leadingBits, saltSize uint) *Handler {
	rand.Seed(time.Now().Unix())

	difficultySize := len(fmt.Sprint(leadingBits))
	extraSize := 0
	versionSize := len("1")
	separatorSize := len(":")
	dateFormatSize := len("060102")
	nonceSize := len(fmt.Sprintf("%x", math.MaxInt))

	maxLengthKeyMsg := versionSize + separatorSize +
		difficultySize + separatorSize +
		dateFormatSize + separatorSize +
		keySize + separatorSize +
		extraSize + separatorSize +
		int(saltSize) + separatorSize +
		nonceSize

	return &Handler{
		validator:       validator,
		keySize:         keySize,
		maxLengthKeyMsg: maxLengthKeyMsg,
	}
}
