INSERT INTO tree_mobile
(name, description, is_active, create_time, creator, geom, tree_attributes) 
VALUES ('MyNewTree', 'Test de création', 't', CURRENT_TIMESTAMP, 999, ST_GeomFromText('POINT(2538221 1152372)', 2056), '{"idvalidation":1, "etatsanitairerem":"Rien à signaler"}');

INSERT INTO tree_mobile
(name, description, is_active, create_time, creator, geom, tree_attributes) 
VALUES ('Mon bel arbre', 'création d''un nouvel arbre', 't', CURRENT_TIMESTAMP, 999, ST_GeomFromText('POINT(2538221 1152372)', 2056), '{"idvalidation":2, "enracinairerem":"En ordre"}');
