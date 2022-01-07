DROP TABLE IF EXISTS characters;
CREATE TABLE characters (
  id         INT AUTO_INCREMENT NOT NULL,
  name      VARCHAR(128) NOT NULL,
  role     VARCHAR(255) NOT NULL,
  level      INT NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO characters
  (name, role, level)
VALUES
	('Hades', 'Emet-Selch', 99),
	('Venat', 'Former Azem', 99),
	('Hythlodaeus', 'Chief of the Bureau of the Architect', 99),
	('Thancred', 'Scion', 80),
	("Y'shtola", 'Scion', 80),
	('Urianger', 'Scion', 80),
	('Lyse', 'Ala Mhigan Resistance', 70);
