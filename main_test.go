package main

import (
	"reflect"
	"testing"

	bip39 "github.com/tyler-smith/go-bip39/wordlists"
)

func TestHashPassword(t *testing.T) {
	password := "my_secure_password"
	expectedLength := 32 // SHA-256 produces a 32-byte hash

	hash := hashPassword(password)
	if len(hash) != expectedLength {
		t.Errorf("Expected hash length %d, got %d", expectedLength, len(hash))
	}
}

func TestGenerateShiftValues(t *testing.T) {
	password := "my_secure_password"
	hash := hashPassword(password)
	expectedCount := 24

	shifts := generateShiftValues(hash)
	if len(shifts) != expectedCount {
		t.Errorf("Expected %d shift values, got %d", expectedCount, len(shifts))
	}

	for _, shift := range shifts {
		if shift < 0 || shift >= 2048 {
			t.Errorf("Shift value %d out of range", shift)
		}
	}
}

func TestShiftMnemonicWords(t *testing.T) {
	mnemonic := []string{"oppose", "duck", "hello", "neglect", "reveal", "key", "humor", "mosquito", "road", "evoke", "flock", "hedgehog"}
	shifts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	shiftedMnemonic := Encrypt(mnemonic, shifts)
	expectedMnemonic := make([]string, len(mnemonic))
	for i, word := range mnemonic {
		index := indexOf(word, bip39.English)
		expectedIndex := (index + shifts[i]) % len(bip39.English)
		expectedMnemonic[i] = bip39.English[expectedIndex]
	}

	if !reflect.DeepEqual(shiftedMnemonic, expectedMnemonic) {
		t.Errorf("Expected shifted mnemonic %v, got %v", expectedMnemonic, shiftedMnemonic)
	}
}

func TestReverseShiftMnemonicWords(t *testing.T) {
	mnemonic := []string{"oppose", "duck", "hello", "neglect", "reveal", "key", "humor", "mosquito", "road", "evoke", "flock", "hedgehog"}
	shifts := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	shiftedMnemonic := Encrypt(mnemonic, shifts)
	reversedMnemonic := Decrypt(shiftedMnemonic, shifts)

	if !reflect.DeepEqual(reversedMnemonic, mnemonic) {
		t.Errorf("Expected reversed mnemonic %v, got %v", mnemonic, reversedMnemonic)
	}
}

func TestReverseWithWrongPassword(t *testing.T) {
	originalMnemonic := []string{"oppose", "duck", "hello", "neglect", "reveal", "key", "humor", "mosquito", "road", "evoke", "flock", "hedgehog"}
	correctPassword := "my_secure_password"
	wrongPassword := "wrong_password"

	// Generate shift values using the correct password
	correctHash := hashPassword(correctPassword)
	correctShifts := generateShiftValues(correctHash)

	// Shift the mnemonic using the correct shifts
	shiftedMnemonic := Encrypt(originalMnemonic, correctShifts)

	// Attempt to reverse using the wrong password
	wrongHash := hashPassword(wrongPassword)
	wrongShifts := generateShiftValues(wrongHash)
	reversedMnemonic := Decrypt(shiftedMnemonic, wrongShifts)

	if reflect.DeepEqual(reversedMnemonic, originalMnemonic) {
		t.Errorf("Reversing with the wrong password should not match the original mnemonic. Got %v", reversedMnemonic)
	}
}
