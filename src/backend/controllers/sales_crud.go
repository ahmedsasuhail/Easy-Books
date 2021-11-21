package controllers

import (
	"fmt"
	"net/http"

	"github.com/ahmedsasuhail/easy-books/models"
	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
)

// Meilisearch index.
var salesIndex = msClient.Index(models.SalesTableName)

// CreateSales creates a new record in the `eb_sales` table.
func CreateSales(c *gin.Context) {
	// Read and parse request body.
	var record models.Sales
	err := parseRequestBody(c, &record)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	// Create record in table and add it to Meilisearch index.
	err = pgClient.Create(&record).Error
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	pgClient.Preload(
		"Relationships",
	).Preload(
		"Purchases",
	).Preload(
		"Purchases.Relationships",
	).Preload(
		"Inventory",
	).Preload(
		"Inventory.Purchases.Relationships",
	).First(&record)
	filteredRecord := map[string]interface{}{
		"id":       record.ID,
		"price":    record.Price,
		"date":     record.Date,
		"credit":   record.Credit,
		"returned": record.Returned,
		"relationships": map[string]interface{}{
			"id":   record.Relationships.ID,
			"name": record.Relationships.Name,
		},
		"purchases": map[string]interface{}{
			"id":           record.Purchases.ID,
			"company_name": record.Purchases.CompanyName,
			"vehicle_name": record.Purchases.VehicleName,
			"price":        record.Purchases.Price,
		},
		"inventory": map[string]interface{}{
			"id":        record.Inventory.ID,
			"part_name": record.Inventory.PartName,
			"quantity":  record.Inventory.Quantity,
		},
	}

	_, err = salesIndex.AddDocuments(filteredRecord)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	successResponse(c, http.StatusOK, "", filteredRecord)
}

// UpdateSales updates a record in the `eb_sales` table.
func UpdateSales(c *gin.Context) {
	// Read and parse request body.
	var record models.Sales
	err := parseRequestBody(c, &record)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	// Update record in table and update Meilisearch index.
	err = pgClient.Model(&record).Updates(&record).Error
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	// Dirty fix to toggle boolean fields.
	if !record.Returned {
		err = pgClient.Model(&record).Select("Returned").Updates(
			models.Sales{Returned: false},
		).Error

		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())

			return
		}
	}

	if !record.Credit {
		err = pgClient.Model(&record).Select("Credit").Updates(
			models.Sales{Credit: false},
		).Error

		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())

			return
		}
	}

	pgClient.Preload(
		"Relationships",
	).Preload(
		"Purchases",
	).Preload(
		"Purchases.Relationships",
	).Preload(
		"Inventory",
	).Preload(
		"Inventory.Purchases",
	).Preload(
		"Inventory.Purchases.Relationships",
	).First(&record)
	filteredRecord := map[string]interface{}{
		"id":       record.ID,
		"price":    record.Price,
		"date":     record.Date,
		"credit":   record.Credit,
		"returned": record.Returned,
		"relationships": map[string]interface{}{
			"id":   record.Relationships.ID,
			"name": record.Relationships.Name,
		},
		"purchases": map[string]interface{}{
			"id":           record.Purchases.ID,
			"company_name": record.Purchases.CompanyName,
			"vehicle_name": record.Purchases.VehicleName,
			"price":        record.Purchases.Price,
		},
		"inventory": map[string]interface{}{
			"id":        record.Inventory.ID,
			"part_name": record.Inventory.PartName,
			"quantity":  record.Inventory.Quantity,
		},
	}

	_, err = salesIndex.UpdateDocuments(filteredRecord)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	successResponse(c, http.StatusOK, "", filteredRecord)
}

// ReadSales returns a paginated list of results from the `eb_sales`
// table.
func ReadSales(c *gin.Context) {
	pagination, err := parsePaginationRequest(c)
	if err != nil {
		errorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	offset := (pagination.Page - 1) * pagination.PageLimit
	stats, err := salesIndex.GetStats()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	records, err := salesIndex.Search(
		pagination.Query,
		&meilisearch.SearchRequest{
			Limit:  int64(pagination.PageLimit),
			Offset: int64(offset),
			Sort: []string{
				fmt.Sprintf("%s:%s", pagination.OrderBy, pagination.SortOrder),
			},
		})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	successResponse(c, http.StatusOK, "", map[string]interface{}{
		"page":                pagination.Page,
		"page_limit":          pagination.PageLimit,
		"order_by":            pagination.OrderBy,
		"sort_order":          pagination.SortOrder,
		"query":               pagination.Query,
		"total_count":         stats.NumberOfDocuments,
		"records":             records.Hits,
		"total_matched_count": len(records.Hits),
	})
}

// DeleteSales deletes a record from the `eb_sales` table.
func DeleteSales(c *gin.Context) {
	// Read and parse request body.
	var record models.Sales
	err := parseRequestBody(c, &record)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	err = pgClient.Where("id = ?", record.ID).Delete(&record).Error
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}
	filteredRecord := map[string]interface{}{
		"id":       record.ID,
		"price":    record.Price,
		"date":     record.Date,
		"credit":   record.Credit,
		"returned": record.Returned,
		"relationships": map[string]interface{}{
			"id":   record.Relationships.ID,
			"name": record.Relationships.Name,
		},
		"purchases": map[string]interface{}{
			"id":           record.Purchases.ID,
			"company_name": record.Purchases.CompanyName,
			"vehicle_name": record.Purchases.VehicleName,
			"price":        record.Purchases.Price,
		},
		"inventory": map[string]interface{}{
			"id":        record.Inventory.ID,
			"part_name": record.Inventory.PartName,
			"quantity":  record.Inventory.Quantity,
		},
	}

	_, err = salesIndex.DeleteDocument(fmt.Sprint(record.ID))
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	successResponse(c, http.StatusOK, "Deleted record.", filteredRecord)
}
