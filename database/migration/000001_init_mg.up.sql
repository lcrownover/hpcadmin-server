CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE pirgs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    owner_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE pirgs_users (
    id SERIAL PRIMARY KEY,
    pirg_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (pirg_id) REFERENCES pirgs(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE pirgs_admins (
    id SERIAL PRIMARY KEY,
    pirg_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (pirg_id) REFERENCES pirgs(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE pirgs_groups (
    id SERIAL PRIMARY KEY,
    pirg_id INT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (pirg_id) REFERENCES pirgs(id)
);

CREATE TABLE groups_users (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (group_id) REFERENCES pirgs_groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE api_keys (
    key TEXT PRIMARY KEY,
	role TEXT NOT NULL,
	user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
