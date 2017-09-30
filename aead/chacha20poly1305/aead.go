package chacha20poly1305

import (
	"encoding/binary"

	"git.schwanenlied.me/yawning/chacha20.git"
	"git.schwanenlied.me/yawning/poly1305.git"
	"go.artemisc.eu/godium"
)

const (
	KeyBytes  = 32
	NSecBytes = 0
	NPubBytes = 8
	ABytes    = 16
)

type aeadImpl struct {
	key []byte
}

func New(key []byte) (impl godium.AEAD) {
	impl = &aeadImpl{
		key: key,
	}
	return
}

func (a *aeadImpl) Wipe() {
	godium.Wipe(a.key)
}

// Seal
func (a *aeadImpl) Seal(dst, nonce, plain, ad []byte) (cipher []byte) {
	block0 := make([]byte, 64)
	slen := make([]byte, 8)

	// get poly key
	ciph, _ := chacha20.NewCipher(a.key, nonce)
	ciph.KeyStream(block0)

	// create poly
	poly, _ := poly1305.New(block0[:poly1305.KeySize])
	godium.Wipe(block0)

	// update tag
	_, _ = poly.Write(ad)
	binary.LittleEndian.PutUint64(slen, uint64(len(ad)))
	_, _ = poly.Write(slen)

	// encrypt with xor
	cipher = append(dst, plain...)
	_ = ciph.Seek(1)
	ciph.XORKeyStream(cipher, cipher)

	// update tag
	_, _ = poly.Write(cipher)
	binary.LittleEndian.PutUint64(slen, uint64(len(cipher)))
	_, _ = poly.Write(slen)

	// add tag
	cipher = poly.Sum(cipher)

	// clear state
	ciph.Reset()
	poly.Clear()

	return
}

// Open
func (a *aeadImpl) Open(dst, nonce, cipher, ad []byte) (plain []byte, err error) {
	return
}

func (a *aeadImpl) Overhead() (c int) {
	c = ABytes
	return
}

func (a *aeadImpl) NonceSize() (c int) {
	c = NPubBytes
	return
}

func (a *aeadImpl) KeyBytes() (c int) {
	c = KeyBytes
	return
}

func (a *aeadImpl) NSecBytes() (c int) {
	c = NSecBytes
	return
}

func (a *aeadImpl) NPubBytes() (c int) {
	c = NPubBytes
	return
}

func (a *aeadImpl) ABytes() (c int) {
	c = ABytes
	return
}
