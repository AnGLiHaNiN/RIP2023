-- Adminer 4.8.1 PostgreSQL 16.1 (Debian 16.1-1.pgdg120+1) dump

DROP TABLE IF EXISTS "components";
CREATE TABLE "public"."components" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "name" character varying(100) NOT NULL,
    "image_url" character varying(100),
    "world_name" character varying(75) NOT NULL,
    "amount" bigint NOT NULL,
    "properties" character varying(200) NOT NULL,
    "is_deleted" boolean DEFAULT false NOT NULL,
    CONSTRAINT "components_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "components" ("uuid", "name", "image_url", "world_name", "amount", "properties", "is_deleted") VALUES
('cb32e347-2af0-4339-b70f-1550d943c382',	'Клавулановая кислота',	'http://localhost:9000/images/cb32e347-2af0-4339-b70f-1550d943c382.jpg',	'Clavulanic acid',	314,	'Cтимулирует антимикробный иммунитет, способствовауя лизису бактериальной стенки',	'f'),
('e8a91e96-6d86-49fd-a17d-4801a1c78642',	'Амоксициллин',	'http://localhost:9000/images/e8a91e96-6d86-49fd-a17d-4801a1c78642.jpg',	'Amoxicillin',	230,	'полусинтетический антибиотик широкого спектра действия для лечения бактериальных инфекции',	'f'),
('45ce0251-3a3e-41c4-9173-3be471c6b590',	'Ибупрофен',	'http://localhost:9000/images/45ce0251-3a3e-41c4-9173-3be471c6b590.jpg',	'Ibuprofenum',	500,	'болеутоляющие и жаропонижающие свойства',	'f'),
('d56c9e29-26a9-4420-a14c-0538d520f73e',	'Хлоропирамин',	'http://localhost:9000/images/d56c9e29-26a9-4420-a14c-0538d520f73e.jpg',	'Chloropyramine',	980,	'Антигистаминный препарат, блокатор H1-гистаминовых рецепторов I-поколения',	'f'),
('5771c2d7-1b00-4b44-9978-0493e1d778d7',	'Мельдоний',	'http://localhost:9000/images/5771c2d7-1b00-4b44-9978-0493e1d778d7.jpg',	'Meldonium',	845,	'Нормализует энергетический метаболизм клеток, поддерживает энергетический метаболизм сердца и других органов',	'f');

DROP TABLE IF EXISTS "medicine_productions";
CREATE TABLE "public"."medicine_productions" (
    "medicine_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "component_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "count" bigint DEFAULT '1' NOT NULL,
    CONSTRAINT "medicine_productions_pkey" PRIMARY KEY ("medicine_id", "component_id")
) WITH (oids = false);

INSERT INTO "medicine_productions" ("medicine_id", "component_id", "count") VALUES
('913afffa-cc3e-4980-98d1-74df9df53882',	'cb32e347-2af0-4339-b70f-1550d943c382',	1),
('94b59b78-fa0f-48ca-ab3e-7ccd840d6c29',	'cb32e347-2af0-4339-b70f-1550d943c382',	1),
('94b59b78-fa0f-48ca-ab3e-7ccd840d6c29',	'e8a91e96-6d86-49fd-a17d-4801a1c78642',	1),
('03f065ca-9eb2-494e-869c-425ff151d34c',	'cb32e347-2af0-4339-b70f-1550d943c382',	1),
('8c24442e-6aeb-438f-9320-0a90beb3e975',	'cb32e347-2af0-4339-b70f-1550d943c382',	1),
('8c24442e-6aeb-438f-9320-0a90beb3e975',	'5771c2d7-1b00-4b44-9978-0493e1d778d7',	1),
('8c24442e-6aeb-438f-9320-0a90beb3e975',	'd56c9e29-26a9-4420-a14c-0538d520f73e',	1),
('413f8397-c6a6-4221-84a7-15626b50ebe2',	'5771c2d7-1b00-4b44-9978-0493e1d778d7',	1),
('c7f5f5ca-7153-478d-b1fb-292b1d24c096',	'cb32e347-2af0-4339-b70f-1550d943c382',	1),
('97fdbd98-a85e-44d3-91c6-451357a89622',	'cb32e347-2af0-4339-b70f-1550d943c382',	1),
('6a13b922-3c4f-45d7-9961-31936a790008',	'cb32e347-2af0-4339-b70f-1550d943c382',	1),
('6a13b922-3c4f-45d7-9961-31936a790008',	'5771c2d7-1b00-4b44-9978-0493e1d778d7',	1),
('913afffa-cc3e-4980-98d1-74df9df53882',	'e8a91e96-6d86-49fd-a17d-4801a1c78642',	3),
('9d9cfa2f-6ea2-4af8-8515-e792ff48868e',	'cb32e347-2af0-4339-b70f-1550d943c382',	2),
('9d9cfa2f-6ea2-4af8-8515-e792ff48868e',	'5771c2d7-1b00-4b44-9978-0493e1d778d7',	3);

