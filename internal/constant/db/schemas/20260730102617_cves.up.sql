CREATE TABLE cves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cve_id VARCHAR(32) NOT NULL UNIQUE,
    published_at TIMESTAMPTZ NOT NULL,
    modified_at TIMESTAMPTZ NOT NULL,
    cvss_version VARCHAR(10),
    cvss_score NUMERIC(3,1),
    severity VARCHAR(20),
    attack_vector VARCHAR(30),
    vector_string TEXT,
    vendor VARCHAR(255),
    product VARCHAR(255),
    year SMALLINT,
    descriptions JSONB,
    metrics JSONB,
    weaknesses JSONB,
    configurations JSONB,
    reference JSONB,
    affected JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_cves_cve_id
    ON cves(cve_id);

CREATE INDEX idx_cves_severity
    ON cves(severity);

CREATE INDEX idx_cves_score
    ON cves(cvss_score);

CREATE INDEX idx_cves_published
    ON cves(published_at DESC);

CREATE INDEX idx_cves_vendor
    ON cves(vendor);

CREATE INDEX idx_cves_product
    ON cves(product);

CREATE INDEX idx_cves_raw
    ON cves USING GIN(affected);

CREATE INDEX idx_cves_metrics
    ON cves USING GIN(metrics);

CREATE INDEX idx_cves_configurations
    ON cves USING GIN(configurations);