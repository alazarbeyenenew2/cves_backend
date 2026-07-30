create table cves (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    cve_id VARCHAR(20) NOT NULL UNIQUE,
    published_at TIMESTAMPTZ NOT NULL,
    modified_at TIMESTAMPTZ NOT NULL,
    severity VARCHAR(20),
    cvss_score NUMERIC(3,1),
    access_vector VARCHAR(20),
    description TEXT NOT NULL,
    vendor VARCHAR(255),
    product VARCHAR(255),
    attack_type VARCHAR(50),
    year SMALLINT,
    cpe_uri TEXT,
    references TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_cve_id
ON cves(cve_id);