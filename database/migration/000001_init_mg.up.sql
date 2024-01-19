CREATE OR REPLACE FUNCTION update_modified_column()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';


CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_users_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE pirgs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    owner_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);
CREATE TRIGGER update_pirgs_modtime BEFORE UPDATE ON pirgs FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE pirgs_users (
    id SERIAL PRIMARY KEY,
    pirg_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (pirg_id) REFERENCES pirgs(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TRIGGER update_pirgs_users_modtime BEFORE UPDATE ON pirgs_users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE pirgs_admins (
    id SERIAL PRIMARY KEY,
    pirg_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (pirg_id) REFERENCES pirgs(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TRIGGER update_pirgs_admins_modtime BEFORE UPDATE ON pirgs_admins FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE pirgs_groups (
    id SERIAL PRIMARY KEY,
    pirg_id INT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (pirg_id) REFERENCES pirgs(id)
);
CREATE TRIGGER update_pirgs_groups_modtime BEFORE UPDATE ON pirgs_groups FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE groups_users (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (group_id) REFERENCES pirgs_groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TRIGGER update_groups_users_modtime BEFORE UPDATE ON groups_users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE api_keys (
    key TEXT PRIMARY KEY,
	role TEXT NOT NULL,
	user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TRIGGER update_api_keys_modtime BEFORE UPDATE ON api_keys FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
