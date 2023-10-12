CREATE TABLE IF NOT EXISTS persons(
                      id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                      email varchar(255) UNIQUE NOT NULL,
                      phone varchar(255) NOT NULL,
                      first_name varchar(255) NOT NULL,
                      created_at timestamp,
                      updated_at timestamp
);