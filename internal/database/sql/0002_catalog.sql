CREATE TABLE canteens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE,
    sort_order integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE floors (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    canteen_id uuid NOT NULL REFERENCES canteens(id) ON DELETE CASCADE,
    name text NOT NULL,
    sort_order integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (canteen_id, name)
);

CREATE TABLE food_windows (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    floor_id uuid NOT NULL REFERENCES floors(id) ON DELETE RESTRICT,
    external_code text NOT NULL UNIQUE,
    name text NOT NULL,
    description text NOT NULL DEFAULT '',
    business_hours text NOT NULL DEFAULT '',
    active boolean NOT NULL DEFAULT true,
    sort_order integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX food_windows_floor_id_idx ON food_windows (floor_id);

CREATE TABLE operating_periods (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    window_id uuid NOT NULL REFERENCES food_windows(id) ON DELETE RESTRICT,
    operator_name text NOT NULL,
    started_at timestamptz NOT NULL,
    ended_at timestamptz,
    created_by uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (ended_at IS NULL OR ended_at > started_at)
);

CREATE UNIQUE INDEX operating_periods_one_current_per_window
    ON operating_periods (window_id) WHERE ended_at IS NULL;

CREATE TABLE departments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE department_members (
    department_id uuid NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (department_id, user_id)
);

CREATE TABLE window_department_assignments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    window_id uuid NOT NULL REFERENCES food_windows(id) ON DELETE RESTRICT,
    department_id uuid NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    started_at timestamptz NOT NULL,
    ended_at timestamptz,
    created_by uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (ended_at IS NULL OR ended_at > started_at)
);

CREATE UNIQUE INDEX window_department_one_current_per_window
    ON window_department_assignments (window_id) WHERE ended_at IS NULL;
