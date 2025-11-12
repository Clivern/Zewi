// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// Gravatar provides methods for generating Gravatar URLs
type Gravatar struct{}

// GetGravatar generates a Gravatar URL for the given email address.
func (g *Gravatar) GetGravatar(email string, size int) string {
	if size <= 0 {
		size = 200
	}

	emailLower := strings.ToLower(email)
	emailBytes := []byte(emailLower)
	hash := sha256.Sum256(emailBytes)
	hashHex := hex.EncodeToString(hash[:])

	return "https://gravatar.com/avatar/" + hashHex + "?s=" + strconv.Itoa(size)
}
