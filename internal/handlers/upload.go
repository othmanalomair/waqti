// internal/handlers/upload.go
package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"waqti/internal/middleware"

	"github.com/labstack/echo/v4"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

type UploadResponse struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	ImageURL string   `json:"image_url,omitempty"`
	Error    string   `json:"error,omitempty"`
	Images   []string `json:"images,omitempty"`
}

// UploadWorkshopImages handles multiple image uploads for workshops
func (h *UploadHandler) UploadWorkshopImages(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.JSON(http.StatusUnauthorized, UploadResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.Logger().Error("Error parsing multipart form:", err)
		return c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Error:   "Invalid form data",
		})
	}

	files := form.File["images[]"] // Note: this matches the name in the frontend
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Error:   "No files uploaded",
		})
	}

	// Ensure upload directory exists
	uploadDir := "web/static/images/upload"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.Logger().Error("Error creating upload directory:", err)
		return c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Error:   "Failed to create upload directory",
		})
	}

	var uploadedImages []string
	maxFiles := 3

	for i, file := range files {
		if i >= maxFiles {
			break // Limit to 3 files
		}

		// Validate file
		if !h.isValidImageFile(file) {
			continue // Skip invalid files
		}

		// Generate random filename
		filename, err := h.generateRandomFilename(file.Filename)
		if err != nil {
			c.Logger().Error("Error generating filename:", err)
			continue
		}

		// Save file
		savedPath, err := h.saveUploadedFile(file, uploadDir, filename)
		if err != nil {
			c.Logger().Error("Error saving file:", err)
			continue
		}

		// Convert to web path
		webPath := strings.Replace(savedPath, "web", "", 1)
		uploadedImages = append(uploadedImages, webPath)
	}

	if len(uploadedImages) == 0 {
		return c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Error:   "No valid images were uploaded",
		})
	}

	return c.JSON(http.StatusOK, UploadResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully uploaded %d images", len(uploadedImages)),
		Images:  uploadedImages,
	})
}

// UploadSingleImage handles single image upload
func (h *UploadHandler) UploadSingleImage(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.JSON(http.StatusUnauthorized, UploadResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Error:   "No file uploaded",
		})
	}

	// Validate file
	if !h.isValidImageFile(file) {
		return c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Error:   "Invalid image file",
		})
	}

	// Ensure upload directory exists
	uploadDir := "web/static/images/upload"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.Logger().Error("Error creating upload directory:", err)
		return c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Error:   "Failed to create upload directory",
		})
	}

	// Generate random filename
	filename, err := h.generateRandomFilename(file.Filename)
	if err != nil {
		c.Logger().Error("Error generating filename:", err)
		return c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Error:   "Failed to generate filename",
		})
	}

	// Save file
	savedPath, err := h.saveUploadedFile(file, uploadDir, filename)
	if err != nil {
		c.Logger().Error("Error saving file:", err)
		return c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Error:   "Failed to save file",
		})
	}

	// Convert to web path
	webPath := strings.Replace(savedPath, "web", "", 1)

	return c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "Image uploaded successfully",
		ImageURL: webPath,
	})
}

// isValidImageFile validates the uploaded file
func (h *UploadHandler) isValidImageFile(file *multipart.FileHeader) bool {
	// Check file size (2MB limit)
	if file.Size > 2*1024*1024 {
		return false
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !validExts[ext] {
		return false
	}

	// Check MIME type
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	// Read first 512 bytes to determine content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer)
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	return validTypes[contentType]
}

// generateRandomFilename creates a random filename while preserving extension
func (h *UploadHandler) generateRandomFilename(originalFilename string) (string, error) {
	// Generate 16 random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string
	randomName := hex.EncodeToString(bytes)

	// Get file extension
	ext := filepath.Ext(originalFilename)

	return randomName + ext, nil
}

// saveUploadedFile saves the multipart file to disk
func (h *UploadHandler) saveUploadedFile(file *multipart.FileHeader, uploadDir, filename string) (string, error) {
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file path
	dst := filepath.Join(uploadDir, filename)

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy file content
	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}

	return dst, nil
}

// DeleteImage deletes an uploaded image
func (h *UploadHandler) DeleteImage(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.JSON(http.StatusUnauthorized, UploadResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	// Get image path from request
	imagePath := c.FormValue("image_path")
	if imagePath == "" {
		return c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Error:   "No image path provided",
		})
	}

	// Security check: ensure path is in upload directory
	if !strings.Contains(imagePath, "/static/images/upload/") {
		return c.JSON(http.StatusForbidden, UploadResponse{
			Success: false,
			Error:   "Invalid image path",
		})
	}

	// Convert web path to file system path
	filePath := "web" + imagePath

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, UploadResponse{
			Success: false,
			Error:   "Image not found",
		})
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		c.Logger().Error("Error deleting file:", err)
		return c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Error:   "Failed to delete image",
		})
	}

	return c.JSON(http.StatusOK, UploadResponse{
		Success: true,
		Message: "Image deleted successfully",
	})
}
