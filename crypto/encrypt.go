package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

type SecureKey [32]byte

func newEncryptionKey() *SecureKey {
	key := SecureKey{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

func newNonce() []byte {
	nonce := [8]byte{}
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		panic(err)
	}
	return nonce[:]
}

func encrypt(plaintext []byte, key *SecureKey) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key *SecureKey) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("Malformed ciphertext was passed")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

func getKeyFromPassword(password []byte, salt []byte) *SecureKey {
	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha1.New)
	var aesKey SecureKey
	copy(aesKey[:], key[:32])
	return &aesKey
}
