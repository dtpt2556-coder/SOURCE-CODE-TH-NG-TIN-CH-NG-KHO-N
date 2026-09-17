-- 0004_seed_demo.sql — Dữ liệu demo lấy NGUYÊN VĂN từ docs/00-inputs/reference-images.md.
--   * 5 tin của "Màn hình B — News Digest Table" (đúng ngày, mã CK, tóm tắt kèm
--     **bold**, domain nguồn, loại tin).
--   * 1 research note CMG của "Màn hình A" với 4 luận điểm.
-- Mục đích: ngay sau khi khởi động, sản phẩm trông ĐÚNG NHƯ ảnh tham chiếu để
-- QA đối chiếu UI. Idempotent qua ON CONFLICT.

-- ------------------------------------------------------------------ articles
INSERT INTO articles (
    source_id, canonical_url, url_hash, title, summary_md, excerpt,
    news_type, status, summary_provider, published_at, fetched_at)
SELECT
    s.id,
    v.url,
    encode(sha256(convert_to(v.url, 'UTF8')), 'hex'),
    v.title,
    v.summary_md,
    replace(v.summary_md, '**', ''),
    v.news_type,
    'published',
    'extractive',
    v.published_at,
    v.published_at
FROM (VALUES
    (
        'vnexpress',
        'https://vnexpress.net/12-ngan-hang-cam-ket-408000-ty-dong-tin-dung-cho-dnnvv-4790001.html',
        '12 ngân hàng cam kết 408.000 tỷ đồng tín dụng cho DNNVV',
        '12 ngân hàng cam kết **408.000 tỷ đồng** tín dụng cho DNNVV theo chương trình NHNN chủ trì. Big4 đăng ký 220.000 tỷ (Agribank 70.000; BIDV, Vietcombank, VietinBank mỗi bên 50.000 tỷ), lãi suất thấp hơn ít nhất 1 điểm % so bình quân cùng kỳ hạn. 8 ngân hàng tư nhân đăng ký 188.000 tỷ, giảm lãi 0,5-2 điểm % tuỳ lĩnh vực kèm miễn giảm phí. Cuối tháng 7 chỉ 8,8% DNNVV tiếp cận được vốn so với trên 47% ở doanh nghiệp lớn.',
        'nganh',
        '2026-08-26T01:30:00Z'::timestamptz
    ),
    (
        'vnexpress',
        'https://vnexpress.net/cbam-doanh-nghiep-thep-nhom-viet-doi-mat-chi-phi-carbon-khi-xuat-sang-eu-4790002.html',
        'CBAM: Doanh nghiệp thép, nhôm Việt đối mặt chi phí carbon lớn khi xuất sang EU',
        'CBAM của EU áp dụng từ đầu 2026, doanh nghiệp có nguy cơ bù đắp **hơn 100 USD/tấn thép** xuất sang EU, nhôm tới cả nghìn USD/tấn. Giá chứng chỉ carbon EU trên 75 USD/tấn CO2. Lô 100.000 tấn thép đơn giá 500-550 USD/tấn phải nộp khoảng 300 tỷ đồng, tương đương 20% giá trị đơn hàng. Cường độ phát thải thép Việt 2,5 tấn CO2/tấn, nhôm 14 tấn so chuẩn EU 1,4 tấn. Sáu nhóm chịu CBAM: sắt thép, nhôm, xi măng, phân bón, điện, hydrogen.',
        'nganh',
        '2026-08-27T02:15:00Z'::timestamptz
    ),
    (
        'cafef',
        'https://cafef.vn/vn-index-dong-cua-1821-diem-lan-dau-vuot-moc-1800-188260826154500123.chn',
        'VN-Index đóng cửa 1.821 điểm, lần đầu vượt mốc 1.800',
        'VN-Index đóng cửa **1.821 điểm (+29,91 điểm, +1,67%)**, lần đầu vượt 1.800, phiên tăng thứ 5 liên tiếp. Giá trị giao dịch ~20.000 tỷ, vượt bình quân 20 phiên; 182 mã tăng / 125 giảm. TCB trần lên 33.450 đ khớp 39,17 triệu cp, BCM trần 44.450 đ, VIC +4,31%, FPT +2,69%. Khối ngoại mua ròng ~14 tỷ toàn thị trường nhưng bán ròng 43 tỷ trên HOSE: mua FPT 203 tỷ, TCB 136 tỷ, PNJ 100 tỷ; bán CTG 78 tỷ, ACB 78 tỷ, VHM 61 tỷ.',
        'nganh',
        '2026-08-26T08:45:00Z'::timestamptz
    ),
    (
        'thanhnien',
        'https://thanhnien.vn/sau-thang-dau-2026-nhap-494000-tan-thit-tri-gia-15-ty-usd-185260827102000456.htm',
        'Sáu tháng đầu 2026, Việt Nam nhập 494.000 tấn thịt trị giá 1,5 tỷ USD',
        'Sáu tháng đầu 2026 nhập **494.000 tấn thịt, ~1,5 tỷ USD (+10% lượng, +36,4% giá trị)**; cả năm 2025 là 1,95 tỷ USD. Gia cầm 186.400 tấn (đùi gà đông lạnh 130.000 tấn), thịt trâu 108.000 tấn, thịt heo ~80.000 tấn. Quý II nhập 42.600 tấn thịt heo, giá bình quân 2.152 USD/tấn giảm 18,5%. Nguồn chính: Mỹ 63.300 tấn, Thổ Nhĩ Kỳ 29.600 tấn (+261%), Tây Ban Nha 13.300 tấn. Giá heo hơi trong nước 56.000 đ/kg chịu sức ép từ hàng đông lạnh giá rẻ.',
        'nganh',
        '2026-08-27T03:20:00Z'::timestamptz
    ),
    (
        'cafef',
        'https://cafef.vn/acbs-uoc-5588-ty-dong-chay-vao-117-co-phieu-viet-nam-ky-co-cau-ftse-geis-188260827110500789.chn',
        'ACBS: 5.588 tỷ đồng có thể chảy vào 117 cổ phiếu Việt Nam kỳ cơ cấu FTSE GEIS',
        'ACBS ước **5.588 tỷ đồng (216 triệu USD)** chảy vào 117 cổ phiếu Việt Nam kỳ cơ cấu FTSE GEIS hiệu lực 21/9/2026, từ 17 quỹ thụ động; phân bổ 3 large-cap, 3 mid-cap, 21 small-cap, 90 micro-cap. Ba mã hút mạnh nhất: VIC 80,67 triệu USD (2.092 tỷ), VHM 25,43 triệu USD (659 tỷ), HPG 13,11 triệu USD (340 tỷ). Hai quỹ tham chiếu lớn nhất là Vanguard Total International Stock Index Fund 646,2 tỷ USD và Vanguard FTSE Emerging Markets ETF 162,3 tỷ USD.',
        'nganh',
        '2026-08-27T04:05:00Z'::timestamptz
    )
) AS v(src_code, url, title, summary_md, news_type, published_at)
JOIN sources s ON s.code = v.src_code
ON CONFLICT (url_hash) DO NOTHING;

