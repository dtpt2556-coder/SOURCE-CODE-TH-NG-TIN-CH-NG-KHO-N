package domain

import "time"

// ResearchNote là bài phân tích dạng "Luận điểm đầu tư" (Màn hình A).
// Theo BA2 §5.6, loại nội dung này KHÔNG auto-publish trong MVP.
type ResearchNote struct {
	ID             int64
	Slug           string
	Symbol         string
	Title          string
	SectionHeading string
	Disclaimer     string
	PublishedAt    time.Time

	Points []ResearchNotePoint
}

// PublishedDateVN định dạng DD/MM/YYYY theo múi giờ truyền vào.
func (n ResearchNote) PublishedDateVN(loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}
	return n.PublishedAt.In(loc).Format("02/01/2006")
}

// ResearchNotePoint là một luận điểm: dòng dẫn in đậm + đoạn phân tích.
type ResearchNotePoint struct {
	ID      int64
	NoteID  int64
	Ordinal int
	Lead    string
	Body    string
}
