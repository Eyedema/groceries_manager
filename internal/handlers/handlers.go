package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ubaldo/api_server/models"

	"github.com/go-chi/chi"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var db *gorm.DB
var itemCache *cache.Cache

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	result := db.Find(&items)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to retrieve items from the database")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid item ID")
		return
	}

	// Check if the item is already cached
	cachedItem, found := itemCache.Get(idStr)
	if found {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedItem)
		return
	}

	var item models.Item
	result := db.First(&item, id)
	if result.Error == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Item not found")
		return
	} else if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to retrieve item from the database")
		return
	}

	// Cache the retrieved item
	itemCache.Set(idStr, item, cache.DefaultExpiration)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func SaveItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request payload")
		return
	}

	// Set the creation timestamp for the new item
	newItem.CreatedAt = time.Now()

	// Create a new record in the database
	result := db.Create(&newItem)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to save item to the database")
		return
	}

	// Return the saved item as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid item ID")
		return
	}

	// Perform the deletion operation
	result := db.Delete(&models.Item{}, id)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to delete item from the database")
		return
	} else if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Item not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
