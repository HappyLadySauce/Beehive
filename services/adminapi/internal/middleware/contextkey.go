package middleware

type contextKey string

const AdminUserIDKey contextKey = "admin_user_id"
const CodeUnauth = 1001
const CodeForbidden = 1003
const CodeInternal = 5000
