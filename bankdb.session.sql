CREATE TABLE managers 
(
      id BIGSERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      salary INTEGER NOT NULL CHECK (salary > 0),
      plan INTEGER NOT NULL CHECK (plan > 0),
      boss_id BIGINIT REFERENCES mana
);