DROP TABLE IF EXISTS "medicines";
CREATE TABLE "public"."medicines" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "status" character varying(20) NOT NULL,
    "creation_date" timestamp NOT NULL,
    "formation_date" timestamp,
    "completion_date" timestamp,
    "moderator_id" uuid,
    "customer_id" uuid NOT NULL,
    "name" character varying(50),
    "verification_status" character varying(40),
    CONSTRAINT "medicines_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "medicines" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "name", "verification_status") VALUES
('913afffa-cc3e-4980-98d1-74df9df53882',	'удалён',	'2024-01-14 14:43:57.005044',	NULL,	NULL,	NULL,	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'Амоксиклав',	NULL),
('94b59b78-fa0f-48ca-ab3e-7ccd840d6c29',	'отклонён',	'2024-01-14 14:47:50.132653',	'2024-01-14 14:52:52.584729',	'2024-01-14 15:00:01.702116',	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'Амоксиклав',	'провалена'),
('03f065ca-9eb2-494e-869c-425ff151d34c',	'удалён',	'2024-01-19 17:37:35.465219',	NULL,	NULL,	NULL,	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	NULL,	NULL),
('8c24442e-6aeb-438f-9320-0a90beb3e975',	'удалён',	'2024-01-19 18:23:51.885101',	NULL,	NULL,	NULL,	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	NULL,	NULL),
('413f8397-c6a6-4221-84a7-15626b50ebe2',	'удалён',	'2024-01-19 18:25:01.021702',	NULL,	NULL,	NULL,	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'123',	NULL),
('97fdbd98-a85e-44d3-91c6-451357a89622',	'удалён',	'2024-01-19 18:52:32.229745',	NULL,	NULL,	NULL,	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'123',	NULL),
('c7f5f5ca-7153-478d-b1fb-292b1d24c096',	'отклонён',	'2024-01-19 18:29:24.776232',	'2024-01-19 18:40:10.82529',	'2024-01-19 19:23:38.000423',	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	NULL,	'провалена'),
('6a13b922-3c4f-45d7-9961-31936a790008',	'завершён',	'2024-01-19 19:14:15.436586',	'2024-01-19 19:14:59.943798',	'2024-01-19 19:23:39.274421',	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'845d1255-b852-4db0-a381-5c7cfc1d4c9e',	'132',	'пройдена'),
('9d9cfa2f-6ea2-4af8-8515-e792ff48868e',	'черновик',	'2024-01-19 19:06:40.292358',	NULL,	NULL,	NULL,	'6138b489-ddea-4b6c-b0ab-912d7ca126e7',	'123',	NULL);

DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "role" bigint,
    "login" character varying(30) NOT NULL,
    "password" character varying(40) NOT NULL,
    "name" character varying(60),
    "email" character varying(40),
    CONSTRAINT "users_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "users" ("uuid", "role", "login", "password", "name", "email") VALUES
('845d1255-b852-4db0-a381-5c7cfc1d4c9e',	1,	'user',	'12dea96fec20593566ab75692c9949596833adc9',	NULL,	NULL),
('b50e737f-05f2-453a-a6d9-9b8cf10e841f',	1,	'test',	'a94a8fe5ccb19ba61c4c0873d391e987982fbbd3',	'pepp',	'test@test.ru'),
('6138b489-ddea-4b6c-b0ab-912d7ca126e7',	2,	'admin',	'd033e22ae348aeb5660fc2140aec35850c4da997',	'Фамилия Имя Отчество',	'test@test.test');

ALTER TABLE ONLY "public"."medicine_productions" ADD CONSTRAINT "fk_medicine_productions_component" FOREIGN KEY (component_id) REFERENCES components(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."medicine_productions" ADD CONSTRAINT "fk_medicine_productions_medicine" FOREIGN KEY (medicine_id) REFERENCES medicines(uuid) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."medicines" ADD CONSTRAINT "fk_medicines_customer" FOREIGN KEY (customer_id) REFERENCES users(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."medicines" ADD CONSTRAINT "fk_medicines_moderator" FOREIGN KEY (moderator_id) REFERENCES users(uuid) NOT DEFERRABLE;

-- 2024-01-29 17:22:24.398052+00
