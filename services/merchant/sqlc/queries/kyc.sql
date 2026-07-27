-- name: GetKYCStatusByStoreID :one
-- GetKYCStatusByStoreID retrieves the verification standing and document URLs for a merchant store joined with core store status.
SELECT 
    k.id, k.store_id, k.aadhaar_number, k.aadhaar_doc_url, k.shop_license_url, k.is_verified, k.verified_at, k.created_at, k.updated_at, s.status::text AS store_status
FROM kyc_documents k
JOIN stores s ON k.store_id = s.id
WHERE k.store_id = $1 LIMIT 1;

-- name: GetStoreStatusByID :one
-- GetStoreStatusByID retrieves the current state machine status string for a specific merchant store.
SELECT status::text AS store_status FROM stores WHERE id = $1 LIMIT 1;

-- name: InsertKYCDocuments :one
-- InsertKYCDocuments inserts a new compliance documentation record during store onboarding or fallback initialization.
INSERT INTO kyc_documents (
    id, store_id, aadhaar_number, aadhaar_doc_url, shop_license_url, is_verified, verified_at, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW()
) RETURNING id, store_id, aadhaar_number, aadhaar_doc_url, shop_license_url, is_verified, verified_at, created_at, updated_at;

-- name: UpdateKYCDocumentsByStoreID :one
-- UpdateKYCDocumentsByStoreID modifies existing compliance file URLs and resets verification standing to unverified.
UPDATE kyc_documents
SET aadhaar_number = COALESCE(NULLIF(sqlc.arg('aadhaar_number')::text, ''), aadhaar_number),
    aadhaar_doc_url = COALESCE(NULLIF(sqlc.arg('aadhaar_doc_url')::text, ''), aadhaar_doc_url),
    shop_license_url = COALESCE(NULLIF(sqlc.arg('shop_license_url')::text, ''), shop_license_url),
    is_verified = FALSE,
    verified_at = NULL,
    updated_at = NOW()
WHERE store_id = sqlc.arg('store_id')
RETURNING id, store_id, aadhaar_number, aadhaar_doc_url, shop_license_url, is_verified, verified_at, created_at, updated_at;

-- name: UpdateStoreStatusByID :one
-- UpdateStoreStatusByID transitions the store's overarching state machine status and returns the updated text.
UPDATE stores
SET status = sqlc.arg('status')::merchant_status, updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING status::text;

