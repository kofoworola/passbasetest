CREATE TABLE IF NOT EXISTS project (
	project_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name varchar(100) NOT NULL,
	api_key varchar UNIQUE,
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now()
)
