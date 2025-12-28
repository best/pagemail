package crypto

import (
	"bytes"
	"testing"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"valid 32-byte key", "12345678901234567890123456789012", false},
		{"short key (will be padded)", "short", false},
		{"long key (will be truncated)", "1234567890123456789012345678901234567890", false},
		{"base64 key", "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoxMjM0NTY=", false},
		{"empty key", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncryptor(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && enc == nil {
				t.Error("NewEncryptor() returned nil encryptor")
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	enc, err := NewEncryptor("test-encryption-key-for-testing!")
	if err != nil {
		t.Fatalf("NewEncryptor() error = %v", err)
	}

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"empty", []byte{}},
		{"short text", []byte("hello")},
		{"long text", []byte("this is a much longer piece of text that should still encrypt and decrypt properly")},
		{"binary data", []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0xfd}},
		{"unicode", []byte("你好世界🌍")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := enc.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			if len(tt.plaintext) > 0 && bytes.Equal(ciphertext, tt.plaintext) {
				t.Error("Ciphertext equals plaintext")
			}

			decrypted, err := enc.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptString(t *testing.T) {
	enc, _ := NewEncryptor("test-key-32-bytes-long-exactly!!")

	original := "password123"
	ciphertext, err := enc.EncryptString(original)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}

	decrypted, err := enc.DecryptToString(ciphertext)
	if err != nil {
		t.Fatalf("DecryptToString() error = %v", err)
	}

	if decrypted != original {
		t.Errorf("DecryptToString() = %v, want %v", decrypted, original)
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	enc, _ := NewEncryptor("test-key")

	tests := []struct {
		name       string
		ciphertext []byte
	}{
		{"too short", []byte{0x01, 0x02}},
		{"empty", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enc.Decrypt(tt.ciphertext)
			if err == nil {
				t.Error("Decrypt() expected error for invalid ciphertext")
			}
		})
	}
}

func TestDifferentKeysProduceDifferentCiphertext(t *testing.T) {
	enc1, _ := NewEncryptor("key-one-32-bytes-long-exactly!!!")
	enc2, _ := NewEncryptor("key-two-32-bytes-long-exactly!!!")

	plaintext := []byte("same plaintext")

	ct1, _ := enc1.Encrypt(plaintext)
	ct2, _ := enc2.Encrypt(plaintext)

	if bytes.Equal(ct1, ct2) {
		t.Error("Different keys should produce different ciphertexts")
	}

	_, err := enc2.Decrypt(ct1)
	if err == nil {
		t.Error("Decrypting with wrong key should fail")
	}
}
