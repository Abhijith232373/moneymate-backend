package http

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

// KYCHandler exposes REST endpoints for viewing business verification standing and submitting updated compliance documentation.
type KYCHandler struct {
	// usecase orchestrates document validation and atomic compliance status updates.
	usecase usecases.KYCUseCase
}

// NewKYCHandler initializes and returns a new KYCHandler instance with injected dependencies.
func NewKYCHandler(uc usecases.KYCUseCase) *KYCHandler {
	return &KYCHandler{usecase: uc}
}

// GetKYCStatus handles GET /merchant/:store_id/kyc (and /status) and returns the full compliance verification dashboard data.
func (h *KYCHandler) GetKYCStatus(c fiber.Ctx) error {
	storeIDStr := c.Params("store_id")
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	kyc, err := h.usecase.GetStatus(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	statusBadge := "Pending Review"
	message := "Your documents are currently under review by our compliance team."
	if kyc.IsVerified || kyc.StoreStatus == "active" || kyc.StoreStatus == "verified" {
		statusBadge = "Verified"
		message = "Your business identity has been verified. You have full access to all merchant features and high transaction limits."
	} else if kyc.StoreStatus == "rejected" {
		statusBadge = "Rejected"
		message = "One or more of your submitted compliance documents could not be verified. Please update your documents below."
	}

	docStatus := "Pending Review"
	switch statusBadge {
	case "Verified":
		docStatus = "Approved"
	case "Rejected":
		docStatus = "Action Required"
	}

	submitDate := kyc.CreatedAt.Format("Oct 02, 2006")
	if kyc.VerifiedAt != nil {
		submitDate = kyc.VerifiedAt.Format("Oct 02, 2006")
	}

	formattedDocStatus := fmt.Sprintf("%s • %s", docStatus, submitDate)

	var documents []KYCDocumentItemResponse
	if kyc.ShopLicenseURL != "" {
		documents = append(documents, KYCDocumentItemResponse{
			Title:           "Shop License / Registration",
			Status:          docStatus,
			SubmittedAt:     submitDate,
			FormattedStatus: formattedDocStatus,
			URL:             kyc.ShopLicenseURL,
			DocType:         "shop_license",
		})
	}
	if kyc.AadhaarDocURL != "" {
		documents = append(documents, KYCDocumentItemResponse{
			Title:           "Aadhaar Card Document",
			Status:          docStatus,
			SubmittedAt:     submitDate,
			FormattedStatus: formattedDocStatus,
			URL:             kyc.AadhaarDocURL,
			DocType:         "aadhaar",
		})
	}
	if documents == nil {
		documents = []KYCDocumentItemResponse{}
	}

	nextReviewYear := time.Now().Year() + 3
	if kyc.VerifiedAt != nil {
		nextReviewYear = kyc.VerifiedAt.Year() + 3
	}
	nextReviewDue := fmt.Sprintf("Next review due: Oct %d", nextReviewYear)

	return c.JSON(KYCStatusResponse{
		StoreID:       kyc.StoreID.String(),
		Status:        statusBadge,
		IsVerified:    statusBadge == "Verified",
		Message:       message,
		NextReviewDue: nextReviewDue,
		Documents:     documents,
	})
}

// UpdateKYCDocuments handles PUT/POST /merchant/:store_id/kyc/update to submit new document URLs and trigger re-verification.
func (h *KYCHandler) UpdateKYCDocuments(c fiber.Ctx) error {
	storeIDStr := c.Params("store_id")
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	var req UpdateKYCDocumentsRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON request payload"})
	}

	kyc, err := h.usecase.UpdateDocuments(c.Context(), storeID, req.AadhaarNumber, req.AadhaarDocURL, req.ShopLicenseURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(UpdateKYCDocumentsResponse{
		StoreID:    kyc.StoreID.String(),
		Status:     "Pending Review",
		IsVerified: false,
		Message:    "Your updated compliance documents have been submitted and are under review by our compliance team.",
	})
}
