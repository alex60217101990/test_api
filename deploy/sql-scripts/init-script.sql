CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS users (
	id INT GENERATED ALWAYS AS IDENTITY,
	public_id UUID NOT NULL DEFAULT gen_random_uuid(),
  	username varchar(45) NOT NULL,
	email varchar(100) NOT NULL,
  	password varchar(450) NOT NULL,
  	is_online bool NOT NULL DEFAULT false,
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at TIME ZONE 'utc'),
	updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at TIME ZONE 'utc'),
	deleted_at TIMESTAMP WITHOUT TIME ZONE,
	PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS users_brin_idx  ON users USING brin (public_id);
CREATE INDEX IF NOT EXISTS username_idx ON users(username);
CREATE INDEX IF NOT EXISTS email_idx ON users(email);

CREATE TABLE IF NOT EXISTS product_categories(
   	id INT GENERATED ALWAYS AS IDENTITY,
	public_id UUID NOT NULL DEFAULT gen_random_uuid(),
	change_by_user INT,
   	name VARCHAR(100) NOT NULL,
	popularity BIGINT,
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at TIME ZONE 'utc'),
	updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at TIME ZONE 'utc'),
	deleted_at TIMESTAMP WITHOUT TIME ZONE,
   	PRIMARY KEY(id),
	CONSTRAINT fk_user_change_categorie
      	FOREIGN KEY(change_by_user) 
	  		REFERENCES users(id)
			ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS name_idx ON product_categories(name);
CREATE INDEX IF NOT EXISTS product_categories_brin_idx  ON product_categories USING brin (public_id);

CREATE TABLE IF NOT EXISTS products(
   	id INT GENERATED ALWAYS AS IDENTITY,
	public_id UUID NOT NULL DEFAULT gen_random_uuid(),
	change_by_user INT,
	categorie_id INT,
   	name VARCHAR(100) NOT NULL,
   	popularity BIGINT,
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at TIME ZONE 'utc'),
	updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP at TIME ZONE 'utc'),
	deleted_at TIMESTAMP WITHOUT TIME ZONE,
   	PRIMARY KEY(id),
   	CONSTRAINT fk_categorie
      	FOREIGN KEY(categorie_id) 
	  		REFERENCES product_categories(id)
			ON DELETE SET NULL,
	CONSTRAINT fk_user_change_product
      	FOREIGN KEY(change_by_user) 
	  		REFERENCES users(id)
			ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS name_idx ON product_categories(name);
CREATE INDEX IF NOT EXISTS product_categories_brin_idx  ON products USING brin (public_id);