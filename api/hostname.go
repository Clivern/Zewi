// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/clivern/zewi/service"

	"github.com/rs/zerolog/log"
)

// HostnameAction handles requests to return the server hostname
func HostnameAction(w http.ResponseWriter, _ *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get hostname")
		hostname = "unknown"
	}

	log.Debug().Str("hostname", hostname).Msg("Hostname endpoint called")

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"hostname": strings.ToLower(hostname),
	})
}
