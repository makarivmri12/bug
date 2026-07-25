-- Add indexes for performance
CREATE INDEX CONCURRENTLY idx_findings_created_at ON findings(created_at DESC);
CREATE INDEX CONCURRENTLY idx_findings_category_severity ON findings(category, severity);
CREATE INDEX CONCURRENTLY idx_scans_created_at ON scans(created_at DESC);
CREATE INDEX CONCURRENTLY idx_scans_workspace_created ON scans(workspace_id, created_at DESC);
CREATE INDEX CONCURRENTLY idx_users_created_at ON users(created_at DESC);
CREATE INDEX CONCURRENTLY idx_workspace_members_user_id ON workspace_members(user_id);

-- Add comments for documentation
COMMENT ON TABLE findings IS 'Stores vulnerability findings from scans';
COMMENT ON COLUMN findings.evidence IS 'JSON containing request/response evidence';
COMMENT ON TABLE scans IS 'Stores scan metadata and status';
COMMENT ON TABLE workspaces IS 'Stores isolated workspaces for multi-tenancy';
COMMENT ON TABLE users IS 'Stores user accounts and authentication';
