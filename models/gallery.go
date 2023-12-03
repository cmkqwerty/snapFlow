package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	// ImagesDir is the path to the directory where images are stored. Defaults to "images/".
	ImagesDir string
}

func (service *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}

	row := service.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}

	row := service.DB.QueryRow(`
	SELECT title, user_id
	FROM galleries
	WHERE id = $1;`, gallery.ID)
	err := row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by ID: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := service.DB.Query(`
	SELECT id, title
	FROM galleries
	WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by userID: %w", err)
	}

	var galleries []Gallery

	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}

		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("scan galleries: %w", err)
		}

		galleries = append(galleries, gallery)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by userID: %w", err)
	}

	return galleries, nil
}

func (service *GalleryService) Update(gallery *Gallery) error {
	_, err := service.DB.Exec(`
	UPDATE galleries
	SET title = $2
	WHERE id = $1;`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}

	return nil
}

func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`
	DELETE FROM galleries
	WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	return nil
}

func (service *GalleryService) Images(id int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(id), "*")

	files, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery: %w", err)
	}

	var images []Image
	for _, file := range files {
		if hasExtension(file, service.extensions()) {
			images = append(images, Image{
				GalleryID: id,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}

	return images, nil
}

func (service *GalleryService) Image(id int, filename string) (Image, error) {
	imagePath := filepath.Join(service.galleryDir(id), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying image: %w", err)
	}

	return Image{
		GalleryID: id,
		Path:      imagePath,
		Filename:  filename,
	}, nil
}

func (service *GalleryService) CreateImage(id int, filename string, contents io.Reader) error {
	galleryDir := service.galleryDir(id)

	err := os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating image directory: %w", err)
	}

	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents of image: %w", err)
	}

	return nil
}

func (service *GalleryService) DeleteImage(id int, filename string) error {
	image, err := service.Image(id, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	return nil
}

func (service *GalleryService) extensions() []string {
	return []string{".jpg", ".jpeg", ".png", ".gif", ".avif"}
}

func (service *GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func hasExtension(file string, extensions []string) bool {
	for _, extension := range extensions {
		file = strings.ToLower(file)
		extension = strings.ToLower(extension)
		if filepath.Ext(file) == extension {
			return true
		}
	}

	return false
}
