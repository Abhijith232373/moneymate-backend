package http

// import (
// 	"time"

// 	"github.com/gofiber/fiber/v3"
// 	"github.com/google/uuid"

// 	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/tokenissuer"
// 	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
// 	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
// 	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
// )

// type PinHandler struct {
// 	pins   usecase.PinUsecase
// 	issuer *tokenissuer.Issuer
// }

// func NewPinHandler(pins usecase.PinUsecase, issuer *tokenissuer.Issuer) *PinHandler {
// 	return &PinHandler{pins: pins, issuer: issuer}
// }

// func (h *PinHandler) Setup(c fiber.Ctx) error {
// 	var req pinRequest
// 	if err := c.Bind().Body(&req); err != nil {
// 		return response.BadRequest(c, nil, "invalid request body")
// 	}
// 	if err := validate.Struct(req); err != nil {
// 		return response.BadRequest(c, nil, "validation failed")
// 	}

// 	userID, err := userIDFromLocals(c)
// 	if err != nil {
// 		return response.Unauthorized(c, "authentication required")
// 	}
// 	if err := h.pins.SetPIN(c.Context(), userID, req.PIN); err != nil {
// 		return handleError(c, err)
// 	}
// 	return response.OK(c, "transaction pin set", nil)
// }

// func (h *PinHandler) Verify(c fiber.Ctx) error {
// 	var req pinRequest
// 	if err := c.Bind().Body(&req); err != nil {
// 		return response.BadRequest(c, nil, "invalid request body")
// 	}
// 	if err := validate.Struct(req); err != nil {
// 		return response.BadRequest(c, nil, "validation failed")
// 	}

// 	userID, err := userIDFromLocals(c)
// 	if err != nil {
// 		return response.Unauthorized(c, "authentication required")
// 	}
// 	if err := h.pins.VerifyPIN(c.Context(), userID, req.PIN); err != nil {
// 		return handleError(c, err)
// 	}

// 	token, expiresAt, err := h.issuer.IssueTransactionToken(userID)
// 	if err != nil {
// 		return response.InternalServerError(c)
// 	}
// 	return response.OK(c, "verified", pinVerifyResponse{
// 		TransactionToken: token,
// 		ExpiresIn: int(time.Until(expiresAt).Seconds()),
// 	})
// }

// func userIDFromLocals(c fiber.Ctx) (uuid.UUID, error) {
// 	idStr, ok := c.Locals("userID").(string)
// 	if !ok || idStr == "" {
// 		return uuid.Nil, apperrors.ErrUnauthorized
// 	}
// 	return uuid.Parse(idStr)
// }
