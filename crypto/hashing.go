package crypto

import "golang.org/x/crypto/bcrypt"

type HashedPassword []byte

func hashPassword(password []byte) (HashedPassword, error) {
	rawHashedPassword, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return nil, err
	}
	return HashedPassword(rawHashedPassword), nil
}

func verifyPasswordHash(hashedPassword HashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
}
