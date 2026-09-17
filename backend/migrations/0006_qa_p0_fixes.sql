-- 0006_qa_p0_fixes.sql — Vá dữ liệu cho các defect P0 trong docs/05-qa/qa-report.md.

-- P0-3 · alias sai: "Điện Máy Xanh" là chuỗi bán lẻ của MWG NHƯNG cũng là tên
-- pháp nhân niêm yết riêng (CTCP Đầu tư Điện Máy Xanh, mã DMX). Alias mơ hồ nên
-- gỡ khỏi MWG và khai báo DMX. PRD ưu tiên precision hơn recall: phân vân thì
-- đừng gắn.
INSERT INTO tickers (symbol, company_name, short_name, exchange, sector, in_vn30) VALUES
    ('DMX', 'CTCP Đầu tư Điện Máy Xanh', 'Điện Máy Xanh', 'HOSE', 'Bán lẻ', false)
ON CONFLICT (symbol) DO NOTHING;

INSERT INTO ticker_aliases (symbol, alias, alias_normalized) VALUES
    ('DMX', 'Điện Máy Xanh', 'dien may xanh')
ON CONFLICT (symbol, alias_normalized) DO NOTHING;

DELETE FROM ticker_aliases WHERE symbol = 'MWG' AND alias_normalized = 'dien may xanh';

UPDATE tickers SET negative_aliases = ARRAY['Điện Máy Xanh', 'Bách Hoá Xanh Online']
    WHERE symbol = 'MWG';

-- P0-3 · trùng ký hiệu: VNM trong "VanEck Vectors Vietnam ETF (VNM ETF)" là mã
-- quỹ Mỹ, không phải Vinamilk.
UPDATE tickers SET negative_aliases = ARRAY['VNM ETF', 'VanEck Vectors Vietnam ETF']
    WHERE symbol = 'VNM';

-- P0-4 · quét lại toàn bộ tin ĐANG PUBLISH bằng danh sách cấm đã mở rộng.
-- Lỗi biên từ cũ đã cho lọt tin thật; không thể để chúng tiếp tục hiển thị
-- trong lúc chờ run ingest kế tiếp.
UPDATE articles
SET status = 'pending_review',
    reject_reason = COALESCE(NULLIF(reject_reason, ''), 'investment_advice_rescan')
WHERE status = 'published'
  AND summary_md ~* ('(nên mua|nên bán|nên nắm giữ|nên duy trì|nên gom|nên giải ngân'
      || '|nên chốt lời|nên cắt lỗ|nên mua vào|nên bán ra|nên tăng|nên giảm'
      || '|nhà đầu tư nên|khuyến nghị|khuyến cáo nhà đầu tư'
      || '|giá mục tiêu|mục tiêu giá|tránh mua đuổi|mua đuổi'
      || '|canh mua|canh bán|bắt đáy|chốt lời|cắt lỗ|lướt sóng'
      || '|khả quan|kém khả quan|outperform|underperform|overweight|underweight'
      || '|tiềm năng tăng giá|cơ hội đầu tư|đáng để đầu tư'
      || '|giải ngân dần|giải ngân thêm|chủ động giải ngân|vùng giải ngân'
      || '|(duy trì|tăng|giảm|nâng|hạ|phân bổ|cắt giảm) tỷ trọng)');

-- P0-2 · các bài đang dùng chung tóm tắt với bài khác cùng nguồn là sản phẩm của
-- lỗi bóc tách, không được để hiển thị.
UPDATE articles a
SET status = 'pending_review',
    reject_reason = COALESCE(NULLIF(a.reject_reason, ''), 'duplicate_summary')
WHERE a.status = 'published'
  AND a.summary_md <> ''
  AND EXISTS (
      SELECT 1 FROM articles b
      WHERE b.source_id = a.source_id
        AND b.summary_md = a.summary_md
        AND b.id <> a.id);
