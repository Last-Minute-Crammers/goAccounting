package initialize

import "strings"

// IsAdminEmail 检查邮箱是否在管理员列表中
func (a *_admin) IsAdminEmail(email string) bool {
    if a.Emails == "" {
        return false
    }
    
    // 分割并清理邮箱列表
    emails := strings.Split(a.Emails, ",")
    for _, adminEmail := range emails {
        trimmedEmail := strings.TrimSpace(adminEmail)
        if strings.EqualFold(trimmedEmail, strings.TrimSpace(email)) {
            return true
        }
    }
    return false
}

// GetAdminEmails 获取所有管理员邮箱列表
func (a *_admin) GetAdminEmails() []string {
    if a.Emails == "" {
        return []string{}
    }
    
    emails := strings.Split(a.Emails, ",")
    result := make([]string, 0, len(emails))
    for _, email := range emails {
        trimmed := strings.TrimSpace(email)
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}