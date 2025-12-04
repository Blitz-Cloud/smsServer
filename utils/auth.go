package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Argon2 parameters (OWASP recommended values for general web applications)
const (
	saltLength  = 16    // length of salt in bytes
	keyLength   = 32    // length of derived key (e.g., for AES-256)
	timeCost    = 1     // number of iterations
	memoryCost  = 64 * 1024 // memory cost in KiB (~64MB)
	parallelism = 4     // number of parallel threads
)

// generateSalt creates a cryptographically secure random salt
func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// HashPassword hashes a plaintext password using Argon2id
// It returns the full Argon2 hash string in the PHC string format
func HashPassword(password string) (string, error) {
	salt, err := generateSalt(saltLength)
	if err != nil {
		return "", err
	}

	hashedKey := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, keyLength)

	// Encode salt and hashed key to base64 for storage and PHC string format
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64HashedKey := base64.RawStdEncoding.EncodeToString(hashedKey)

	// Format as PHC string: $argon2id$v=19$m=<memory>,t=<time>,p=<threads>$<salt>$<hash>
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memoryCost, timeCost, parallelism, b64Salt, b64HashedKey), nil
}

// VerifyPassword compares a plaintext password with an Argon2 hash string
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash string to extract parameters, salt, and the stored hash
	var (
		version     int
		mem, time   uint32
		threads     uint8
		b64Salt     string
		b64HashedKey string
	)

	_, err := fmt.Sscanf(encodedHash, "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		&version, &mem, &time, &threads, &b64Salt, &b64HashedKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse Argon2 hash: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(b64HashedKey)
	if err != nil {
		return false, fmt.Errorf("failed to decode hashed key: %w", err)
	}

	// Re-hash the provided password with the extracted parameters and salt
	// Use the original keyLength when calling IDKey
	computedHash := argon2.IDKey([]byte(password), salt, time, mem, threads, uint32(len(storedHash)))

	// Use constant-time comparison to prevent timing attacks
	match := subtle.ConstantTimeCompare(storedHash, computedHash) == 1
	return match, nil
}