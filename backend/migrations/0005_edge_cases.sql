-- 0005_edge_cases.sql — Cột và dữ liệu cho docs/03-sa/ingest-edge-cases.md.
-- Không sửa 0001–0004 vì các file đó đã chạy ở môi trường khác.

ALTER TABLE articles
    ADD COLUMN IF NOT EXISTS is_demo                boolean     NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS content_hash           text,
    ADD COLUMN IF NOT EXISTS revision               integer     NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS updated_at             timestamptz,
    ADD COLUMN IF NOT EXISTS last_seen_at           timestamptz,
    ADD COLUMN IF NOT EXISTS published_at_estimated boolean     NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS is_live                boolean     NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS reject_reason          text,
    ADD COLUMN IF NOT EXISTS duplicate_of           bigint REFERENCES articles(id),
    -- U-07: bản tóm tắt mới chưa qua validator nằm ở đây, KHÔNG ghi đè summary_md
    -- đang publish. Cột này ngoài §12 nhưng không có nó thì U-07 buộc phải hoặc
    -- vứt bản mới, hoặc publish nội dung chưa kiểm chứng.
    ADD COLUMN IF NOT EXISTS pending_summary_md     text;

ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_status_check;
ALTER TABLE articles ADD CONSTRAINT articles_status_check
    CHECK (status IN ('published', 'pending_review', 'rejected', 'source_gone'));

-- PRD mục 5 chốt 9 nhãn; architecture.md mục 6 chỉ liệt kê 7. Bổ sung 2 nhãn
-- còn thiếu để FE dựng dropdown lọc từ API thay vì hardcode.
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_news_type_check;
ALTER TABLE articles ADD CONSTRAINT articles_news_type_check
    CHECK (news_type IN (
        'vi_mo', 'nganh', 'doanh_nghiep', 'thi_truong', 'khoi_ngoai',
        'co_tuc_phat_hanh', 'phap_ly', 'phan_tich', 'trai_phieu_tin_dung'));

CREATE INDEX IF NOT EXISTS idx_articles_real_published
    ON articles (published_at DESC) WHERE status = 'published' AND is_demo = false;
CREATE INDEX IF NOT EXISTS idx_articles_content_hash
    ON articles (source_id, content_hash) WHERE content_hash IS NOT NULL;

ALTER TABLE sources
    ADD COLUMN IF NOT EXISTS list_selector          text,
    ADD COLUMN IF NOT EXISTS etag                   text,
    ADD COLUMN IF NOT EXISTS last_modified          text,
    ADD COLUMN IF NOT EXISTS consecutive_empty_runs integer     NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS backoff_until          timestamptz,
    ADD COLUMN IF NOT EXISTS notes                  text        NOT NULL DEFAULT '';

ALTER TABLE tickers
    ADD COLUMN IF NOT EXISTS delisted_at            date,
    ADD COLUMN IF NOT EXISTS negative_aliases       text[];

ALTER TABLE crawl_run_sources
    ADD COLUMN IF NOT EXISTS truncated              boolean     NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS skipped_by_robots      integer     NOT NULL DEFAULT 0;

-- G-04: chuỗi KHÔNG được map vào mã, chặn "Vinhomes" -> VIC.
UPDATE tickers SET negative_aliases = ARRAY['Vinhomes', 'Vincom Retail', 'VinFast', 'Vincom'] WHERE symbol = 'VIC';
UPDATE tickers SET negative_aliases = ARRAY['Masan Consumer', 'Masan MEATLife', 'Masan High-Tech'] WHERE symbol = 'MSN';
UPDATE tickers SET negative_aliases = ARRAY['FPT Retail', 'FPT Telecom', 'FPT Shop'] WHERE symbol = 'FPT';
UPDATE tickers SET negative_aliases = ARRAY['Bảo Việt Bank', 'BVBank'] WHERE symbol = 'BVH';
UPDATE tickers SET negative_aliases = ARRAY['Đất Xanh Services'] WHERE symbol = 'DXG';

-- URL RSS đã được probe thực tế ngày 16/09/2026: giữ bật nguồn nào trả feed hợp
-- lệ có item, tắt nguồn nào không. Không để nguồn hỏng ở trạng thái bật.
UPDATE sources SET rss_url = 'https://vnexpress.net/rss/kinh-doanh.rss', enabled = true,
    notes = 'Verified 16/09/2026: 60 item. Duong /rss/kinh-doanh/chung-khoan.rss trong BA2 tra ve HTML, khong phai feed.'
    WHERE code = 'vnexpress';
UPDATE sources SET rss_url = 'https://vietstock.vn/145/chung-khoan.rss', enabled = true,
    notes = 'Verified 16/09/2026: 30 item. Duong /rss trong BA2 la trang index HTML liet ke feed con.'
    WHERE code = 'vietstock';
UPDATE sources SET rss_url = 'https://www.tinnhanhchungkhoan.vn/rss/chung-khoan-124.rss', enabled = true,
    notes = 'Verified 16/09/2026: 50 item. Duong /rss/chung-khoan.rss trong BA2 tra ve 404.'
    WHERE code = 'tinnhanhchungkhoan';
UPDATE sources SET enabled = true, notes = 'Verified 16/09/2026: 50 item.'
    WHERE code IN ('cafef', 'vneconomy', 'thanhnien', 'tuoitre');
UPDATE sources SET enabled = false,
    notes = 'Verified 16/09/2026: feed tra 200 va XML hop le nhung channel RONG (0 item) o phia nguon. Bat lai sau khi nguon sua.'
    WHERE code = 'baodautu';
UPDATE sources SET enabled = false,
    notes = 'Tier 0 khong co RSS cong khai, can crawler HTML rieng (D-02). Chua implement o v1.'
    WHERE code IN ('hsx', 'hnx', 'ssc');

-- T2 mục 0: 5 bài seed từ ảnh tham chiếu dùng URL đúng định dạng nhưng KHÔNG
-- tồn tại. Đánh dấu là demo và tắt short link để không ai bấm vào URL bịa.
UPDATE articles SET is_demo = true, last_seen_at = fetched_at
    WHERE canonical_url IN (
        'https://vnexpress.net/12-ngan-hang-cam-ket-408000-ty-dong-tin-dung-cho-dnnvv-4790001.html',
        'https://vnexpress.net/cbam-doanh-nghiep-thep-nhom-viet-doi-mat-chi-phi-carbon-khi-xuat-sang-eu-4790002.html',
        'https://cafef.vn/vn-index-dong-cua-1821-diem-lan-dau-vuot-moc-1800-188260826154500123.chn',
        'https://thanhnien.vn/sau-thang-dau-2026-nhap-494000-tan-thit-tri-gia-15-ty-usd-185260827102000456.htm',
        'https://cafef.vn/acbs-uoc-5588-ty-dong-chay-vao-117-co-phieu-viet-nam-ky-co-cau-ftse-geis-188260827110500789.chn'
    );

UPDATE short_links SET target_alive = false, last_checked_at = now()
    WHERE article_id IN (SELECT id FROM articles WHERE is_demo = true);
