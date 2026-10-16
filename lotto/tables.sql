CREATE DATABASE lotto;

\c lotto

CREATE TABLE IF not exists eurojackpot (
    id      SERIAL NOT NULL PRIMARY KEY,
    one     TEXT NOT NULL,
    two     TEXT NOT NULL,
    three   TEXT NOT NULL,
    four    TEXT NOT NULL,
    five    TEXT NOT NULL,
    sup1    TEXT NOT NULL,
    sup2    TEXT NOT NULL,
    datum   TEXT NOT NULL,
    join_date DATE NOT NULL DEFAULT CURRENT_DATE 
);

CREATE TABLE IF not exists lotto6aus49 (
    id      SERIAL NOT NULL PRIMARY KEY,
    one     TEXT NOT NULL,
    two     TEXT NOT NULL,
    three   TEXT NOT NULL,
    four    TEXT NOT NULL,
    five    TEXT NOT NULL,
    six     TEXT NOT NULL,
    sup1    TEXT NOT NULL,
    datum   TEXT NOT NULL,
    join_date DATE NOT NULL DEFAULT CURRENT_DATE 
);