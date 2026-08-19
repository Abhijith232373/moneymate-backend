package s3util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type Client struct {
	s3         *s3.Client
	presign    *s3.PresignClient
	bucket     string
	publicBase string
}

func New(ctx context.Context, bucket, region, publicBase string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	return &Client{
		s3:         client,
		presign:    s3.NewPresignClient(client),
		bucket:     bucket,
		publicBase: strings.TrimRight(publicBase, "/"),
	}, nil
}

var allowedImageTypes = map[string]string{
	"image/png":  "png",
	"image/jpeg": "jpg",
	"image/webp": "webp",
}

type PresignedUpload struct {
	UploadURL string
	PublicURL string
	ExpiresAt time.Time
}

// PresignProfilePictureUpload — public-read prefix (users/ or merchants/).
func (c *Client) PresignProfilePictureUpload(ctx context.Context, ownerPrefix string, ownerID uuid.UUID, contentType string) (*PresignedUpload, error) {
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	key := fmt.Sprintf("%s/%s/profile-%s.%s", ownerPrefix, ownerID, uuid.New(), ext)
	expiresIn := 5 * time.Minute

	req, err := c.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return nil, fmt.Errorf("presign put object: %w", err)
	}

	return &PresignedUpload{
		UploadURL: req.URL,
		PublicURL: fmt.Sprintf("%s/%s", c.publicBase, key),
		ExpiresAt: time.Now().Add(expiresIn),
	}, nil
}

// IsOwnedURL verifies a public URL actually falls under this owner's prefix
// before trusting it enough to persist to the DB.
func (c *Client) IsOwnedURL(url, ownerPrefix string, ownerID uuid.UUID) bool {
	expectedPrefix := fmt.Sprintf("%s/%s/%s/", c.publicBase, ownerPrefix, ownerID)
	return strings.HasPrefix(url, expectedPrefix)
}

// ── KYC (private) — not wired into routes yet, ready for when merchant KYC is tackled ──

var allowedDocTypes = map[string]string{
	"application/pdf": "pdf",
	"image/png":        "png",
	"image/jpeg":       "jpg",
}

// PresignKYCDocumentUpload — private prefix, no public bucket-policy grant covers "kyc/*".
func (c *Client) PresignKYCDocumentUpload(ctx context.Context, storeID uuid.UUID, docKind, contentType string) (*PresignedUpload, string, error) {
	ext, ok := allowedDocTypes[contentType]
	if !ok {
		return nil, "", fmt.Errorf("unsupported content type: %s", contentType)
	}
	key := fmt.Sprintf("kyc/%s/%s-%s.%s", storeID, docKind, uuid.New(), ext)
	expiresIn := 5 * time.Minute

	req, err := c.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return nil, "", fmt.Errorf("presign kyc put object: %w", err)
	}
	return &PresignedUpload{UploadURL: req.URL, ExpiresAt: time.Now().Add(expiresIn)}, key, nil
}

// PresignKYCDocumentGet — for admin review screens only; short-lived, key-based (not a public URL).
func (c *Client) PresignKYCDocumentGet(ctx context.Context, key string) (string, error) {
	expiresIn := 10 * time.Minute
	req, err := c.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", fmt.Errorf("presign kyc get object: %w", err)
	}
	return req.URL, nil
}