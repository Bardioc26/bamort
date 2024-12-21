CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);


CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE characters (
    character_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(100),
    rasse VARCHAR(50),
    typ VARCHAR(50),
    alter INT,
    anrede VARCHAR(20),
    grad INT,
    groesse INT,
    gewicht INT,
    glaube VARCHAR(100),
    hand VARCHAR(20),
    image VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE eigenschaften (
    eigenschaften_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    au INT,
    gs INT,
    gw INT,
    `in` INT,
    ko INT,
    pa INT,
    st INT,
    wk INT,
    zt INT,
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE ausruestung (
    ausruestung_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    name VARCHAR(100),
    anzahl INT,
    gewicht FLOAT,
    wert FLOAT,
    beinhaltet_in VARCHAR(255),
    beschreibung TEXT,
    bonus INT DEFAULT NULL,
    ist_magisch BOOLEAN,
    abw INT,
    ausgebrannt BOOLEAN,
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE behaeltnisse (
    behaeltnisse_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    name VARCHAR(100),
    gewicht FLOAT,
    volumen FLOAT,
    tragkraft FLOAT,
    wert FLOAT,
    beinhaltet_in VARCHAR(255),
    beschreibung TEXT,
    ist_magisch BOOLEAN,
    abw INT,
    ausgebrannt BOOLEAN,
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE fertigkeiten (
    fertigkeiten_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    name VARCHAR(100),
    fertigkeitswert INT,
    pp INT DEFAULT NULL,
    bonus INT DEFAULT NULL,
    beschreibung TEXT,
    quelle VARCHAR(100),
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE waffen (
    waffen_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    name VARCHAR(100),
    gewicht FLOAT,
    wert FLOAT,
    anzahl INT,
    beinhaltet_in VARCHAR(255),
    beschreibung TEXT,
    abwb INT,
    anb INT,
    schb INT,
    name_fuer_spezialisierung VARCHAR(100),
    ist_magisch BOOLEAN,
    abw INT,
    ausgebrannt BOOLEAN,
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE zauber (
    zauber_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    name VARCHAR(100),
    bonus INT DEFAULT NULL,
    beschreibung TEXT,
    quelle VARCHAR(100),
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE bennies (
    bennies_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    gg INT,
    gp INT,
    sg INT,
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE transportmittel (
    transportmittel_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    name VARCHAR(100),
    gewicht FLOAT,
    wert FLOAT,
    tragkraft FLOAT,
    beinhaltet_in VARCHAR(255),
    beschreibung TEXT,
    ist_magisch BOOLEAN,
    abw INT,
    ausgebrannt BOOLEAN,
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);

CREATE TABLE erfahrungsschatz (
    erfahrungsschatz_id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    value INT,
    FOREIGN KEY (character_id) REFERENCES characters(character_id)
);


/*
Summary

This schema reflects the modular structure of the CharType definition. Each substructure is normalized into its own table, ensuring:

    Scalability.
    Clear relationships using foreign keys.
    Multi-user handling through the users table.
*/