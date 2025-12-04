// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/zewi/db"
	"github.com/clivern/zewi/service"

	"github.com/rs/zerolog/log"
)

const stateKey = "state"

// GetStateAction handles GET requests to retrieve the state
func GetStateAction(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("Get state endpoint called")

	database := db.GetDB()
	if database == nil {
		log.Error().Msg("Database not initialized")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Database not initialized",
		})
		return
	}

	repo := db.NewOptionRepository(database)
	option, err := repo.Get(stateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get state from database")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve state",
		})
		return
	}

	if option == nil {
		// State doesn't exist yet, return empty/null value
		service.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"state": nil,
		})
		return
	}

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"state": option.Value,
	})
}

// UpdateStateRequest represents the request body for updating state
type UpdateStateRequest struct {
	Value string `json:"value"`
}

// UpdateStateAction handles PUT/POST requests to update the state
func UpdateStateAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Update state endpoint called")

	database := db.GetDB()
	if database == nil {
		log.Error().Msg("Database not initialized")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Database not initialized",
		})
		return
	}

	var req UpdateStateRequest
	if err := service.DecodeJSON(r, &req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		service.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	repo := db.NewOptionRepository(database)

	// Check if state exists
	existing, err := repo.Get(stateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check existing state")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to check state",
		})
		return
	}

	if existing == nil {
		// Create new state option
		if err := repo.Create(stateKey, req.Value); err != nil {
			log.Error().Err(err).Msg("Failed to create state")
			service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to create state",
			})
			return
		}
		log.Info().Msg("State created successfully")
	} else {
		// Update existing state
		if err := repo.Update(stateKey, req.Value); err != nil {
			log.Error().Err(err).Msg("Failed to update state")
			service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to update state",
			})
			return
		}
		log.Info().Msg("State updated successfully")
	}

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "State updated successfully",
		"state":   req.Value,
	})
}
