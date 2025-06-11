package schema

import "net/http"

// CreateHTTPResponseRFC9457 创建符合RFC9457格式的HTTP响应。
//
// 参数:
//   - title (string): 响应的标题。
//   - status (int): HTTP状态码。
//   - r (*http.Request): 与响应关联的HTTP请求。
//
// 返回:
//   - map[string]interface{}: HTTP响应的映射。
func CreateHTTPResponseRFC9457(title string, status int, r *http.Request) map[string]any {
	return map[string]interface{}{
		"type":     "about:blank",
		"title":    title,
		"status":   status,
		"detail":   "",
		"instance": r.Method + " " + r.RequestURI,
	}
}
