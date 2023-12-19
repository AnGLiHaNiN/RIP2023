-- Adminer 4.8.1 PostgreSQL 16.0 (Debian 16.0-1.pgdg120+1) dump

DROP TABLE IF EXISTS "components";
CREATE TABLE "public"."components" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "status" character varying(20) NOT NULL,
    "creation_date" timestamp NOT NULL,
    "formation_date" timestamp,
    "completion_date" timestamp,
    "moderator_id" uuid,
    "customer_id" uuid NOT NULL,
    "component_name" character varying(50) NOT NULL,
    CONSTRAINT "components_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "components" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "component_name") VALUES
('520e181a-81d8-4e9c-9289-20b75747fe4a',	'черновик',	'2023-12-10 14:53:35.910002',	NULL,	NULL,	NULL,	'a0807342-c445-48b0-aeee-a66f06ecf01f',	'');

DROP TABLE IF EXISTS "medicine_productions";
CREATE TABLE "public"."medicine_productions" (
    "medicine_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "component_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    CONSTRAINT "medicine_productions_pkey" PRIMARY KEY ("medicine_id", "component_id")
) WITH (oids = false);

INSERT INTO "medicine_productions" ("medicine_id", "component_id") VALUES
('0ed4ab2e-ae3c-447f-a18f-747f7784daa4',	'520e181a-81d8-4e9c-9289-20b75747fe4a');

DROP TABLE IF EXISTS "medicines";
CREATE TABLE "public"."medicines" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "name" character varying(100) NOT NULL,
    "image_url" character varying(100),
    "dosage" character varying(75) NOT NULL,
    "amount" bigint NOT NULL,
    "manufacturer" character varying(100) NOT NULL,
    "is_deleted" boolean DEFAULT false NOT NULL,
    CONSTRAINT "medicines_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "medicines" ("uuid", "name", "image_url", "dosage", "amount", "manufacturer", "is_deleted") VALUES
('0ed4ab2e-ae3c-447f-a18f-747f7784daa4',	'Нурофен',	'localhost:9000/images/0ed4ab2e-ae3c-447f-a18f-747f7784daa4.jpg',	'210 мг',	100,	'OOO "МедСан"',	'f'),
('ae2bb0ef-f1d9-476e-885f-d961dc64d436',	'Лоперамид',	'localhost:9000/images/Парацетомол.jpg',	'400 мг',	200,	'ООО "Калинка"',	'f'),
('96c5a987-5a8e-4ddb-9dcc-581cefcabd5c',	'Фистал',	'localhost:9000/images/Амоксициллин.jpg',	'325 мг',	245,	'ООО "Здоровье"',	'f'),
('7454a9f0-b23d-4469-b9c0-1693a8c7169f',	'Терафлю',	'localhost:9000/images/Ибупрофен.jpg',	'200 мг',	130,	'ООО "Лечись"',	'f'),
('d3b94ccb-9266-4a06-8bb3-c81a1db266c4',	'Арбидол',	'localhost:9000/images/6350910648.jpg',	'500 мг',	100,	'ООО "Фармоцевт"',	'f');

DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "role" bigint,
    "login" character varying(30) NOT NULL,
    "password" character varying(40) NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "users" ("uuid", "role", "login", "password") VALUES
('4c39276c-ace2-40bb-b3ac-417dc8bf5573',	1,	'user1',	'b3daa77b4c04a9551b8781d03191fe098f325e67'),
('a0807342-c445-48b0-aeee-a66f06ecf01f',	2,	'admin',	'd033e22ae348aeb5660fc2140aec35850c4da997');

ALTER TABLE ONLY "public"."components" ADD CONSTRAINT "fk_components_customer" FOREIGN KEY (customer_id) REFERENCES users(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."components" ADD CONSTRAINT "fk_components_moderator" FOREIGN KEY (moderator_id) REFERENCES users(uuid) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."medicine_productions" ADD CONSTRAINT "fk_medicine_productions_component" FOREIGN KEY (component_id) REFERENCES components(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."medicine_productions" ADD CONSTRAINT "fk_medicine_productions_medicine" FOREIGN KEY (medicine_id) REFERENCES medicines(uuid) NOT DEFERRABLE;

-- 2023-12-19 22:04:48.492465+00
