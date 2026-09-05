package devseed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type account struct {
	id       string
	email    string
	nickname string
	role     string
}

var accounts = []account{
	{id: "00000000-0000-0000-0000-000000000001", email: "student@mail.ustc.edu.cn", nickname: "开发同学", role: "member"},
	{id: "00000000-0000-0000-0000-000000000002", email: "manager@ustc.edu.cn", nickname: "开发管理人员", role: "manager"},
	{id: "00000000-0000-0000-0000-000000000003", email: "admin@ustc.edu.cn", nickname: "开发平台管理员", role: "admin"},
}

func Run(ctx context.Context, pool *pgxpool.Pool, password string) error {
	if password == "" {
		return fmt.Errorf("DEV_SEED_PASSWORD is required")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash development password: %w", err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, item := range accounts {
		if _, err := tx.Exec(ctx, `
			INSERT INTO users (id, email, password_hash, nickname, role, email_verified_at)
			VALUES ($1, $2, $3, $4, $5, now())
			ON CONFLICT (id) DO UPDATE SET
				email = EXCLUDED.email,
				password_hash = EXCLUDED.password_hash,
				nickname = EXCLUDED.nickname,
				role = EXCLUDED.role,
				email_verified_at = EXCLUDED.email_verified_at,
				disabled_at = NULL
		`, item.id, item.email, string(passwordHash), item.nickname, item.role); err != nil {
			return fmt.Errorf("seed account %s: %w", item.email, err)
		}
	}

	statements := []string{
		`INSERT INTO canteens (id, name, sort_order)
		 VALUES ('10000000-0000-0000-0000-000000000001', '示例食堂', 1)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, sort_order = EXCLUDED.sort_order`,
		`INSERT INTO floors (id, canteen_id, name, sort_order)
		 VALUES ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', '一层', 1),
		        ('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', '二层', 2)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, sort_order = EXCLUDED.sort_order`,
		`INSERT INTO food_windows (id, floor_id, external_code, name, description, business_hours, sort_order)
		 VALUES ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'DEMO-01', '示例窗口一', '用于开发食堂目录和评价功能', '10:30-13:30 / 16:30-19:30', 1),
		        ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000002', 'DEMO-02', '示例窗口二', '用于开发管理范围和流水导入', '10:30-13:30 / 16:30-19:30', 1)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, business_hours = EXCLUDED.business_hours, sort_order = EXCLUDED.sort_order`,
		`INSERT INTO departments (id, name)
		 VALUES ('50000000-0000-0000-0000-000000000001', '示例管理部门')
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`,
		`INSERT INTO department_members (department_id, user_id)
		 VALUES ('50000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002')
		 ON CONFLICT DO NOTHING`,
		`INSERT INTO operating_periods (id, window_id, operator_name, started_at, created_by)
		 VALUES ('40000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', '示例经营方', '2026-01-01T00:00:00+08:00', '00000000-0000-0000-0000-000000000003'),
		        ('40000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', '示例经营方', '2026-01-01T00:00:00+08:00', '00000000-0000-0000-0000-000000000003')
		 ON CONFLICT (id) DO UPDATE SET operator_name = EXCLUDED.operator_name`,
		`INSERT INTO window_department_assignments (id, window_id, department_id, started_at, created_by)
		 VALUES ('60000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', '2026-01-01T00:00:00+08:00', '00000000-0000-0000-0000-000000000003'),
		        ('60000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000001', '2026-01-01T00:00:00+08:00', '00000000-0000-0000-0000-000000000003')
		 ON CONFLICT (id) DO NOTHING`,
	}
	for _, statement := range statements {
		if _, err := tx.Exec(ctx, statement); err != nil {
			return fmt.Errorf("seed catalog: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}
	return nil
}
