package crypto

import (
	"fmt"

	"../storage"
)

const (
	hashPasswordDescriptor = ".password_hash"
	encodedKeyDescriptor   = ".encoded_key"
	saltDescriptor         = ".salt"
	contentExtension       = ".dat"
)

type CryptoStorage interface {
	Init(password []byte) error
	SaveContent(password []byte, code string, content []byte) error
	GetContent(password []byte, code string) ([]byte, error)
}

type persistentCryptoProvider struct {
	storage      storage.Storage
	passwordHash HashedPassword
	encodedKey   []byte
	salt         []byte
}

func (cryptoStorage *persistentCryptoProvider) Init(password []byte) error {
	storedPasswordHash, errStoredPasswordHash := cryptoStorage.storage.Read(hashPasswordDescriptor)
	if errStoredPasswordHash == nil {
		errHashVerification := verifyPasswordHash(HashedPassword(storedPasswordHash), password)
		if errHashVerification != nil {
			return fmt.Errorf("Password hash verification failed")
		}
	}
	storedEncodedKey, errStoredEncodedKey := cryptoStorage.storage.Read(encodedKeyDescriptor)
	storedSalt, errStoredSalt := cryptoStorage.storage.Read(saltDescriptor)
	cryptoStorage.validateCryptoStorageInitState(errStoredEncodedKey, errStoredSalt, errStoredPasswordHash)

	if errStoredEncodedKey != nil {
		hashedPassword, encodedKey, salt, err := cryptoStorage.generateNewKeys(password)
		if err != nil {
			return fmt.Errorf("Unable to generate new keys")
		}
		cryptoStorage.storeNewKeys(hashedPassword, encodedKey, salt)
		cryptoStorage.setInternalCryptoSystem(hashedPassword, encodedKey, salt)
	} else {
		cryptoStorage.setInternalCryptoSystem(HashedPassword(storedPasswordHash), storedEncodedKey, storedSalt)
	}
	return nil
}

func (cryptoStorage *persistentCryptoProvider) SaveContent(password []byte, code string, content []byte) error {
	if cryptoStorage.encodedKey == nil || cryptoStorage.passwordHash == nil || cryptoStorage.salt == nil {
		panic("Crypto storage was not initialized")
	}
	errHashVerification := verifyPasswordHash(cryptoStorage.passwordHash, password)
	if errHashVerification != nil {
		return fmt.Errorf("Password hash verification failed")
	}
	secureKey, err := cryptoStorage.decryptEncodedKey(password)
	if err != nil {
		return err
	}
	encryptedContent, err := encrypt(content, secureKey)
	if err != nil {
		return err
	}
	return cryptoStorage.storage.Save(cryptoStorage.wrapCode(code), encryptedContent)
}

func (cryptoStorage *persistentCryptoProvider) GetContent(password []byte, code string) ([]byte, error) {
	cryptoStorage.validateCryptoStorageInitialized()
	errHashVerification := verifyPasswordHash(cryptoStorage.passwordHash, password)
	if errHashVerification != nil {
		return nil, fmt.Errorf("Password hash verification failed")
	}
	secureKey, err := cryptoStorage.decryptEncodedKey(password)
	if err != nil {
		return nil, err
	}
	encryptedContent, err := cryptoStorage.storage.Read(cryptoStorage.wrapCode(code))
	if err != nil {
		return nil, err
	}

	content, err := decrypt(encryptedContent, secureKey)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (cryptoStorage *persistentCryptoProvider) generateNewKeys(password []byte) (HashedPassword, []byte, []byte, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, nil, nil, err
	}

	secureKey := newEncryptionKey()
	salt := newNonce()
	passwordKey := getKeyFromPassword(password, salt)
	encodedKey, err := encrypt(secureKey[:], passwordKey)
	if err != nil {
		return nil, nil, nil, err
	}
	return hashedPassword, encodedKey, salt, nil
}

func (cryptoStorage *persistentCryptoProvider) storeNewKeys(hashedPassword HashedPassword, encodedKey []byte, salt []byte) {
	errSaveHashedPassword := cryptoStorage.storage.Save(hashPasswordDescriptor, []byte(hashedPassword))
	errSaveSalt := cryptoStorage.storage.Save(saltDescriptor, salt)
	errSaveEncodedKey := cryptoStorage.storage.Save(encodedKeyDescriptor, encodedKey)
	if errSaveHashedPassword != nil || errSaveSalt != nil || errSaveEncodedKey != nil {
		cryptoStorage.storage.Remove(hashPasswordDescriptor)
		cryptoStorage.storage.Remove(saltDescriptor)
		cryptoStorage.storage.Remove(encodedKeyDescriptor)
		panic("Unable to store newly generated keys")
	}
}

func (cryptoStorage *persistentCryptoProvider) validateCryptoStorageInitState(errStoredEncodedKey error, errStoredSalt error, errStoredPasswordHash error) {
	if (errStoredEncodedKey == nil && errStoredSalt != nil) || (errStoredEncodedKey != nil && errStoredSalt == nil) {
		panic("Salt and encoded key are desynchonized!")
	} else if errStoredEncodedKey == nil && errStoredPasswordHash != nil {
		panic("Password is lost!")
	} else if errStoredEncodedKey != nil && errStoredPasswordHash == nil {
		panic("Encoded key is lost!")
	}
}

func (cryptoStorage *persistentCryptoProvider) decryptEncodedKey(password []byte) (*SecureKey, error) {
	passwordKey := getKeyFromPassword(password, cryptoStorage.salt)
	secureKeySlice, err := decrypt(cryptoStorage.encodedKey, passwordKey)
	if err != nil {
		return nil, err
	}
	var secureKey SecureKey
	copy(secureKey[:], secureKeySlice[:32])
	return &secureKey, nil
}

func (cryptoStorage *persistentCryptoProvider) setInternalCryptoSystem(hashedPassword HashedPassword, encodedKey []byte, salt []byte) {
	cryptoStorage.passwordHash = hashedPassword
	cryptoStorage.encodedKey = encodedKey
	cryptoStorage.salt = salt
}

func (cryptoStorage *persistentCryptoProvider) validateCryptoStorageInitialized() {
	if cryptoStorage.encodedKey == nil || cryptoStorage.passwordHash == nil || cryptoStorage.salt == nil {
		panic("Crypto storage was not initialized")
	}
}

func (cryptoStorage *persistentCryptoProvider) wrapCode(code string) string {
	return code + contentExtension
}

func NewPersistentCryptoProvider(storage storage.Storage) CryptoStorage {
	return &persistentCryptoProvider{
		storage,
		nil,
		nil,
		nil,
	}
}
