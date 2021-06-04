// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vault

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/argon2"
)

const (
	saltLength         = 16
	nonceLength        = 24
	keyLength          = 32
	shamirSecretLength = 32
)

func deriveKey(passphrase string, salt []byte) [keyLength]byte {
	secretKeyBytes := argon2.IDKey([]byte(passphrase), salt, 4, 64*1024, 4, 32)
	var secretKey [keyLength]byte
	copy(secretKey[:], secretKeyBytes)
	return secretKey
}


type PassphraseBoxer struct {
	passphrase string
}

func NewPassphraseBoxer(password string) *PassphraseBoxer {
	return &PassphraseBoxer{
		passphrase: password,
	}
}

func (b *PassphraseBoxer) WrapType() string {
	return "passphrase"
}

func (b *PassphraseBoxer) Seal(in []byte) (string, error) {
	var nonce [nonceLength]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		return "", err
	}

	salt := make([]byte, saltLength)
	if _, err := crypto_rand.Read(salt); err != nil {
		return "", err
	}
	secretKey := deriveKey(b.passphrase, salt)
	prefix := append(salt, nonce[:]...)

	cipherText := secretbox.Seal(prefix, in, &nonce, &secretKey)

	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

func (b *PassphraseBoxer) Open(in string) ([]byte, error) {
	buf, err := base64.RawStdEncoding.DecodeString(in)
	if err != nil {
		return []byte{}, err
	}

	salt := make([]byte, saltLength)
	copy(salt, buf[:saltLength])
	var nonce [nonceLength]byte
	copy(nonce[:], buf[saltLength:nonceLength+saltLength])

	secretKey := deriveKey(b.passphrase, salt)
	decrypted, ok := secretbox.Open(nil, buf[nonceLength+saltLength:], &nonce, &secretKey)
	if !ok {
		return []byte{}, fmt.Errorf("failed to decrypt")
	}
	return decrypted, nil
}
