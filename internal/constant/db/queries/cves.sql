-- name: UpsertCVEs :exec 
INSERT INTO cves (
    cve_id,
    published_at,
    modified_at,
    cvss_version,
    cvss_score,
    severity,
    attack_vector,
    vector_string,
    descriptions,
    vendor,
    product,
    year,
    metrics,
    weaknesses,
    configurations,
    reference,
    affected
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14,
    $15, $16, $17
)
ON CONFLICT (cve_id)
DO UPDATE SET
    published_at   = EXCLUDED.published_at,
    modified_at    = EXCLUDED.modified_at,
    cvss_version   = EXCLUDED.cvss_version,
    cvss_score     = EXCLUDED.cvss_score,
    severity       = EXCLUDED.severity,
    attack_vector  = EXCLUDED.attack_vector,
    vector_string  = EXCLUDED.vector_string,
    vendor         = EXCLUDED.vendor,
    product        = EXCLUDED.product,
    year           = EXCLUDED.year,
    descriptions   = EXCLUDED.descriptions,
    metrics        = EXCLUDED.metrics,
    weaknesses     = EXCLUDED.weaknesses,
    configurations = EXCLUDED.configurations,
    reference      = EXCLUDED.reference,
    affected            = EXCLUDED.affected,
    updated_at     = NOW();