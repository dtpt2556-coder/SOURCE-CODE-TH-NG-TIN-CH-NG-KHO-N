-- 0002_seed_sources.sql — Danh mục nguồn tin đã duyệt (§1 BA2).
-- Domain trong bảng này đồng thời là ALLOWLIST chống open redirect của /r/{code}.
-- Idempotent qua ON CONFLICT (code).

INSERT INTO sources (code, name, domain, rss_url, list_url, tier, enabled, rate_limit_ms) VALUES
    -- Tier 0 — nguồn sơ cấp (CBTT & cơ quan quản lý), không có RSS công khai.
    ('hsx',  'HOSE — Sở GDCK TP.HCM',           'hsx.vn',     '', 'https://www.hsx.vn/Modules/Cms/Web/NewsByCat/e4b9b0d1-c1b9-4b1e-9d4e-000000000001', 0, false, 3000),
    ('hnx',  'HNX — Sở GDCK Hà Nội',            'hnx.vn',     '', 'https://www.hnx.vn/vi-vn/thong-tin-cong-bo.html', 0, false, 3000),
    ('ssc',  'UBCKNN — Uỷ ban Chứng khoán NN',  'ssc.gov.vn', '', 'https://www.ssc.gov.vn/ubck/faces/vi/vimenu/vipages_vitintuc', 0, false, 3000),

    -- Tier 1 — báo/chuyên trang tài chính uy tín.
    ('cafef',              'CafeF',                  'cafef.vn',              'https://cafef.vn/thi-truong-chung-khoan.rss',                  'https://cafef.vn/thi-truong-chung-khoan.chn',        1, true, 2000),
    ('vietstock',          'Vietstock',              'vietstock.vn',          'https://vietstock.vn/rss',                                     'https://vietstock.vn/chung-khoan.htm',               1, true, 2000),
    ('vnexpress',          'VnExpress',              'vnexpress.net',         'https://vnexpress.net/rss/kinh-doanh/chung-khoan.rss',         'https://vnexpress.net/kinh-doanh/chung-khoan',       1, true, 2000),
    ('tinnhanhchungkhoan', 'Tin nhanh Chứng khoán',  'tinnhanhchungkhoan.vn', 'https://www.tinnhanhchungkhoan.vn/rss/chung-khoan.rss',        'https://www.tinnhanhchungkhoan.vn/chung-khoan/',     1, true, 2000),
    ('baodautu',           'Báo Đầu tư',             'baodautu.vn',           'https://baodautu.vn/rss/tai-chinh-chung-khoan.rss',            'https://baodautu.vn/tai-chinh-chung-khoan/',         1, true, 2000),
    ('vneconomy',          'VnEconomy',              'vneconomy.vn',          'https://vneconomy.vn/chung-khoan.rss',                         'https://vneconomy.vn/chung-khoan.htm',               1, true, 2500),

    -- Tier 2 — báo tổng hợp bổ trợ độ phủ.
    ('thanhnien',          'Thanh Niên',             'thanhnien.vn',          'https://thanhnien.vn/rss/kinh-te/chung-khoan.rss',             'https://thanhnien.vn/kinh-te/chung-khoan.htm',       2, true, 2500),
    ('tuoitre',            'Tuổi Trẻ',               'tuoitre.vn',            'https://tuoitre.vn/rss/kinh-doanh.rss',                        'https://tuoitre.vn/kinh-doanh.htm',                  2, true, 2500)
ON CONFLICT (code) DO UPDATE SET
    name          = EXCLUDED.name,
    domain        = EXCLUDED.domain,
    rss_url       = EXCLUDED.rss_url,
    list_url      = EXCLUDED.list_url,
    tier          = EXCLUDED.tier,
    rate_limit_ms = EXCLUDED.rate_limit_ms;