-- --------------------------------------------------------------- short_links
-- target_url thuộc đúng domain nguồn => qua được allowlist chống open redirect.
INSERT INTO short_links (code, article_id, target_url, target_alive)
SELECT v.code, a.id, a.canonical_url, true
FROM (VALUES
    ('a7Kx2p', 'https://vnexpress.net/12-ngan-hang-cam-ket-408000-ty-dong-tin-dung-cho-dnnvv-4790001.html'),
    ('b3Qw9m', 'https://vnexpress.net/cbam-doanh-nghiep-thep-nhom-viet-doi-mat-chi-phi-carbon-khi-xuat-sang-eu-4790002.html'),
    ('c8Zt4r', 'https://cafef.vn/vn-index-dong-cua-1821-diem-lan-dau-vuot-moc-1800-188260826154500123.chn'),
    ('d2Lp6v', 'https://thanhnien.vn/sau-thang-dau-2026-nhap-494000-tan-thit-tri-gia-15-ty-usd-185260827102000456.htm'),
    ('e5Nh1s', 'https://cafef.vn/acbs-uoc-5588-ty-dong-chay-vao-117-co-phieu-viet-nam-ky-co-cau-ftse-geis-188260827110500789.chn')
) AS v(code, url)
JOIN articles a ON a.canonical_url = v.url
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------ article_tickers
-- `score` giảm dần đúng theo thứ tự mã hiển thị trong ảnh tham chiếu
-- (cột "Mã CK"), vì API sắp xếp ORDER BY score DESC, symbol.
INSERT INTO article_tickers (article_id, symbol, relevance, score)
SELECT sl.article_id, v.symbol, 'primary', v.score
FROM (VALUES
    -- Tin 1: BID, VCB, CTG, SHB, MSB, STB, BVB, NAB, NVB, SGB, TPB
    ('a7Kx2p', 'BID', 0.99::real),
    ('a7Kx2p', 'VCB', 0.98),
    ('a7Kx2p', 'CTG', 0.97),
    ('a7Kx2p', 'SHB', 0.96),
    ('a7Kx2p', 'MSB', 0.95),
    ('a7Kx2p', 'STB', 0.94),
    ('a7Kx2p', 'BVB', 0.93),
    ('a7Kx2p', 'NAB', 0.92),
    ('a7Kx2p', 'NVB', 0.91),
    ('a7Kx2p', 'SGB', 0.90),
    ('a7Kx2p', 'TPB', 0.89),
    -- Tin 2: HPG, HSG, NKG, TIS
    ('b3Qw9m', 'HPG', 0.99),
    ('b3Qw9m', 'HSG', 0.98),
    ('b3Qw9m', 'NKG', 0.97),
    ('b3Qw9m', 'TIS', 0.96),
    -- Tin 3: TCB, BCM, VIC, FPT
    ('c8Zt4r', 'TCB', 0.99),
    ('c8Zt4r', 'BCM', 0.98),
    ('c8Zt4r', 'VIC', 0.97),
    ('c8Zt4r', 'FPT', 0.96),
    -- Tin 4: DBC, BAF, MML, HAG
    ('d2Lp6v', 'DBC', 0.99),
    ('d2Lp6v', 'BAF', 0.98),
    ('d2Lp6v', 'MML', 0.97),
    ('d2Lp6v', 'HAG', 0.96),
    -- Tin 5: VIC, VHM, HPG
    ('e5Nh1s', 'VIC', 0.99),
    ('e5Nh1s', 'VHM', 0.98),
    ('e5Nh1s', 'HPG', 0.97)
) AS v(code, symbol, score)
JOIN short_links sl ON sl.code = v.code
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------- research_notes
INSERT INTO research_notes (slug, symbol, title, section_heading, disclaimer, published_at)
VALUES (
    'cmg-quy-2-2026',
    'CMG',
    'CMG – QUÝ 2/2026: DOANH THU TIẾP TỤC TĂNG NHƯNG LỢI NHUẬN CĐ MẸ GIẢM – CHU KỲ ĐẦU TƯ DATA CENTER CHƯA TẠO RA LỢI NHUẬN TƯƠNG XỨNG',
    'LUẬN ĐIỂM ĐẦU TƯ',
    'Nội dung mang tính tổng hợp thông tin, không phải khuyến nghị mua/bán chứng khoán. Nhà đầu tư tự chịu trách nhiệm với quyết định của mình.',
    '2026-08-28T02:00:00Z'::timestamptz
)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO research_note_points (note_id, ordinal, lead, body)
SELECT n.id, v.ordinal, v.lead, v.body
FROM (VALUES
    (
        1,
        'Hoạt động kinh doanh cốt lõi chưa xấu đi, nhưng Q2/2026 cho thấy doanh thu và lợi nhuận bắt đầu đi lệch nhau.',
        'Q2/2026, tương ứng quý đầu tiên của niên độ tài chính 2026–2027 của CMC, doanh thu đạt 2.323 tỷ đồng, tăng 5,1% YoY; biên gộp tăng nhẹ từ 17,7% lên 17,9%. Tuy nhiên, LNST hợp nhất giảm 13,1% còn khoảng 101 tỷ đồng và LNST CĐ mẹ giảm 20,3% còn 75 tỷ đồng. Nguyên nhân chính không nằm ở biên gộp mà ở chi phí vận hành và đặc biệt chi phí lãi vay tăng 48%. Điều này phản ánh giai đoạn CMC phải bỏ vốn trước cho hạ tầng số trong khi doanh thu mới chưa theo kịp.'
    ),
    (
        2,
        'Chất lượng tài sản công nghệ của CMC đáng chú ý hơn tốc độ lợi nhuận hiện tại.',
        'CMC có vị thế tương đối mạnh trong ba thị trường: tích hợp hệ thống và chuyển đổi số cho doanh nghiệp lớn, hạ tầng Cloud/Data Center trong nước và xuất khẩu dịch vụ CNTT. Khách hàng có thể xác định gồm Samsung SDS, IBM, Honda, Panasonic, VPBank, BIDV, VIB, ABBank cùng nhiều cơ quan Chính phủ. Theo MBS Research/CBRE, CMC Telecom chiếm khoảng 12% thị trường Data Center Việt Nam; CMC cũng công bố CMC Cloud chiếm trên 25% thị trường cloud nội địa.'
    ),
    (
        3,
        'Hyperscale Data Center là tài sản có khả năng thay đổi quy mô CMG, nhưng chưa thể coi là lợi nhuận đã chắc chắn.',
        'Dự án tại Khu Công nghệ cao TP.HCM có tổng vốn giai đoạn đầu trên 250 triệu USD, công suất thiết kế ban đầu 30 MW và có khả năng mở rộng trên 100 MW. Nhu cầu AI, Cloud và lưu trữ dữ liệu tạo thị trường thuận lợi, nhưng hiệu quả đầu tư cuối cùng còn phụ thuộc vào tỷ lệ lấp đầy, giá thuê, chi phí điện, cấu trúc vốn và thời gian đưa dự án vào khai thác.'
    ),
    (
        4,
        'Ở 24.150 đồng/cp, định giá đã giảm đáng kể nhưng chưa phải trường hợp “rẻ rõ ràng”.',
        'P/E TTM hiện khoảng 14 lần, thấp hơn đáng kể vùng định giá trước đây của CMG. Tuy nhiên, ROE chỉ quanh 12% và lợi nhuận đang chịu áp lực trong khi CAPEX tăng mạnh. Biên an toàn chỉ thực sự xuất hiện nếu CMC chứng minh được rằng giai đoạn đầu tư hiện nay sẽ chuyển thành tăng trưởng LNST CĐ mẹ từ 2027 trở đi.'
    )
) AS v(ordinal, lead, body)
CROSS JOIN (SELECT id FROM research_notes WHERE slug = 'cmg-quy-2-2026') AS n
ON CONFLICT (note_id, ordinal) DO NOTHING;
