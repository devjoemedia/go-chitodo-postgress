package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/devjoemedia/scrumpilot-go-api/database"
	"github.com/devjoemedia/scrumpilot-go-api/middleware"
	"github.com/devjoemedia/scrumpilot-go-api/models"
	"github.com/devjoemedia/scrumpilot-go-api/types"
	api_response "github.com/devjoemedia/scrumpilot-go-api/types/response"
	"github.com/devjoemedia/scrumpilot-go-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// CreateTicket godoc
// @Summary      Create a new ticket
// @Description  Create a new ticket item
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security 		 BearerAuth
// @Param        body  body     models.CreateTicketRequest  true  "Ticket object"
// @Success      200   {object} api_response.CreateTicketResponse
// @Failure      400   {string} string      "Invalid JSON"
// @Router       /api/v1/tickets [post]
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid input")
		return
	}

	userID, _, err := middleware.GetUserIDAndEmailFromRequest(r)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get user ID and email")
		return
	}

	ticket := models.Ticket{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		ReporterID:  userID,
		AssigneeID:  nil,
	}
	if req.AssigneeID != nil {
		ticket.AssigneeID = req.AssigneeID
	}

	ctx := r.Context()
	result := database.DB.WithContext(ctx).Create(&ticket)

	if result.Error != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create ticket")
		return
	}

	// 🔥 IMPORTANT: Reload with relations
	if err := database.DB.WithContext(ctx).
		Preload("Reporter").
		Preload("Assignee").
		First(&ticket, ticket.ID).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load ticket relations")
		return
	}

	response := api_response.CreateTicketResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Ticket created successfully",
		Ticket:  &ticket,
	}

	utils.JSON(w, http.StatusOK, response)
}

// GetTickets godoc
// @Summary      Get tickets with pagination and filters
// @Description  Retrieve tickets with optional pagination (page, size) and filter by status, priority, assignee_id, reporter_id
// @Tags         tickets
// @Produce      json
// @Security 		 BearerAuth
// @Param        page        query    int    false  "Page number (default: 1)"
// @Param        size        query    int    false  "Page size (default: 10, max: 100)"
// @Param        status      query    string   false  "Filter by status"
// @Param        priority    query    string   false  "Filter by priority"
// @Param        assignee_id query    int      false  "Filter by assignee_id"
// @Param        reporter_id query    int      false  "Filter by reporter_id"
// @Success      200 {object} api_response.GetTicketsResponse
// @Router       /api/v1/tickets [get]
func GetTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Pagination
	page := 1
	size := 10

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil && s > 0 {
		size = s
	}

	if size > 100 {
		size = 100
	}

	// Filters
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	var assigneeID *uint
	if v := r.URL.Query().Get("assignee_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			u := uint(id)
			assigneeID = &u
		}
	}

	var reporterID *uint
	if v := r.URL.Query().Get("reporter_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			u := uint(id)
			reporterID = &u
		}
	}

	// Initialize tickets slice
	var tickets []models.Ticket

	// Build Base Query
	query := database.DB.WithContext(ctx).Model(&models.Ticket{})

	// Apply filter
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assigneeID != nil {
		query = query.Where("assignee_id = ?", *assigneeID)
	}
	if reporterID != nil {
		query = query.Where("reporter_id = ?", *reporterID)
	}

	// Count total tickets
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Execute query
	if err := query.
		Preload("Reporter").
		Preload("Assignee").
		Offset((page - 1) * size).
		Limit(size).
		Find(&tickets).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Return response
	utils.JSON(w, http.StatusOK, api_response.GetTicketsResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "success",
		Tickets: tickets,
		Pagination: types.Pagination{
			Total: int(total),
			Size:  size,
			Page:  page,
		},
	})
}

// GetTicketByID godoc
// @Summary      Get ticket by ID
// @Description  Fetch a specific ticket by its ID
// @Tags         tickets
// @Produce      json
// @Security 		 BearerAuth
// @Param        id   path    int  true  "Ticket ID"
// @Success      200  {object} api_response.GetTicketResponse
// @Failure      400  {string} string     "Invalid ID"
// @Failure      404  {string} string     "Ticket not found"
// @Router       /api/v1/tickets/{id} [get]
func GetTicketByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Ticket Slice
	var ticket models.Ticket
	if err := database.DB.WithContext(ctx).First(&ticket, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Ticket not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Return response
	utils.JSON(w, http.StatusOK, api_response.GetTicketResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Ticket retrieved successfully",
		Ticket:  &ticket,
	})
}

// UpdateTodo godoc
// @Summary      Update an existing todo
// @Description  Update a todo item by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security 		 BearerAuth
// @Param        id     path    int  true  "Todo ID"
// @Param        body   body    models.UpdateTodoRequest  true  "Todo object"
// @Success      200    {object} api_response.UpdateTodoResponse
// @Failure      400    {string} string  "Invalid JSON"
// @Failure      404    {string} string  "Todo not found"
// @Router       /api/v1/todos/{id} [put]
// @Router       /api/v1/todos/{id} [patch]
func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req models.UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Ticket Slice
	var ticket models.Ticket
	if err := database.DB.WithContext(ctx).First(&ticket, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Ticket not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Update fields
	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if req.Status != nil {
		ticket.Status = *req.Status
	}
	if req.Priority != nil {
		ticket.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		ticket.AssigneeID = req.AssigneeID
	}

	// save updates
	result := database.DB.WithContext(ctx).Save(&ticket).Error
	if result != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update ticket")
		return
	}

	// Response
	response := api_response.UpdateTicketResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Ticket updated successfully",
		Ticket:  &ticket,
	}

	utils.JSON(w, http.StatusOK, response)
}

// DeleteTicket godoc
// @Summary      Delete ticket by ID
// @Description  Delete a specific ticket by its ID
// @Tags         tickets
// @Produce      json
// @Security 		 BearerAuth
// @Param        id   path    int  true  "Ticket ID"
// @Success      200  {object} api_response.DeleteTicketResponse
// @Failure      400  {string} string             "Invalid ID"
// @Failure      404  {string} string             "Ticket not found"
// @Router       /api/v1/tickets/{id} [delete]
func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Delete ticket
	if err := database.DB.WithContext(ctx).Delete(&models.Ticket{}, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Ticket not found")
			return
		}

		utils.Error(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	// Response
	response := api_response.DeleteTicketResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Ticket deleted successfully",
	}

	utils.JSON(w, http.StatusOK, response)
}
