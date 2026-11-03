package memo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSearchRealTitleAndBody(t *testing.T) {
	s := testStore(t)
	for p, body := range map[string]string{"项目/产品计划.md": "# 产品计划\n\n只记录真实资料。", "研究/中文.md": "# 研究记录\n\n" + strings.Repeat("前文段落", 50) + "本地资料库支持 UTF-8 和中文全文检索。"} {
		if _, e := s.Write(p, body, "", true); e != nil {
			t.Fatal(e)
		}
	}
	for query, path := range map[string]string{"产品": "项目/产品计划.md", "全文检索": "研究/中文.md", "utf-8": "研究/中文.md"} {
		rows, e := s.Search(query)
		if e != nil || len(rows) != 1 || rows[0].Path != path {
			t.Fatalf("query %q: %v %v", query, rows, e)
		}
		if rows[0].Content != "" || !utf8.ValidString(rows[0].Preview) {
			t.Fatal("search must return safe summaries rather than full bodies")
		}
	}
	rows, e := s.Search("不存在的关键词")
	if e != nil || len(rows) != 0 {
		t.Fatal("false positive")
	}
	rows, e = s.Search("")
	if e != nil || len(rows) != 2 {
		t.Fatal("empty query should list real notes")
	}
}
func TestSearchEndpointAuthorization(t *testing.T) {
	s, ts := testServer(t)
	if _, e := s.store.Write("a.md", "# 示例\n\n中文搜索内容", "", true); e != nil {
		t.Fatal(e)
	}
	r := request(t, ts, "GET", "/api/search?q="+url.QueryEscape("搜索"), "", nil, nil)
	if r.StatusCode != http.StatusUnauthorized {
		t.Fatal("unauthorized search accepted")
	}
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	r = request(t, ts, "GET", "/api/search?q="+url.QueryEscape("搜索"), "", h, nil)
	var rows []Note
	if e := json.NewDecoder(r.Body).Decode(&rows); e != nil || r.StatusCode != 200 || len(rows) != 1 {
		t.Fatal("search response invalid")
	}
	r = request(t, ts, "GET", "/api/search?q="+strings.Repeat("a", 513), "", h, nil)
	if r.StatusCode != 400 {
		t.Fatal("oversized query accepted")
	}
}
