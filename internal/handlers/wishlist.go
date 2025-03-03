package handlers

import (
	"ace/internal/models"
	"ace/internal/repository"
	"ace/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type WishlistHandler struct {
	Repo *repository.WishlistRepository
}

// NewWishlistHandler creates a new handler instance
func NewWishlistHandler(repo *repository.WishlistRepository) *WishlistHandler {
	return &WishlistHandler{Repo: repo}
}

// CreateWishlist Create new wishlist item
func (h *WishlistHandler) CreateWishlist(c *gin.Context) {
	var wishlist models.Wishlist

	if err := c.ShouldBindJSON(&wishlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	if errors := utils.ValidateStruct(wishlist); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	id, err := h.Repo.Create(c.Request.Context(), &wishlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed: " + err.Error()})
		return
	}
	wishlist.ID = id

	c.JSON(http.StatusOK, wishlist)
}

// UpdateWishlist Update wishlist item
func (h *WishlistHandler) UpdateWishlist(c *gin.Context) {
	id := c.Param("id")
	var wishlist models.Wishlist

	if err := c.ShouldBindJSON(&wishlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	if errors := utils.ValidateStruct(wishlist); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	err := h.Repo.Update(c.Request.Context(), id, &wishlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wishlist updated successfully"})
}

// GetAllWishlist get with pagination
func (h *WishlistHandler) GetAllWishlist(c *gin.Context) {
	pagination := utils.GetPagination(c)
	wishlists, total, err := h.Repo.GetAll(c.Request.Context(), pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlists: " + err.Error()})
		return
	}

	pagination.CalculateTotalPage(total)
	c.JSON(http.StatusOK, gin.H{
		"data":       wishlists,
		"pagination": pagination,
	})
}

// GetOneWishlist get detail wishlist item
func (h *WishlistHandler) GetOneWishlist(c *gin.Context) {
	id := c.Param("id")
	wishlist, err := h.Repo.GetOne(c.Request.Context(), id)
	if err != nil {
		print(err.Error())
		if strings.Contains(err.Error(), "invalid ID") || strings.Contains(err.Error(), "cannot decode") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		if strings.Contains(err.Error(), "no document") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, wishlist)
}

func (h *WishlistHandler) DeleteWishlist(c *gin.Context) {
	id := c.Param("id")
	err := h.Repo.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error()[:12] == "invalid ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		if err.Error()[:15] == "no document" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Wishlist deleted successfully"})
}
