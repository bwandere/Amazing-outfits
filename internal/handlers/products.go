package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"amazing-outfits/internal/models"
)

type ProductHandler struct {
	DB *sql.DB
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	brand := q.Get("brand")
	size := q.Get("size")
	priceStr := q.Get("price")
	sort := q.Get("sort")
	search := q.Get("search")
	category := q.Get("category")
	prodType := q.Get("type")

	query := `
		SELECT p.id, p.name, p.brand, p.type, p.category, p.price, p.old_price, 
		       p.rating, p.reviews, p.description, p.tags, p.created_at, p.updated_at
		FROM products p
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if brand != "" {
		query += fmt.Sprintf(" AND LOWER(p.brand) = LOWER($%d)", argIdx)
		args = append(args, brand)
		argIdx++
	}

	if prodType != "" {
		query += fmt.Sprintf(" AND p.type = $%d", argIdx)
		args = append(args, prodType)
		argIdx++
	}

	if category != "" && category != "All" {
		query += fmt.Sprintf(" AND LOWER(p.category) = LOWER($%d)", argIdx)
		args = append(args, category)
		argIdx++
	}

	if priceStr != "" {
		if maxP, err := strconv.ParseFloat(priceStr, 64); err == nil && maxP > 0 {
			query += fmt.Sprintf(" AND p.price <= $%d", argIdx)
			args = append(args, maxP)
			argIdx++
		}
	}

	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(p.name) LIKE LOWER($%d) OR LOWER(p.description) LIKE LOWER($%d))", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if size != "" {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM product_sizes ps WHERE ps.product_id = p.id AND ps.size = $%d)", argIdx)
		args = append(args, size)
		argIdx++
	}

	switch sort {
	case "price-asc":
		query += " ORDER BY p.price ASC"
	case "price-desc":
		query += " ORDER BY p.price DESC"
	case "rating":
		query += " ORDER BY p.rating DESC"
	default:
		query += " ORDER BY p.created_at ASC"
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Query error: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var tagsJSON []byte
		err := rows.Scan(
			&p.ID, &p.Name, &p.Brand, &p.Type, &p.Category, &p.Price, &p.OldPrice,
			&p.Rating, &p.Reviews, &p.Description, &tagsJSON, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			http.Error(w, `{"error":"Error scanning product"}`, http.StatusInternalServerError)
			return
		}
		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &p.Tags)
		}
		if p.Tags == nil {
			p.Tags = []string{}
		}

		// Fetch Images
		imgRows, _ := h.DB.Query("SELECT url FROM product_images WHERE product_id = $1 ORDER BY display_order ASC", p.ID)
		p.Images = []string{}
		if imgRows != nil {
			for imgRows.Next() {
				var u string
				if err := imgRows.Scan(&u); err == nil {
					p.Images = append(p.Images, u)
				}
			}
			imgRows.Close()
		}

		// Fetch Sizes
		sizeRows, _ := h.DB.Query("SELECT size FROM product_sizes WHERE product_id = $1 ORDER BY id ASC", p.ID)
		p.Sizes = []string{}
		if sizeRows != nil {
			for sizeRows.Next() {
				var s string
				if err := sizeRows.Scan(&s); err == nil {
					p.Sizes = append(p.Sizes, s)
				}
			}
			sizeRows.Close()
		}

		products = append(products, p)
	}

	if products == nil {
		products = []models.Product{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, `{"error":"Product ID required"}`, http.StatusBadRequest)
		return
	}

	var p models.Product
	var tagsJSON []byte
	err := h.DB.QueryRow(`
		SELECT id, name, brand, type, category, price, old_price, 
		       rating, reviews, description, tags, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.Name, &p.Brand, &p.Type, &p.Category, &p.Price, &p.OldPrice,
		&p.Rating, &p.Reviews, &p.Description, &tagsJSON, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Product not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &p.Tags)
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}

	// Images
	imgRows, _ := h.DB.Query("SELECT url FROM product_images WHERE product_id = $1 ORDER BY display_order ASC", p.ID)
	p.Images = []string{}
	if imgRows != nil {
		for imgRows.Next() {
			var u string
			if err := imgRows.Scan(&u); err == nil {
				p.Images = append(p.Images, u)
			}
		}
		imgRows.Close()
	}

	// Sizes
	sizeRows, _ := h.DB.Query("SELECT size FROM product_sizes WHERE product_id = $1 ORDER BY id ASC", p.ID)
	p.Sizes = []string{}
	if sizeRows != nil {
		for sizeRows.Next() {
			var s string
			if err := sizeRows.Scan(&s); err == nil {
				p.Sizes = append(p.Sizes, s)
			}
		}
		sizeRows.Close()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	if p.ID == "" {
		p.ID = strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))
	}
	if p.Brand == "" {
		p.Brand = "Amazing"
	}
	if p.Type == "" {
		p.Type = "sneaker"
	}

	tagsJSON, _ := json.Marshal(p.Tags)

	_, err := h.DB.Exec(`
		INSERT INTO products (id, name, brand, type, category, price, old_price, rating, reviews, description, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, p.ID, p.Name, p.Brand, p.Type, p.Category, p.Price, p.OldPrice, p.Rating, p.Reviews, p.Description, tagsJSON)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Failed to create product: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Insert sizes
	for _, s := range p.Sizes {
		h.DB.Exec("INSERT INTO product_sizes (product_id, size) VALUES ($1, $2) ON CONFLICT DO NOTHING", p.ID, s)
	}

	// Insert images
	for idx, u := range p.Images {
		h.DB.Exec("INSERT INTO product_images (product_id, url, display_order) VALUES ($1, $2, $3)", p.ID, u, idx)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	tagsJSON, _ := json.Marshal(p.Tags)

	_, err := h.DB.Exec(`
		UPDATE products
		SET name = $1, brand = $2, type = $3, category = $4, price = $5, old_price = $6,
		    description = $7, tags = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`, p.Name, p.Brand, p.Type, p.Category, p.Price, p.OldPrice, p.Description, tagsJSON, id)

	if err != nil {
		http.Error(w, `{"error":"Failed to update product"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := h.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete product"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}
