-- 0007_qa_round2_fixes.sql — Vá theo nghiệm thu QA vòng 2.

-- P2-8 (vòng 1) và P1-NEW-2: khoảng trống dữ liệu và nhiễu chủ đề phải nhìn
-- thấy được ở mức vận hành, không chỉ nằm trong log.
ALTER TABLE crawl_run_sources
    ADD COLUMN IF NOT EXISTS off_topic integer NOT NULL DEFAULT 0;

-- P1-NEW-4 / P2-NEW-3: công ty con mang tên công ty mẹ phải được khai báo phủ
-- định, nếu không "FPT Retail (FRT)" sẽ gắn luôn FPT.
UPDATE tickers SET negative_aliases = ARRAY['FPT Retail', 'FPT Telecom', 'FPT Shop', 'FPT Software', 'FPT Securities', 'Chứng khoán FPT']
    WHERE symbol = 'FPT';
UPDATE tickers SET negative_aliases = ARRAY['Masan Consumer', 'Masan MEATLife', 'Masan High-Tech', 'Masan High-Tech Materials']
    WHERE symbol = 'MSN';
UPDATE tickers SET negative_aliases = ARRAY['Vinhomes', 'Vincom Retail', 'VinFast', 'Vincom', 'VinBigdata']
    WHERE symbol = 'VIC';
UPDATE tickers SET negative_aliases = ARRAY['Techcom Securities', 'TCBS', 'TechcomInsurance', 'Techcom Life']
    WHERE symbol = 'TCB';
UPDATE tickers SET negative_aliases = ARRAY['VPBankS', 'VPBank Securities', 'VPBank Finance']
    WHERE symbol = 'VPB';
UPDATE tickers SET negative_aliases = ARRAY['ACB Securities', 'ACBS', 'ACB Insurance']
    WHERE symbol = 'ACB';
UPDATE tickers SET negative_aliases = ARRAY['MB Securities', 'MBS', 'MB Capital', 'MBAgeas']
    WHERE symbol = 'MBB';
UPDATE tickers SET negative_aliases = ARRAY['Bảo Việt Bank', 'BVBank', 'Bảo Việt Securities', 'BVSC']
    WHERE symbol = 'BVH';

-- P1-NEW-1: bộ lọc tư vấn cũ bỏ dấu trước khi so khớp nên bắt nhầm "nền tảng"
-- (trùng "nên tăng") và "căn bản" (trùng "cần bán"). Thả các bài bị cách ly oan
-- trở lại publish: chỉ thả bài mà bản lọc ĐÃ SỬA (so trên text có dấu) không
-- còn bắt, và chỉ thả bài bị giữ vì đúng lý do đó.
UPDATE articles
SET status = 'published',
    reject_reason = NULL
WHERE status = 'pending_review'
  AND reject_reason IN ('summary_rejected', 'investment_advice_rescan')
  AND summary_md !~* ('(nên mua|nên bán|nên nắm giữ|nên duy trì|nên gom|nên giải ngân'
      || '|nên chốt lời|nên cắt lỗ|nên mua vào|nên bán ra|nên hạn chế|nên ưu tiên|nên tránh'
      || '|nhà đầu tư nên|khuyến nghị|khuyến cáo nhà đầu tư'
      || '|giá mục tiêu|mục tiêu giá|tránh mua đuổi|mua đuổi'
      || '|canh mua|canh bán|bắt đáy|chốt lời|cắt lỗ|lướt sóng'
      || '|khả quan|kém khả quan|outperform|underperform|overweight|underweight'
      || '|tiềm năng tăng giá|cơ hội đầu tư|đáng để đầu tư'
      || '|giải ngân dần|giải ngân thêm|chủ động giải ngân|vùng giải ngân'
      || '|(nên|cần|hãy) (mua|bán|nắm giữ|duy trì|gom|giải ngân|chốt lời|cắt lỗ)'
      || '|(duy trì|tăng|giảm|nâng|hạ|phân bổ|gia tăng|cắt giảm) tỷ trọng)');

-- P0-3: tầng T3 ngữ cảnh ngành đã bị tắt (env TAGGER_INDUSTRY_TIER=false).
-- Gỡ các tag do tầng này sinh ra: QA đếm toàn bộ 38 tag score=0.40 và cả 38 đều
-- sai. Giữ lại trong DB thì chúng vẫn hiển thị trên cột Mã CK.
DELETE FROM article_tickers WHERE score = 0.40;
