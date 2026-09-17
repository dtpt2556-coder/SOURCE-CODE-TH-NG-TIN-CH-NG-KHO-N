-- 0001_init.sql — Lược đồ khởi tạo TenPoint (mục 4 architecture.md).
-- Idempotent: chạy lại nhiều lần không lỗi.

CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- unaccent() không IMMUTABLE nên không dùng trực tiếp trong GENERATED column.
-- Bọc lại bằng hàm IMMUTABLE trỏ tới đúng dictionary (kỹ thuật chuẩn của PG).
CREATE OR REPLACE FUNCTION tenpoint_unaccent(text)
RETURNS text
LANGUAGE sql IMMUTABLE PARALLEL SAFE STRICT AS
$$ SELECT public.unaccent('public.unaccent'::regdictionary, $1) $$;

-- ---------------------------------------------------------------- sources
CREATE TABLE IF NOT EXISTS sources (
    id            serial PRIMARY KEY,
    code          text     NOT NULL UNIQUE,
    name          text     NOT NULL,
    domain        text     NOT NULL UNIQUE,
    rss_url       text     NOT NULL DEFAULT '',
    list_url      text     NOT NULL DEFAULT '',
    tier          smallint NOT NULL DEFAULT 1,
    enabled       boolean  NOT NULL DEFAULT true,
    rate_limit_ms integer  NOT NULL DEFAULT 2000,
    created_at    timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT sources_tier_check CHECK (tier BETWEEN 0 AND 3),
    CONSTRAINT sources_rate_limit_check CHECK (rate_limit_ms >= 0)
);

-- ---------------------------------------------------------------- tickers
CREATE TABLE IF NOT EXISTS tickers (
    symbol       text PRIMARY KEY,
    company_name text NOT NULL,
    short_name   text NOT NULL DEFAULT '',
    exchange     text NOT NULL DEFAULT 'HOSE',
    sector       text NOT NULL DEFAULT '',
    in_vn30      boolean NOT NULL DEFAULT false,
    status       text NOT NULL DEFAULT 'active',
    CONSTRAINT tickers_exchange_check CHECK (exchange IN ('HOSE','HNX','UPCOM')),
    CONSTRAINT tickers_status_check CHECK (status IN ('active','suspended','delisted'))
);

CREATE TABLE IF NOT EXISTS ticker_aliases (
    id               bigserial PRIMARY KEY,
    symbol           text NOT NULL REFERENCES tickers(symbol) ON DELETE CASCADE,
    alias            text NOT NULL,
    alias_normalized text NOT NULL,
    CONSTRAINT ticker_aliases_uniq UNIQUE (symbol, alias_normalized)
);

CREATE INDEX IF NOT EXISTS idx_ticker_aliases_normalized
    ON ticker_aliases (alias_normalized);

