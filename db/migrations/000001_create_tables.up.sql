CREATE TABLE IF NOT EXISTS first_names(
    id serial PRIMARY KEY,
    name VARCHAR (50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS people(
    id serial PRIMARY KEY,
    p_name VARCHAR (50) NOT NULL,
    surname VARCHAR (50) NOT NULL,
    patronymic VARCHAR (50),
    age INT,
    gender VARCHAR (20),
    FOREIGN KEY (p_name) REFERENCES first_names(name)
);

CREATE TABLE IF NOT EXISTS nation(
    rec_id serial PRIMARY KEY,
    user_name VARCHAR (50) NOT NULL,
    country_id VARCHAR (10) NOT NULL,
    probability REAL,
    FOREIGN KEY (user_name) REFERENCES first_names(name)
);