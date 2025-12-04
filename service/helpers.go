// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteJSON writes a JSON response with the given status code and data
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("failed to write JSON response: %w", err)
	}
	return nil
}

// DecodeJSON decodes JSON from request body
func DecodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	return nil
}
