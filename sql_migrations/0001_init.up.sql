CREATE TABLE user_data (
	id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v1(),
	display_name character varying(200) NOT NULL CONSTRAINT displayNameLength CHECK (char_length(display_name) > 0),
	password character varying(100) NOT NULL CONSTRAINT passwordLength CHECK (char_length(password) > 0),
	email character varying(200) UNIQUE NOT NULL CONSTRAINT emailLength CHECK (char_length(email) > 0),
	created_at timestamp DEFAULT NOW() NOT NULL,
	updated_at timestamp DEFAULT NOW() NOT NULL,
	deleted_at timestamp
);