-- ---------------------------------------------------------- article_clusters
CREATE TABLE IF NOT EXISTS article_clusters (
    id                     bigserial PRIMARY KEY,
    representative_simhash text NOT NULL,
    window_start           timestamptz NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------- articles
CREATE TABLE IF NOT EXISTS articles (
    id                      bigserial PRIMARY KEY,
    source_id               integer NOT NULL REFERENCES sources(id),
    canonical_url           text NOT NULL UNIQUE,
    url_hash                text NOT NULL UNIQUE,
    title                   text NOT NULL,
    summary_md              text NOT NULL DEFAULT '',
    -- excerpt là dữ liệu NỘI BỘ cho validator, KHÔNG BAO GIỜ trả ra API (ADR-006).
    excerpt                 text NOT NULL DEFAULT '',
    news_type               text NOT NULL,
    title_simhash           text NOT NULL DEFAULT '',
    cluster_id              bigint REFERENCES article_clusters(id),
    is_canonical_in_cluster boolean NOT NULL DEFAULT true,
    status                  text NOT NULL DEFAULT 'published',
    summary_provider        text NOT NULL DEFAULT 'extractive',
    published_at            timestamptz NOT NULL,
    fetched_at              timestamptz NOT NULL DEFAULT now(),
    search_vec tsvector GENERATED ALWAYS AS (
        to_tsvector('simple',
            tenpoint_unaccent(coalesce(title, '') || ' ' || coalesce(summary_md, '')))
    ) STORED,
    CONSTRAINT articles_news_type_check CHECK (news_type IN (
        'vi_mo','nganh','doanh_nghiep','thi_truong','khoi_ngoai',
        'co_tuc_phat_hanh','phap_ly')),
    CONSTRAINT articles_status_check CHECK (status IN ('published','pending_review','rejected')),
    CONSTRAINT articles_provider_check CHECK (summary_provider IN ('anthropic','extractive'))
);

-- Truy vấn nóng nhất: trang digest.
CREATE INDEX IF NOT EXISTS idx_articles_published_desc
    ON articles (published_at DESC) WHERE status = 'published';
CREATE INDEX IF NOT EXISTS idx_articles_search_vec
    ON articles USING GIN (search_vec);
CREATE INDEX IF NOT EXISTS idx_articles_newstype_published
    ON articles (news_type, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_cluster
    ON articles (cluster_id) WHERE cluster_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_articles_simhash
    ON articles (title_simhash);

-- ---------------------------------------------------------- article_tickers
CREATE TABLE IF NOT EXISTS article_tickers (
    article_id bigint NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    symbol     text   NOT NULL REFERENCES tickers(symbol) ON DELETE CASCADE,
    relevance  text   NOT NULL DEFAULT 'mentioned',
    score      real   NOT NULL DEFAULT 0,
    PRIMARY KEY (article_id, symbol),
    CONSTRAINT article_tickers_relevance_check CHECK (relevance IN ('primary','mentioned')),
    CONSTRAINT article_tickers_score_check CHECK (score >= 0 AND score <= 1)
);

CREATE INDEX IF NOT EXISTS idx_article_tickers_symbol
    ON article_tickers (symbol, article_id);

-- ---------------------------------------------------------------- short_links
CREATE TABLE IF NOT EXISTS short_links (
    id              bigserial PRIMARY KEY,
    code            text NOT NULL UNIQUE,
    article_id      bigint NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    target_url      text NOT NULL,
    click_count     bigint NOT NULL DEFAULT 0,
    last_checked_at timestamptz,
    target_alive    boolean NOT NULL DEFAULT true,
    created_at      timestamptz NOT NULL DEFAULT now()
);

-- Một article = một short_code vĩnh viễn (R7.10 của BA2).
CREATE UNIQUE INDEX IF NOT EXISTS idx_short_links_article
    ON short_links (article_id);

-- ------------------------------------------------------------ research_notes
CREATE TABLE IF NOT EXISTS research_notes (
    id              bigserial PRIMARY KEY,
    slug            text NOT NULL UNIQUE,
    symbol          text REFERENCES tickers(symbol) ON DELETE SET NULL,
    title           text NOT NULL,
    section_heading text NOT NULL DEFAULT 'LUẬN ĐIỂM ĐẦU TƯ',
    disclaimer      text NOT NULL DEFAULT '',
    published_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_research_notes_symbol
    ON research_notes (symbol, published_at DESC);

CREATE TABLE IF NOT EXISTS research_note_points (
    id      bigserial PRIMARY KEY,
    note_id bigint NOT NULL REFERENCES research_notes(id) ON DELETE CASCADE,
    ordinal integer NOT NULL,
    lead    text NOT NULL,
    body    text NOT NULL,
    CONSTRAINT research_note_points_uniq UNIQUE (note_id, ordinal)
);

-- ---------------------------------------------------------------- crawl_runs
CREATE TABLE IF NOT EXISTS crawl_runs (
    id                bigserial PRIMARY KEY,
    trigger           text NOT NULL DEFAULT 'cron',
    started_at        timestamptz NOT NULL DEFAULT now(),
    finished_at       timestamptz,
    articles_found    integer NOT NULL DEFAULT 0,
    articles_new      integer NOT NULL DEFAULT 0,
    articles_rejected integer NOT NULL DEFAULT 0,
    errors            integer NOT NULL DEFAULT 0,
    CONSTRAINT crawl_runs_trigger_check CHECK (trigger IN ('cron','manual'))
);

CREATE INDEX IF NOT EXISTS idx_crawl_runs_started
    ON crawl_runs (started_at DESC);

CREATE TABLE IF NOT EXISTS crawl_run_sources (
    id            bigserial PRIMARY KEY,
    run_id        bigint NOT NULL REFERENCES crawl_runs(id) ON DELETE CASCADE,
    source_id     integer NOT NULL REFERENCES sources(id),
    found         integer NOT NULL DEFAULT 0,
    new_count     integer NOT NULL DEFAULT 0,
    error_message text NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_crawl_run_sources_run
    ON crawl_run_sources (run_id);
