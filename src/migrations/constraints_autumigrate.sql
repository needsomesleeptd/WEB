CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    page_count INT NOT NULL,
    document_name VARCHAR(255) NOT NULL,
    checks_count INT NOT NULL,
    creator_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS markup_types (
    id serial  PRIMARY KEY,
    description VARCHAR(1000) NOT NULL,
    creator_id INT NOT NULL,
    class_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
   id serial PRIMARY KEY,
   login VARCHAR(255) UNIQUE NOT NULL,
   password VARCHAR(255) NOT NULL,
   name VARCHAR(255) NOT NULL,
   surname VARCHAR(255) NOT NULL,
   role INT NOT NULL,
   group VARCHAR(255) NOT NULL 
);

CREATE TABLE IF NOT EXISTS markups (
    id serial PRIMARY KEY NOT NULL,
    page_data BYTEA NOT NULL,
    error_bb JSONB DEFAULT '[]' NOT NULL,
    class_label BIGINT NOT NULL,
    creator_id BIGINT NOT NULL,
);



CREATE TABLE IF NOT EXISTS document_queues (
    id serial PRIMARY KEY NOT NULL,
    doc_id UUID NOT NULL,
    status SMALLINT NOT NULL,
    CONSTRAINT fk_doc_queue_docs FOREIGN KEY (doc_id) REFERENCES documents(id) ON DELETE CASCADE
);




CREATE TABLE IF NOT EXISTS comments (
    id serial PRIMARY KEY NOT NULL,
    doc_id UUID NOT NULL,
    description text NOT NULL,
    creator_id INT NOT NULL,
    CONSTRAINT fk_comment_docs FOREIGN KEY (doc_id) REFERENCES documents(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_users FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS achievments (
    id serial PRIMARY KEY NOT NULL,
    description text NOT NULL,
    creator_id INT NOT NULL,
    granted_to_id INT NOT NULL,
    page_data BYTEA NOT NULL,
    CONSTRAINT fk_comment_creator FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_granted FOREIGN KEY (granted_to_id) REFERENCES users(id) ON DELETE CASCADE
);


ALTER TABLE markups ADD CONSTRAINT fk_markup_markup_type FOREIGN KEY ( class_label ) REFERENCES markup_types( id );

ALTER TABLE documents ADD CONSTRAINT fk_document_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE markups ADD CONSTRAINT fk_markup_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE markup_types ADD CONSTRAINT fk_markup_types_user FOREIGN KEY ( creator_id ) REFERENCES users( id );


--admin
INSERT INTO users (id, login, password, name, surname, role, "group")
VALUES (1, 'admin', '$2a$10$nV.DqaVAtAr9EhCRqseU6OikgPC1GCIYsmb3Enh5pGwTwa/VntK8K', '<string>', '<string>', 2, '<string>');




-- graphs
INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (191,'Отсутствует легенда на графике',1,'no_graph_leg');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (192,'Отсутствует подпись к осям на графике',1,'no_graph_annot');


-- schemes


INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (161,'Неверное расположение стрелок на графиках',1,'wrong_scheme_arrows');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (162,'Отсутствует подпись (да, нет) к блоку ветвления',1,'wrong_scheme_if');


INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (163,'Неверный формат терминаторов схемы алгоритма',1,'wrong_terminators');

-- tables

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (171,'Отсутствует подпись таблицы',1,'no_table_annot');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (172,'Подпись таблицы неверна',1,'wrong_table_annot');

-- formulas 
INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (151,'Выравнивание формулы неверно',1,'wrong_formula_pos');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (152,'Отсутствует/неверна нумерация формулы',1,'wrong_formula_num');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (153,'В формуле неверно сопоставлены скобки (не бьются по 2)',1,'wrong_formula_bounds');


INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (154,'Отсутствует знак препинания в конце формулы',1,'wrong_formula_ending');

INSERT INTO markup_types (id,description,creator_id,class_name)
VALUES (0,'Ошибок нет, все хорошо))',1,'no_errors');



--autoinc MarkupTypes

SELECT pg_get_serial_sequence('markup_types', 'id'); -- Get the sequence 

ALTER SEQUENCE markup_types_id_seq RESTART WITH 200; -- Set starting value to 1000