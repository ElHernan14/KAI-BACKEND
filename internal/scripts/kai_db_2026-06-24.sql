--
-- PostgreSQL database dump
--

\restrict oJfded7zOp2pVAWm1W9Uv6x5WEk5GFZ781NlmWPEISPM6LMbJpHEYOZRWuaGQ2O

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3

-- Started on 2026-06-24 20:51:14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 2 (class 3079 OID 59595)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 5247 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- TOC entry 922 (class 1247 OID 59966)
-- Name: acto_cuidado_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.acto_cuidado_enum AS ENUM (
    'tomar_agua',
    'respirar',
    'caminar',
    'escribir',
    'descansar_mejor',
    'comer_mejor',
    'pedir_ayuda'
);


ALTER TYPE public.acto_cuidado_enum OWNER TO postgres;

--
-- TOC entry 919 (class 1247 OID 59956)
-- Name: autoexigencia_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.autoexigencia_enum AS ENUM (
    'muy_duro_conmigo',
    'a_veces_me_castigo',
    'trato_de_entenderme',
    'estoy_aprendiendo_a_cuidarme'
);


ALTER TYPE public.autoexigencia_enum OWNER TO postgres;

--
-- TOC entry 910 (class 1247 OID 59918)
-- Name: dificultad_principal_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.dificultad_principal_enum AS ENUM (
    'cuidarme',
    'sostener_habitos',
    'descansar',
    'motivarme',
    'organizarme',
    'sentirme_bien_conmigo'
);


ALTER TYPE public.dificultad_principal_enum OWNER TO postgres;

--
-- TOC entry 907 (class 1247 OID 59904)
-- Name: estado_actual_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.estado_actual_enum AS ENUM (
    'agotado',
    'desconectado',
    'ansioso',
    'estancado',
    'con_ganas_de_empezar',
    'no_lo_se'
);


ALTER TYPE public.estado_actual_enum OWNER TO postgres;

--
-- TOC entry 913 (class 1247 OID 59932)
-- Name: necesidad_principal_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.necesidad_principal_enum AS ENUM (
    'energia',
    'calma',
    'constancia',
    'confianza',
    'conexion_conmigo',
    'bienestar_fisico'
);


ALTER TYPE public.necesidad_principal_enum OWNER TO postgres;

--
-- TOC entry 916 (class 1247 OID 59946)
-- Name: ritmo_vida_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.ritmo_vida_enum AS ENUM (
    'agotador',
    'variable',
    'tranquilo',
    'dificil_de_organizar'
);


ALTER TYPE public.ritmo_vida_enum OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 225 (class 1259 OID 59846)
-- Name: animales_catalogo; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.animales_catalogo (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    nombre character varying(50) NOT NULL,
    simbolo_principal character varying(50),
    tipo_animal character varying(30),
    rareza character varying(20),
    descripcion text,
    imagen_base character varying(100),
    atributo_principal character varying(30),
    activo boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.animales_catalogo OWNER TO postgres;

--
-- TOC entry 241 (class 1259 OID 60278)
-- Name: atributos_kai; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.atributos_kai (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    atributo_kai_id uuid NOT NULL,
    valor integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.atributos_kai OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 60217)
-- Name: categorias_xp; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categorias_xp (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    nombre character varying(30) NOT NULL,
    descripcion text
);


ALTER TABLE public.categorias_xp OWNER TO postgres;

--
-- TOC entry 234 (class 1259 OID 60124)
-- Name: configuracion_usuario; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.configuracion_usuario (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    notificaciones_activas boolean DEFAULT true,
    sonidos_activos boolean DEFAULT true,
    mostrar_rachas boolean DEFAULT true,
    modo_discreto boolean DEFAULT false,
    intensidad_kai character varying(20) DEFAULT 'normal'::character varying,
    horario_recordatorio time without time zone,
    bloquear_con_pin boolean DEFAULT false,
    permitir_mensajes_emocionales boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.configuracion_usuario OWNER TO postgres;

--
-- TOC entry 240 (class 1259 OID 60262)
-- Name: desbloqueos_kai; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.desbloqueos_kai (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    tipo_desbloqueo character varying(50),
    referencia_id uuid,
    descripcion text,
    desbloqueado_en timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    evento_usuario_id uuid
);


ALTER TABLE public.desbloqueos_kai OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 59720)
-- Name: diario_estoico; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.diario_estoico (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    fecha date NOT NULL,
    que_hice_bien text,
    que_mejorar text,
    como_me_senti text,
    agradecimiento_1 text,
    agradecimiento_2 text,
    agradecimiento_3 text,
    estado_animo character varying(20),
    xp_otorgada integer DEFAULT 15,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    evento_usuario_id uuid
);


ALTER TABLE public.diario_estoico OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 59860)
-- Name: encuentros_animales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.encuentros_animales (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    animal_id uuid NOT NULL,
    titulo_encuentro character varying(100),
    reflexion text,
    etapa_emocional character varying(50),
    nivel_encuentro integer DEFAULT 1,
    desbloqueado_por character varying(100),
    imagen_variante character varying(100),
    visto boolean DEFAULT false,
    fecha_encuentro timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.encuentros_animales OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 59753)
-- Name: estado_kai; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.estado_kai (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    estado_actual character varying(30) NOT NULL,
    etapa_actual character varying(30) NOT NULL,
    energia integer DEFAULT 100,
    imagen_kai character varying(100),
    ultimo_mensaje text,
    ultima_interaccion timestamp without time zone,
    dias_sin_actividad integer DEFAULT 0,
    modo_recuperacion boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    modo_actual character varying(30),
    atributo_dominante_id uuid,
    nivel_vinculo integer DEFAULT 1,
    ultima_evolucion timestamp without time zone
);


ALTER TABLE public.estado_kai OWNER TO postgres;

--
-- TOC entry 235 (class 1259 OID 60148)
-- Name: eventos_usuario; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.eventos_usuario (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    tipo_evento character varying(50) NOT NULL,
    referencia_id uuid,
    referencia_tabla character varying(50),
    descripcion text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    xp_otorgada integer DEFAULT 0
);


ALTER TABLE public.eventos_usuario OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 60004)
-- Name: habitos_catalogo; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.habitos_catalogo (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    nombre character varying(100) NOT NULL,
    descripcion text,
    categoria character varying(50),
    tipo_cuidado character varying(50),
    dificultad character varying(20),
    xp_base integer DEFAULT 10,
    es_premium boolean DEFAULT false,
    activo boolean DEFAULT true,
    imagen_habito character varying(100),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    categoria_xp_id uuid
);


ALTER TABLE public.habitos_catalogo OWNER TO postgres;

--
-- TOC entry 229 (class 1259 OID 60018)
-- Name: habitos_usuario; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.habitos_usuario (
    id uuid DEFAULT public.uuid_generate_v4() CONSTRAINT habitos_usuario_id_not_null1 NOT NULL,
    usuario_id uuid NOT NULL,
    habito_catalogo_id uuid NOT NULL,
    personalizado boolean DEFAULT false,
    activo boolean DEFAULT true,
    fecha_inicio date DEFAULT CURRENT_DATE,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.habitos_usuario OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 60242)
-- Name: mensaje_atributos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mensaje_atributos (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    mensaje_kai_id uuid NOT NULL,
    atributo_kai_id uuid NOT NULL,
    nivel_minimo integer DEFAULT 1
);


ALTER TABLE public.mensaje_atributos OWNER TO postgres;

--
-- TOC entry 233 (class 1259 OID 60110)
-- Name: mensajes_kai; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mensajes_kai (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    tipo character varying(50) NOT NULL,
    contexto character varying(50),
    tono character varying(30),
    mensaje text NOT NULL,
    rareza character varying(20) DEFAULT 'comun'::character varying,
    desbloqueado_por character varying(100),
    activo boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.mensajes_kai OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 60165)
-- Name: mensajes_usuario; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mensajes_usuario (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    mensaje_kai_id uuid NOT NULL,
    leido boolean DEFAULT false,
    mostrado_en timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.mensajes_usuario OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 59699)
-- Name: misiones_diarias; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.misiones_diarias (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    titulo character varying(150) NOT NULL,
    descripcion text,
    area character varying(30),
    xp_recompensa integer DEFAULT 10,
    es_mision_minima boolean DEFAULT false,
    completada boolean DEFAULT false,
    fecha date NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    habito_usuario_id uuid
);


ALTER TABLE public.misiones_diarias OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 60091)
-- Name: momentos_memorables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.momentos_memorables (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    titulo character varying(150) NOT NULL,
    descripcion text,
    tipo_momento character varying(50),
    imagen_referencia character varying(100),
    reflexion text,
    relacionado_con character varying(50),
    importante boolean DEFAULT true,
    fecha_momento timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    encuentro_animal_id uuid,
    mensaje_usuario_id uuid,
    evento_usuario_id uuid
);


ALTER TABLE public.momentos_memorables OWNER TO postgres;

--
-- TOC entry 227 (class 1259 OID 59981)
-- Name: onboarding_respuestas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.onboarding_respuestas (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    estado_actual public.estado_actual_enum NOT NULL,
    dificultad_principal public.dificultad_principal_enum NOT NULL,
    necesidad_principal public.necesidad_principal_enum NOT NULL,
    ritmo_vida public.ritmo_vida_enum NOT NULL,
    autoexigencia public.autoexigencia_enum NOT NULL,
    acto_cuidado_inicial public.acto_cuidado_enum NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.onboarding_respuestas OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 59738)
-- Name: progreso_fisico; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.progreso_fisico (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    fecha date NOT NULL,
    peso numeric(5,2),
    cintura numeric(5,1),
    cadera numeric(5,1),
    pecho numeric(5,1),
    brazo numeric(5,1),
    muslo numeric(5,1),
    pantorrilla numeric(5,1),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    evento_usuario_id uuid
);


ALTER TABLE public.progreso_fisico OWNER TO postgres;

--
-- TOC entry 231 (class 1259 OID 60066)
-- Name: rachas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rachas (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    habito_usuario_id uuid NOT NULL,
    dias_actuales integer DEFAULT 0,
    record_historico integer DEFAULT 0,
    inicio_racha date,
    ultima_actividad date,
    activa boolean DEFAULT true,
    protegida boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.rachas OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 60041)
-- Name: registros_habito; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.registros_habito (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    habito_usuario_id uuid NOT NULL,
    fecha date NOT NULL,
    completado boolean DEFAULT false,
    valor_registrado text,
    xp_ganada integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.registros_habito OWNER TO postgres;

--
-- TOC entry 238 (class 1259 OID 60229)
-- Name: tipos_atributo_kai; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tipos_atributo_kai (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    nombre character varying(30) NOT NULL,
    descripcion text
);


ALTER TABLE public.tipos_atributo_kai OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 59606)
-- Name: usuarios; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.usuarios (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    nombre character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    password_hash text NOT NULL,
    fecha_registro timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    perfil_base character varying(30),
    etapa_kai character varying(30) DEFAULT 'cachorro'::character varying,
    racha_global integer DEFAULT 0,
    dias_inactivo integer DEFAULT 0
);


ALTER TABLE public.usuarios OWNER TO postgres;

--
-- TOC entry 242 (class 1259 OID 60300)
-- Name: xp_atributos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.xp_atributos (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    categoria_xp_id uuid NOT NULL,
    atributo_kai_id uuid NOT NULL,
    multiplicador numeric(4,2) DEFAULT 1.0
);


ALTER TABLE public.xp_atributos OWNER TO postgres;

--
-- TOC entry 243 (class 1259 OID 60331)
-- Name: xp_usuario; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.xp_usuario (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    usuario_id uuid NOT NULL,
    categoria_xp_id uuid NOT NULL,
    valor integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.xp_usuario OWNER TO postgres;

--
-- TOC entry 5223 (class 0 OID 59846)
-- Dependencies: 225
-- Data for Name: animales_catalogo; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.animales_catalogo (id, nombre, simbolo_principal, tipo_animal, rareza, descripcion, imagen_base, atributo_principal, activo, created_at) FROM stdin;
45bb0b0f-67a3-4456-8275-1cacdbc8b609	Búho	sabiduria	reflexivo	raro	Representa la capacidad de observarse y aprender del propio camino.	buho_base	sabiduria	t	2026-05-06 21:08:55.641677
247f2b44-fa44-42b4-9341-1d28e67bed0a	Lobo	resiliencia	protector	epico	Representa la fuerza de seguir adelante incluso después de perderse.	lobo_base	resistencia	t	2026-05-06 21:08:55.641677
364c7ebd-fafd-4cd7-9ce9-c524e4223121	Zorro	adaptabilidad	astuto	raro	Representa la inteligencia para adaptarse y encontrar nuevos caminos.	zorro_base	conciencia	t	2026-05-06 21:08:55.641677
3e99b7a3-0781-4c8e-b75f-6974d09e626a	Ciervo	calma	sereno	comun	Representa la sensibilidad, la tranquilidad y el avance sin violencia.	ciervo_base	equilibrio	t	2026-05-06 21:08:55.641677
ee9d8a80-0f1d-46e6-8d5e-b3fa27f1bce8	Oso	fortaleza	guardian	epico	Representa la fuerza interior y la protección emocional.	oso_base	fuerza	t	2026-05-06 21:08:55.641677
23b15958-ded3-406c-8f2d-084cbd5375cd	Colibrí	pequenos_pasos	ligero	comun	Representa los pequeños avances diarios que terminan transformando una vida.	colibri_base	vitalidad	t	2026-05-06 21:08:55.641677
e68eaf8e-2837-49cb-9cd2-241dc065834a	Tortuga	constancia	persistente	comun	Representa el progreso lento pero constante.	tortuga_base	disciplina	t	2026-05-06 21:08:55.641677
f41ae20c-e70f-4322-b4d8-718eb56d1326	Cuervo	transformacion	mistico	legendario	Representa los cambios profundos que nacen después de atravesar oscuridad.	cuervo_base	conciencia	t	2026-05-06 21:08:55.641677
a2fa0867-a243-4052-9dee-91a007ca3d49	León	coraje	lider	legendario	Representa el valor de enfrentar el miedo y avanzar igualmente.	leon_base	fuerza	t	2026-05-06 21:08:55.641677
b4cf2d23-53aa-4904-b034-e6bfe190faf0	Mariposa	renacimiento	espiritual	raro	Representa la transformación personal y los nuevos comienzos.	mariposa_base	equilibrio	t	2026-05-06 21:08:55.641677
c35b7a29-e728-4a54-b109-3688a07d73f5	Nutria	alegria	social	comun	Representa la importancia del juego, el descanso y la conexión emocional.	nutria_base	vitalidad	t	2026-05-06 21:08:55.641677
f06e2460-b575-43ae-9004-e26192cf32c5	Águila	vision	visionario	legendario	Representa la capacidad de ver más allá del momento presente.	aguila_base	sabiduria	t	2026-05-06 21:08:55.641677
\.


--
-- TOC entry 5239 (class 0 OID 60278)
-- Dependencies: 241
-- Data for Name: atributos_kai; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.atributos_kai (id, usuario_id, atributo_kai_id, valor, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 5235 (class 0 OID 60217)
-- Dependencies: 237
-- Data for Name: categorias_xp; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categorias_xp (id, nombre, descripcion) FROM stdin;
08f90c37-7224-4da5-b796-d83b323be9ae	Vitalidad	Energía física y autocuidado corporal
320a374a-6c0a-49d3-ba51-e646ca27d906	Movimiento	Actividad física y movimiento consciente
41d43ad1-0f76-41a4-8298-b41f374f1f11	Disciplina	Constancia y cumplimiento de compromisos personales
c02beabc-14c5-47c1-8278-35786865367f	Sabiduria	Reflexión, aprendizaje y autoconocimiento
2b87c63a-cfa7-4a30-a73d-d98b80eb0a6c	Equilibrio	Bienestar emocional y regulación interna
17dcb0db-37dc-42bb-8865-394e669e9572	Coraje	Acciones difíciles realizadas a pesar del miedo
5dd6a8b4-18b3-4cea-aef8-c0d8cf7fe9e5	Conexion	Relaciones positivas y apoyo social
c770a520-57aa-4143-9cb1-098d425584bd	Constancia	Capacidad de volver y sostener el proceso
\.


--
-- TOC entry 5232 (class 0 OID 60124)
-- Dependencies: 234
-- Data for Name: configuracion_usuario; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.configuracion_usuario (id, usuario_id, notificaciones_activas, sonidos_activos, mostrar_rachas, modo_discreto, intensidad_kai, horario_recordatorio, bloquear_con_pin, permitir_mensajes_emocionales, created_at, updated_at) FROM stdin;
d33403fc-d72c-40fe-8736-8a6eed1f5a10	87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7	t	t	t	f	suave	\N	f	t	2026-05-07 16:47:26.504542	2026-05-07 16:47:26.504542
\.


--
-- TOC entry 5238 (class 0 OID 60262)
-- Dependencies: 240
-- Data for Name: desbloqueos_kai; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.desbloqueos_kai (id, usuario_id, tipo_desbloqueo, referencia_id, descripcion, desbloqueado_en, evento_usuario_id) FROM stdin;
\.


--
-- TOC entry 5220 (class 0 OID 59720)
-- Dependencies: 222
-- Data for Name: diario_estoico; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.diario_estoico (id, usuario_id, fecha, que_hice_bien, que_mejorar, como_me_senti, agradecimiento_1, agradecimiento_2, agradecimiento_3, estado_animo, xp_otorgada, created_at, evento_usuario_id) FROM stdin;
dca8924a-8329-473d-9081-605e8e78545d	87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7	2026-05-06	Caminé hoy	Dormir más temprano	Me sentí tranquilo	Mi familia	Haber descansado	La comida de hoy	bien	15	2026-05-06 20:21:40.141051	\N
\.


--
-- TOC entry 5224 (class 0 OID 59860)
-- Dependencies: 226
-- Data for Name: encuentros_animales; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.encuentros_animales (id, usuario_id, animal_id, titulo_encuentro, reflexion, etapa_emocional, nivel_encuentro, desbloqueado_por, imagen_variante, visto, fecha_encuentro, created_at) FROM stdin;
\.


--
-- TOC entry 5222 (class 0 OID 59753)
-- Dependencies: 224
-- Data for Name: estado_kai; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.estado_kai (id, usuario_id, estado_actual, etapa_actual, energia, imagen_kai, ultimo_mensaje, ultima_interaccion, dias_sin_actividad, modo_recuperacion, created_at, updated_at, modo_actual, atributo_dominante_id, nivel_vinculo, ultima_evolucion) FROM stdin;
16a11ef0-38c3-4ed9-b752-71f146ccddf1	87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7	activo	cachorro	100	cachorro_activo	Empezamos. Eso ya es algo grande para mí.	2026-05-06 20:38:18.556979	0	f	2026-05-06 20:38:18.556979	2026-05-06 20:38:18.556979	\N	\N	1	\N
\.


--
-- TOC entry 5233 (class 0 OID 60148)
-- Dependencies: 235
-- Data for Name: eventos_usuario; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.eventos_usuario (id, usuario_id, tipo_evento, referencia_id, referencia_tabla, descripcion, created_at, xp_otorgada) FROM stdin;
\.


--
-- TOC entry 5226 (class 0 OID 60004)
-- Dependencies: 228
-- Data for Name: habitos_catalogo; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.habitos_catalogo (id, nombre, descripcion, categoria, tipo_cuidado, dificultad, xp_base, es_premium, activo, imagen_habito, created_at, categoria_xp_id) FROM stdin;
a48b647b-fe23-4163-bac2-04092ec9eb0e	Registrar mi progreso	Reconocer que incluso los pequeños pasos cuentan.	Disciplina	Mental	facil	15	f	t	facil.webp	2026-05-06 21:46:28.878065	41d43ad1-0f76-41a4-8298-b41f374f1f11
f08197b8-3a81-4b85-af8a-75bd383784cf	Tender la cama	Un pequeño acto de orden puede cambiar el inicio del día.	Disciplina	Mental	facil	10	f	t	facil.webp	2026-05-06 21:46:28.878065	41d43ad1-0f76-41a4-8298-b41f374f1f11
82a4c3f3-1bb4-48d1-bef7-870e598632a5	Leer una reflexión	Tomarte un momento para pensar distinto.	Sabiduria	Mental	facil	15	f	t	facil.webp	2026-05-06 21:46:28.878065	c02beabc-14c5-47c1-8278-35786865367f
8d5fa949-8f54-4346-9de5-4cfc2891cc2d	Desconectarme unos minutos	Alejarse un momento del ruido también ayuda.	Vitalidad	Emocional	facil	10	f	t	facil.webp	2026-05-06 21:46:28.878065	08f90c37-7224-4da5-b796-d83b323be9ae
68582ea2-9726-4282-a0d4-878839a96d69	Entrenar	Realizar entrenamiento planificado	Movimiento	Fisico	media	25	f	t	medio.webp	2026-05-21 19:28:47.683141	320a374a-6c0a-49d3-ba51-e646ca27d906
6c6bdaf4-60d0-4e60-ba85-fe234c7f9f1c	Leer 10 minutos	Lectura para crecimiento personal	Sabiduria	Mental	facil	15	f	t	facil.webp	2026-05-21 19:30:17.766069	c02beabc-14c5-47c1-8278-35786865367f
c2d0d659-09c3-4949-8e43-79bf0c590908	Diario estoico	Reflexion personal diaria	Sabiduria	Mental	facil	15	f	t	facil.webp	2026-05-21 19:29:55.613741	c02beabc-14c5-47c1-8278-35786865367f
1639535d-ba7e-4b8e-9b86-196917b98caa	Registrar emociones	Identificar y registrar estado emocional	Equilibrio	Emocional	facil	15	f	t	facil.webp	2026-05-21 19:31:20.798563	2b87c63a-cfa7-4a30-a73d-d98b80eb0a6c
3f5e41a2-25d8-4a22-a8c2-a046c872991b	Meditacion breve	Momento corto de meditacion	Equilibrio	Emocional	facil	15	f	t	facil.webp	2026-05-21 19:31:02.683292	2b87c63a-cfa7-4a30-a73d-d98b80eb0a6c
8f167195-54cb-48c7-849a-26884ea4a462	Intentar algo nuevo	Salir de la zona de confort	Coraje	Mental	media	20	f	t	medio.webp	2026-05-21 19:32:12.405826	17dcb0db-37dc-42bb-8865-394e669e9572
139c6544-6556-47e0-acaf-8ead4f0d5ca9	Resolver tarea pendiente	Completar algo postergado	Coraje	Mental	media	20	f	t	medio.webp	2026-05-21 19:31:53.074861	17dcb0db-37dc-42bb-8865-394e669e9572
89f02213-80a0-43c5-a490-fffc0aff8a7d	Dormir un poco antes	Descansar también es parte del crecimiento.	Vitalidad	Fisico	media	20	f	t	medio.webp	2026-05-06 21:46:28.878065	08f90c37-7224-4da5-b796-d83b323be9ae
63de49fd-97a1-4a8e-8b4a-5fd10e88654c	Caminar 15 minutos	Actividad fisica ligera sostenida	Movimiento	Fisico	facil	15	f	t	facil.webp	2026-05-21 19:28:31.998092	320a374a-6c0a-49d3-ba51-e646ca27d906
271bf66d-f914-4a63-8707-33b5422f5064	Caminar 5 minutos	Movimiento mínimo diario	Movimiento	Fisico	facil	10	f	t	facil.webp	2026-05-21 18:55:24.701962	320a374a-6c0a-49d3-ba51-e646ca27d906
1025cf1f-141d-46b1-84b3-1c711b451119	Hablar con alguien querido	Fortalecer vinculos personales	Conexion	Social	facil	15	f	t	facil.webp	2026-05-21 19:33:19.205575	5dd6a8b4-18b3-4cea-aef8-c0d8cf7fe9e5
42776a23-91b6-41a5-9499-9336afda91cd	Agradecer a alguien	Expresar gratitud a otra persona	Conexion	Social	facil	15	f	t	facil.webp	2026-05-21 19:32:55.322355	5dd6a8b4-18b3-4cea-aef8-c0d8cf7fe9e5
c67ab0ad-2de2-4d49-aad2-105ed550fed1	Volver despues de una ausencia	Retomar el camino luego de varios dias	Constancia	Identidad	media	50	f	t	medio.webp	2026-05-21 19:34:40.858698	c770a520-57aa-4143-9cb1-098d425584bd
5e160ad6-d11f-4fbb-8edf-27c9d02e3bf4	Mision minima	Completar una accion minima posible	Constancia	Identidad	facil	10	f	t	facil.webp	2026-05-21 19:34:24.34734	c770a520-57aa-4143-9cb1-098d425584bd
0ceb4237-c71e-4caa-92de-d60cdb3edd9d	Hablarme con amabilidad	Intentar tratarte como tratarías a alguien que querés.	Equilibrio	Emocional	media	25	f	t	medio.webp	2026-05-06 21:46:28.878065	2b87c63a-cfa7-4a30-a73d-d98b80eb0a6c
46361232-531b-4181-93ef-7684d11cb942	Salir al sol	Conectar unos minutos con la luz del día.	Vitalidad	Emocional	facil	15	f	t	facil.webp	2026-05-06 21:46:28.878065	08f90c37-7224-4da5-b796-d83b323be9ae
1f119006-d510-4b57-af30-95b05761b459	Escuchar música tranquila	Permitirte bajar el ritmo por un momento.	Vitalidad	Emocional	facil	10	f	t	facil.webp	2026-05-06 21:46:28.878065	08f90c37-7224-4da5-b796-d83b323be9ae
f3cd56c3-889b-43d0-bfea-00ced6a81a49	Escribir cómo me siento	Poner en palabras lo que llevas dentro.	Sabiduria	Emocional	facil	20	f	t	facil.webp	2026-05-06 21:46:28.878065	c02beabc-14c5-47c1-8278-35786865367f
95cc5a64-86d2-404d-afc1-9eeecb5857cd	Tomar agua	Un pequeño acto de cuidado para tu cuerpo.	Vitalidad	Fisico	facil	10	f	t	facil.webp	2026-05-06 21:46:28.878065	08f90c37-7224-4da5-b796-d83b323be9ae
658170fa-4544-4296-8e9d-309590489837	Comer algo nutritivo	Nutrirte también es una forma de cuidarte.	Vitalidad	Fisico	facil	15	f	t	facil.webp	2026-05-06 21:46:28.878065	08f90c37-7224-4da5-b796-d83b323be9ae
865710a2-84c8-464a-899b-fa99733200aa	Estirar el cuerpo	Escuchar el cuerpo y darle un poco de atención.	Movimiento	Fisico	facil	15	f	t	facil.webp	2026-05-06 21:46:28.878065	320a374a-6c0a-49d3-ba51-e646ca27d906
0fe0b1cc-4752-4b53-8d91-35b6366222a8	Caminar unos minutos	Mover el cuerpo suavemente también es autocuidado.	Movimiento	Fisico	facil	15	f	t	facil.webp	2026-05-06 21:46:28.878065	320a374a-6c0a-49d3-ba51-e646ca27d906
abf029df-ec3b-472a-baa9-b1ef6cfcbfb1	Registrar tres agradecimientos	Reconocer pequeños momentos valiosos del día.	Sabiduria	Emocional	facil	20	f	t	facil.webp	2026-05-06 21:46:28.878065	c02beabc-14c5-47c1-8278-35786865367f
a988c92a-a043-46c4-8e7f-9b8dee36d502	Respirar conscientemente	Detenerse unos minutos también es avanzar.	Equilibrio	Emocional	facil	10	f	t	facil.webp	2026-05-06 21:46:28.878065	2b87c63a-cfa7-4a30-a73d-d98b80eb0a6c
\.


--
-- TOC entry 5227 (class 0 OID 60018)
-- Dependencies: 229
-- Data for Name: habitos_usuario; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.habitos_usuario (id, usuario_id, habito_catalogo_id, personalizado, activo, fecha_inicio, created_at) FROM stdin;
\.


--
-- TOC entry 5237 (class 0 OID 60242)
-- Dependencies: 239
-- Data for Name: mensaje_atributos; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mensaje_atributos (id, mensaje_kai_id, atributo_kai_id, nivel_minimo) FROM stdin;
5abe185b-26c8-42e3-ba87-e15ef37f2f0a	86b7f894-7ec2-42d5-a7e9-caa0ed3d02c4	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
32ed286c-72d1-41b7-9594-b76828801c16	02cc02d3-892f-4428-862e-225f4e55992c	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
a2a595b8-7198-4744-b910-71e3999a53eb	de0fdf95-1d4b-4bb4-a4ee-26438c9f0309	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
6369f8ff-263c-4087-a44a-c874bc1ad051	6e807be7-53cc-4512-b585-0c9962a95ebd	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
c877f1df-8bb6-46f6-98cd-77c0060b70ef	55bd9046-fe0b-40d0-83af-2babbddb9e22	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
5785edcb-8e05-4fc8-ada3-fb0ae1f84588	ec9538e9-ec0f-479d-9fbd-bbf298f630fe	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
178654c4-494d-4acc-ae88-6e0e90f7f2b9	508a7d08-7e8c-4367-90e9-390387c5f2f5	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
d39cb2e9-2af0-4196-942a-7c87583b6807	6c3baa1a-467b-420c-88cf-3f49ed98341e	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
8562651f-6d55-49e3-a0e9-19fafe501c0f	a7d7445d-10f9-4ae3-915b-965b3f1394e8	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
6d8ce419-2891-4d22-a380-99f330c7d543	b42223d9-36b1-4992-bd15-2438fe53ddc1	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
4ec4b1d6-cad8-45dc-aeae-44522d70f4f1	62e93688-5ae4-4993-858d-3712e842e0fc	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
f9ef013e-9973-4329-baba-ca7ae766b61c	007996f8-df7a-48a2-bc84-6211ee7deaf9	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
5f3a8e1a-d97e-47b1-8bc2-fcb7909e3d84	82a7efbb-5b9b-4f2d-9dd4-4590f8c57d3a	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	10
0063f405-6979-4621-8b48-9b61265f35c9	d2ad73a8-fb55-494e-a2ba-a3df1d185223	da217a67-56c5-4a19-aa76-f97663a2a065	10
1e41612f-947e-4bb9-821e-fea5a2e25327	25856ab0-61ab-485f-ab71-220618b9daa0	da217a67-56c5-4a19-aa76-f97663a2a065	10
7a2d6251-b4aa-4b20-9216-bd3f169ba9b5	d59e8a8b-8724-4d79-bc2d-3d64c5309be2	da217a67-56c5-4a19-aa76-f97663a2a065	10
1e1b1225-bf6d-4115-8ca5-70a242fa7754	7c99af56-5a7d-47df-ad44-86aaa5abe55a	da217a67-56c5-4a19-aa76-f97663a2a065	10
b0380f96-a21a-41d7-aa6e-e3f5681ea8de	00d7114e-1dd2-40dd-8c16-602c34fd81ac	da217a67-56c5-4a19-aa76-f97663a2a065	10
2be202e3-ea93-4a91-a20b-80ff203bc316	f35bfb3e-9866-4f7a-8be2-ff5acb81ea87	da217a67-56c5-4a19-aa76-f97663a2a065	10
b11269d5-a6ed-42ca-83fa-e3c52bc504da	92bc0995-b5fe-44be-9cbf-4d349c11a4c5	da217a67-56c5-4a19-aa76-f97663a2a065	10
2ef0d991-5cc0-48bb-b9b1-b6fd44b502e6	8f3e6093-5816-4e95-a40f-53a78a1085a7	da217a67-56c5-4a19-aa76-f97663a2a065	10
f8cb581a-7570-4d80-8cc2-0326023db73e	e7d51059-8912-42e9-89f9-e15ba00f18bf	da217a67-56c5-4a19-aa76-f97663a2a065	10
936d7830-9e96-438a-80fe-2c345907acd7	f6bdbd2f-26a1-4ef9-b5ec-84124b8e1fbd	da217a67-56c5-4a19-aa76-f97663a2a065	10
de8742f8-b420-43b0-a6bf-62cf30fa8ef8	44b4c59a-b0b2-456b-8009-dd6cc867cd51	da217a67-56c5-4a19-aa76-f97663a2a065	10
c187186f-9594-4309-be16-9c9ae17cdce7	39ee4b75-ac5b-40b8-b65c-5d1c15a6af7b	84eda4c2-31a8-463b-b272-75a7fbf6421d	10
aeb937a5-9dbd-45bc-8fcf-c3c6783cca2b	b07979be-c3b2-49be-ba41-51ec952eea34	84eda4c2-31a8-463b-b272-75a7fbf6421d	10
9da9e188-5cc4-42b0-ac2f-c417b6a32aa6	333f4b9b-7639-40b4-a227-57d7d03f6951	84eda4c2-31a8-463b-b272-75a7fbf6421d	10
91a83674-cd30-4310-b839-dee81d029def	c5d0efa4-f8b5-4144-be05-e0d99675e02f	84eda4c2-31a8-463b-b272-75a7fbf6421d	10
6c8e3b88-d80e-4ae4-8b11-8b5338aac750	824a27a6-e3f8-4bb9-9a90-c7d246d0451e	84eda4c2-31a8-463b-b272-75a7fbf6421d	10
d57bf709-6d02-4eaf-803f-64ff5fed0791	438da8e4-481b-4711-b714-94cd43f700d3	d0962734-7507-44d7-b276-2374be80c9f0	10
6c6d8873-5aea-4a5e-8267-429b18c218b3	df2cb650-c1eb-44cc-8966-2d5a0c9ed76f	d0962734-7507-44d7-b276-2374be80c9f0	10
470c2634-a4c9-4383-9230-612925722a74	ef6fb547-3886-42a5-bd41-9638e0dd3736	d0962734-7507-44d7-b276-2374be80c9f0	10
44a96eb3-efe1-4aab-bdad-23262935e6c6	3b9b792d-98f7-4f56-9351-7e3d83698883	d0962734-7507-44d7-b276-2374be80c9f0	10
5da76014-21b8-4e8a-869a-84b15b2a21bd	65380409-6645-4894-8f2f-e64de3dd26b6	d0962734-7507-44d7-b276-2374be80c9f0	10
8cf3c32e-e1fb-4b1a-b8f8-5fa152b61b07	b2c54c4f-536c-4149-a611-7c399e0d5e8e	e7319970-d7f0-4abc-a069-63e430c75a95	10
9eda94f1-857c-4547-9623-5f52ce414618	5964d19c-1a5e-410c-95ed-300a056e316c	e7319970-d7f0-4abc-a069-63e430c75a95	10
2dbd64e1-aff3-4289-8109-b98d73dfba59	83796f02-3b44-4856-8019-9e907b5b4b1e	e7319970-d7f0-4abc-a069-63e430c75a95	10
a5e1c100-d1da-4e97-b70b-cd0595e0cdd7	4bb16f39-47ed-40f8-b51a-2d1c0a9eb0fc	e7319970-d7f0-4abc-a069-63e430c75a95	10
5504f09c-31a6-4912-8b69-df39025db461	3507e965-063f-40a6-93d7-21dcbee27c8c	e7319970-d7f0-4abc-a069-63e430c75a95	10
548b3657-e004-4a1e-9205-4f7c04324b03	508a7d08-7e8c-4367-90e9-390387c5f2f5	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	10
56058d11-8221-48a3-bb7e-03245606e33f	6c3baa1a-467b-420c-88cf-3f49ed98341e	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	10
baed8f94-b0b5-4fc4-98c1-4aa259aaf15c	a7d7445d-10f9-4ae3-915b-965b3f1394e8	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	10
4c203ef2-9094-4a3d-ad25-9742aa489038	1c26a69b-8a93-4f27-a17b-9f5961771835	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	10
cde86ab9-c928-41b9-afd2-9f7e20e80b4d	b135ec79-7cbc-420c-a5fb-2ccf24447095	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	10
e12af810-b887-4448-8864-f2f2e924a993	a2c3ca41-4d7b-4fe1-b7a0-c7270a0705d1	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	10
e4a90bcb-1b76-42c8-aed2-9536c052a2a3	fdfc10ed-65e3-4af3-8824-b1d63fb59165	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	10
e5bd2950-379a-4cfd-86e9-9b9f952971ea	27967fd8-d26c-4a5f-a69a-942e315f5bf4	0d90fffa-fb4e-4905-91ce-51869671bbcb	10
df797e34-a901-43d5-a433-40564610eb61	17f55f12-0565-4f2f-841f-1bc0c65b378a	0d90fffa-fb4e-4905-91ce-51869671bbcb	10
7d253a2f-f18d-4cf1-a57a-e22bc7cbb909	bb2e52b3-b151-4d63-8942-c84e55419ff0	0d90fffa-fb4e-4905-91ce-51869671bbcb	10
e3567f28-2be4-4b3b-9340-ed3ae65a5d24	e16c2426-3979-4d4e-83f0-80243b0c0cd4	0d90fffa-fb4e-4905-91ce-51869671bbcb	10
35e882fa-9a81-49c1-8bd3-addbbf9f46de	e8a0a771-845c-4d18-a3cd-5b32e1694445	42a645e9-7064-41b6-93ca-2a14266aee6d	10
9748ae11-60b1-4aa1-8dc1-36e36b94b697	48b5e797-8f6f-443e-8a09-e39627221ee7	42a645e9-7064-41b6-93ca-2a14266aee6d	10
826fac8d-3b15-47f8-8aa4-ae23a6741fdd	46249506-baa6-4f0b-aa79-8d23c340c6d5	42a645e9-7064-41b6-93ca-2a14266aee6d	10
6c73e0a9-9377-4a32-bc77-827ced4206bc	5e8f12ec-dc50-4285-ae70-ae6e455397c5	42a645e9-7064-41b6-93ca-2a14266aee6d	10
\.


--
-- TOC entry 5231 (class 0 OID 60110)
-- Dependencies: 233
-- Data for Name: mensajes_kai; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mensajes_kai (id, tipo, contexto, tono, mensaje, rareza, desbloqueado_por, activo, created_at) FROM stdin;
b2c54c4f-536c-4149-a611-7c399e0d5e8e	regreso	abandono	suave	Volver también es una forma de avanzar.	comun	\N	t	2026-05-07 16:45:57.550318
438da8e4-481b-4711-b714-94cd43f700d3	cansancio	agotamiento	compania	No hace falta hacerlo perfecto hoy.	comun	\N	t	2026-05-07 16:45:57.550318
86b7f894-7ec2-42d5-a7e9-caa0ed3d02c4	reflexion	diario	reflexivo	A veces entenderse ya es un gran paso.	comun	\N	t	2026-05-07 16:45:57.550318
d2ad73a8-fb55-494e-a2ba-a3df1d185223	racha	continuidad	celebracion	Pequeños pasos repetidos crean caminos.	comun	\N	t	2026-05-07 16:45:57.550318
39ee4b75-ac5b-40b8-b65c-5d1c15a6af7b	autocuidado	descanso	suave	Descansar también es parte del crecimiento.	comun	\N	t	2026-05-07 16:45:57.550318
25856ab0-61ab-485f-ab71-220618b9daa0	motivacion	general	calido	Una accion pequena sigue siendo una accion importante.	comun	\N	t	2026-05-21 19:00:27.766486
d59e8a8b-8724-4d79-bc2d-3d64c5309be2	motivacion	general	calido	No necesitas hacerlo perfecto. Solo necesitas aparecer.	comun	\N	t	2026-05-21 19:00:27.766486
02cc02d3-892f-4428-862e-225f4e55992c	estoico	accion	reflexivo	Concentrate en lo que puedes controlar hoy.	raro	\N	t	2026-05-21 19:00:27.766486
de0fdf95-1d4b-4bb4-a4ee-26438c9f0309	estoico	dificultad	sereno	Los obstaculos tambien forman parte del camino.	raro	\N	t	2026-05-21 19:05:13.395496
6e807be7-53cc-4512-b585-0c9962a95ebd	estoico	accion	reflexivo	Lo importante es la accion que puedes realizar hoy.	comun	\N	t	2026-05-21 19:06:26.262342
55bd9046-fe0b-40d0-83af-2babbddb9e22	estoico	control	reflexivo	No todo depende de ti. Tus acciones si.	comun	\N	t	2026-05-21 19:06:26.262342
ec9538e9-ec0f-479d-9fbd-bbf298f630fe	estoico	presente	sereno	Un paso presente vale mas que cien planes futuros.	comun	\N	t	2026-05-21 19:06:26.262342
508a7d08-7e8c-4367-90e9-390387c5f2f5	conciencia	reflexion	reflexivo	La conciencia crece cuando te permites observar sin juzgar.	raro	\N	t	2026-05-21 21:49:51.875923
6c3baa1a-467b-420c-88cf-3f49ed98341e	conciencia	presencia	reflexivo	Observar lo que sientes es el primer paso para comprenderlo.	comun	\N	t	2026-05-21 21:50:17.3844
a7d7445d-10f9-4ae3-915b-965b3f1394e8	conciencia	autoconocimiento	sereno	Conocerte mejor también es una forma de cuidarte.	comun	\N	t	2026-05-21 21:50:17.3844
b42223d9-36b1-4992-bd15-2438fe53ddc1	reflexion	aprendizaje	reflexivo	Cada experiencia puede enseñarte algo si te permites observarla.	comun	\N	t	2026-05-21 21:53:48.156843
62e93688-5ae4-4993-858d-3712e842e0fc	reflexion	introspeccion	sereno	A veces la pausa aporta más claridad que la velocidad.	comun	\N	t	2026-05-21 21:53:48.156843
007996f8-df7a-48a2-bc84-6211ee7deaf9	reflexion	emociones	reflexivo	Comprender una emoción es mejor que ignorarla.	comun	\N	t	2026-05-21 21:53:48.156843
82a7efbb-5b9b-4f2d-9dd4-4590f8c57d3a	reflexion	autoconocimiento	reflexivo	Mirar hacia dentro requiere valentía.	raro	\N	t	2026-05-21 21:53:48.156843
7c99af56-5a7d-47df-ad44-86aaa5abe55a	motivacion	progreso	calido	Pequeños avances también son avances.	comun	\N	t	2026-05-21 21:53:57.95464
00d7114e-1dd2-40dd-8c16-602c34fd81ac	motivacion	continuidad	calido	No necesitas hacerlo perfecto para seguir creciendo.	comun	\N	t	2026-05-21 21:53:57.95464
f35bfb3e-9866-4f7a-8be2-ff5acb81ea87	motivacion	accion	calido	Hoy ya hiciste algo por ti.	comun	\N	t	2026-05-21 21:53:57.95464
92bc0995-b5fe-44be-9cbf-4d349c11a4c5	motivacion	persistencia	calido	Lo importante es volver a intentarlo.	raro	\N	t	2026-05-21 21:53:57.95464
5964d19c-1a5e-410c-95ed-300a056e316c	regreso	reencuentro	suave	Me alegra verte nuevamente.	comun	\N	t	2026-05-21 21:54:08.422666
83796f02-3b44-4856-8019-9e907b5b4b1e	regreso	continuidad	suave	El camino continúa desde donde estás.	comun	\N	t	2026-05-21 21:54:08.422666
4bb16f39-47ed-40f8-b51a-2d1c0a9eb0fc	regreso	habitos	companero	Nunca es tarde para retomar un hábito.	comun	\N	t	2026-05-21 21:54:08.422666
3507e965-063f-40a6-93d7-21dcbee27c8c	regreso	crecimiento	companero	Volver también es progreso.	raro	\N	t	2026-05-21 21:54:08.422666
b07979be-c3b2-49be-ba41-51ec952eea34	autocuidado	bienestar	suave	Cuidarte también es una prioridad.	comun	\N	t	2026-05-21 21:54:19.899707
333f4b9b-7639-40b4-a227-57d7d03f6951	autocuidado	salud	calido	Tu bienestar merece atención.	comun	\N	t	2026-05-21 21:54:19.899707
c5d0efa4-f8b5-4144-be05-e0d99675e02f	autocuidado	habitos	suave	Los pequeños cuidados tienen impacto.	comun	\N	t	2026-05-21 21:54:19.899707
824a27a6-e3f8-4bb9-9a90-c7d246d0451e	autocuidado	escucha	calido	Escucharte también es importante.	raro	\N	t	2026-05-21 21:54:19.899707
df2cb650-c1eb-44cc-8966-2d5a0c9ed76f	cansancio	descanso	sereno	Descansar no significa rendirse.	comun	\N	t	2026-05-21 21:54:32.736802
ef6fb547-3886-42a5-bd41-9638e0dd3736	cansancio	recuperacion	sereno	La recuperación también forma parte del crecimiento.	comun	\N	t	2026-05-21 21:54:32.736802
3b9b792d-98f7-4f56-9351-7e3d83698883	cansancio	energia	companero	No necesitas exigir más de lo que puedes dar hoy.	comun	\N	t	2026-05-21 21:54:32.736802
65380409-6645-4894-8f2f-e64de3dd26b6	cansancio	pausa	suave	A veces la mejor decisión es bajar el ritmo.	raro	\N	t	2026-05-21 21:54:32.736802
8f3e6093-5816-4e95-a40f-53a78a1085a7	racha	constancia	celebracion	Tu constancia está dando frutos.	comun	\N	t	2026-05-21 21:54:44.494487
e7d51059-8912-42e9-89f9-e15ba00f18bf	racha	identidad	celebracion	Cada día suma a la construcción de tu identidad.	comun	\N	t	2026-05-21 21:54:44.494487
f6bdbd2f-26a1-4ef9-b5ec-84124b8e1fbd	racha	disciplina	celebracion	Has demostrado compromiso contigo mismo.	comun	\N	t	2026-05-21 21:54:44.494487
44b4c59a-b0b2-456b-8009-dd6cc867cd51	racha	logro	celebracion	La disciplina se construye repitiendo pequeños actos.	raro	\N	t	2026-05-21 21:54:44.494487
1c26a69b-8a93-4f27-a17b-9f5961771835	conciencia	presencia	reflexivo	Cuando prestas atención a lo que ocurre dentro de ti, aparecen nuevas posibilidades.	comun	\N	t	2026-05-21 21:59:40.080807
b135ec79-7cbc-420c-a5fb-2ccf24447095	conciencia	observacion	sereno	Observar sin juzgar es una forma de aprender.	comun	\N	t	2026-05-21 21:59:40.080807
a2c3ca41-4d7b-4fe1-b7a0-c7270a0705d1	conciencia	autoconocimiento	reflexivo	Cada emoción puede enseñarte algo importante.	comun	\N	t	2026-05-21 21:59:40.080807
fdfc10ed-65e3-4af3-8824-b1d63fb59165	conciencia	claridad	sereno	La claridad suele llegar cuando dejamos de correr por un momento.	raro	\N	t	2026-05-21 21:59:40.080807
27967fd8-d26c-4a5f-a69a-942e315f5bf4	fuerza	accion	energico	Tu fuerza crece cada vez que decides actuar.	comun	\N	t	2026-05-21 21:59:49.863947
17f55f12-0565-4f2f-841f-1bc0c65b378a	fuerza	movimiento	energico	El movimiento de hoy construye la energía de mañana.	comun	\N	t	2026-05-21 21:59:49.863947
bb2e52b3-b151-4d63-8942-c84e55419ff0	fuerza	progreso	calido	Incluso los avances pequeños fortalecen tu camino.	comun	\N	t	2026-05-21 21:59:49.863947
e16c2426-3979-4d4e-83f0-80243b0c0cd4	fuerza	desafio	celebracion	Superar un desafío deja una huella de fortaleza.	raro	\N	t	2026-05-21 21:59:49.863947
e8a0a771-845c-4d18-a3cd-5b32e1694445	resistencia	constancia	calido	Persistir también es una forma de valentía.	comun	\N	t	2026-05-21 22:00:02.42955
48b5e797-8f6f-443e-8a09-e39627221ee7	resistencia	dificultad	sereno	No necesitas avanzar rápido, solo seguir avanzando.	comun	\N	t	2026-05-21 22:00:02.42955
46249506-baa6-4f0b-aa79-8d23c340c6d5	resistencia	perseverancia	reflexivo	La resistencia se construye en los días que parecen más difíciles.	comun	\N	t	2026-05-21 22:00:02.42955
5e8f12ec-dc50-4285-ae70-ae6e455397c5	resistencia	superacion	celebracion	Cada obstáculo superado fortalece tu capacidad de continuar.	raro	\N	t	2026-05-21 22:00:02.42955
0e79b33c-0179-41d7-baad-cb51dc4aae3e	evolucion_joven	cambio_etapa	emocionado	He cambiado un poco. Gracias por ayudarme a crecer.	raro	\N	t	2026-06-24 20:39:33.069798
6cc97c0f-5be0-4ad0-9d35-011943ef8f28	evolucion_joven	cambio_etapa	calido	Cada pequeño paso que diste me trajo hasta aquí.	raro	\N	t	2026-06-24 20:39:33.069798
253d6c49-e47a-4ce1-ae29-dbc56ba43fa2	evolucion_joven	cambio_etapa	reflexivo	Ya no soy el mismo de antes. Y vos tampoco.	raro	\N	t	2026-06-24 20:39:33.069798
64b53f84-2bbf-4c61-83e8-ab6ae207a93b	evolucion_joven	cambio_etapa	companero	Estamos creciendo juntos.	raro	\N	t	2026-06-24 20:39:33.069798
ff20db71-a051-4ebf-9e76-deb34e233a92	evolucion_joven	cambio_etapa	calido	Recuerdo cuando recién comenzamos este camino. Mirá todo lo que logramos.	epico	\N	t	2026-06-24 20:39:33.069798
5b4360d6-66f9-4f17-9706-d921689fec46	evolucion_joven	cambio_etapa	celebracion	Este cambio es mío, pero también es tuyo.	epico	\N	t	2026-06-24 20:39:33.069798
7788df02-cc9b-4b6e-b6f5-5b7a2edfcfea	bienvenida	primer_ingreso	calido	Hola. Soy Kai. Gracias por invitarme a acompañarte.	comun	\N	t	2026-06-24 20:40:20.11326
47128881-a5ff-4971-b173-c054903b3b5c	bienvenida	primer_ingreso	companero	Qué bueno verte por aquí. A partir de hoy vamos a construir pequeños cambios, un paso a la vez.	comun	\N	t	2026-06-24 20:40:20.11326
fe80086a-6cdc-4d19-b2b6-897e5ac82346	bienvenida	primer_ingreso	calido	Hola. Soy Kai. Voy a estar con vos en cada pequeño avance, incluso en los días difíciles.	comun	\N	t	2026-06-24 20:40:20.11326
a6749aeb-a4c2-4454-8555-5926d8b30570	bienvenida	primer_ingreso	sereno	Bienvenido. Este lugar no es para competir con nadie. Es para crecer a tu propio ritmo.	comun	\N	t	2026-06-24 20:40:20.11326
9362d2ea-d2a7-48da-8aa7-2109918829d5	bienvenida	onboarding	companero	Cada hábito que completes me ayudará a crecer. Y mientras yo crezco, vos también.	raro	\N	t	2026-06-24 20:40:20.11326
a0d9a5d4-b9f1-4cca-8564-1def0f42db96	bienvenida	onboarding	calido	No importa cuán pequeño parezca un paso. Los cambios importantes suelen empezar así.	comun	\N	t	2026-06-24 20:40:20.11326
2c65231a-c4c6-458f-aadf-3d3b8795544d	bienvenida	onboarding	companero	Tus acciones van construyendo mi evolución. Y también la tuya.	raro	\N	t	2026-06-24 20:40:20.11326
11ab0fc4-f5fb-4480-a99e-bd438849a134	bienvenida	primer_ingreso	calido	Me alegra que estés aquí. Hoy puede ser un buen día para empezar algo nuevo.	comun	\N	t	2026-06-24 20:40:20.11326
1d376c06-9559-4b39-8dab-fda8f3f25267	bienvenida	primer_ingreso	companero	Gracias por darme una oportunidad. Vamos a descubrir juntos de qué somos capaces.	raro	\N	t	2026-06-24 20:40:20.11326
978969fa-b664-4280-b487-3ec9067b183b	saludo	apertura	calido	Hola. Me alegra volver a verte.	comun	\N	t	2026-06-24 20:40:45.514595
93ce4488-9797-4312-aeb8-90701b30b363	saludo	apertura	companero	¿Qué pequeño paso te gustaría dar hoy?	comun	\N	t	2026-06-24 20:40:45.514595
dc36452a-6eb8-4681-b5c1-89c9cf6c3790	saludo	apertura	calido	Estoy listo para acompañarte un día más.	comun	\N	t	2026-06-24 20:40:45.514595
add1cb93-e4e4-4ec4-b799-bef893d2e2ae	saludo	apertura	sereno	Un nuevo día también es una nueva oportunidad para intentarlo.	comun	\N	t	2026-06-24 20:40:45.514595
a85f7b10-1d73-4e2c-b4f3-7632ddba4e2a	saludo	apertura	companero	No importa cómo haya ido ayer. Hoy podemos empezar de nuevo.	raro	\N	t	2026-06-24 20:40:45.514595
2389e918-a23a-426d-9aea-2a94b7a0be75	saludo	apertura	calido	Gracias por volver. Me gusta compartir este camino con vos.	raro	\N	t	2026-06-24 20:40:45.514595
7386849a-dfb8-4ecf-aa74-598e4391c52a	evolucion_adulto	cambio_etapa	celebracion	He alcanzado una nueva etapa.	epico	\N	t	2026-06-24 20:41:23.051268
33156f9e-76a8-4403-9ff0-b8a943e1981c	evolucion_adulto	cambio_etapa	reflexivo	Todo lo que logramos hasta aquí fue gracias a tu constancia.	epico	\N	t	2026-06-24 20:41:23.051268
c77685ed-32ee-43a9-a128-5df7645140ef	evolucion_adulto	cambio_etapa	sabio	Los cambios importantes ocurren paso a paso.	raro	\N	t	2026-06-24 20:41:23.051268
d0abfe1a-49bc-4b64-98f1-4222423cb74b	evolucion_adulto	cambio_etapa	calido	Gracias por acompañarme en este camino.	raro	\N	t	2026-06-24 20:41:23.051268
bffdac93-5734-4639-b688-e7a4257c2356	evolucion_adulto	cambio_etapa	reflexivo	Ahora veo el mundo de una forma diferente.	epico	\N	t	2026-06-24 20:41:23.051268
6999f1b0-864b-44dd-b0ec-00e84297f000	evolucion_adulto	cambio_etapa	calido	Tu esfuerzo dejó huellas en mí.	epico	\N	t	2026-06-24 20:41:23.051268
eb3d4cf0-97ba-4fda-b24d-b78ef758df52	evolucion_adulto	cambio_etapa	sabio	Lo que construimos juntos sigue creciendo.	epico	\N	t	2026-06-24 20:41:23.051268
020c8afb-0c9f-4f32-aae9-3d6400c76bd1	evolucion_adulto	cambio_etapa	celebracion	Esto es solo el comienzo.	epico	\N	t	2026-06-24 20:41:23.051268
\.


--
-- TOC entry 5234 (class 0 OID 60165)
-- Dependencies: 236
-- Data for Name: mensajes_usuario; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mensajes_usuario (id, usuario_id, mensaje_kai_id, leido, mostrado_en, created_at) FROM stdin;
\.


--
-- TOC entry 5219 (class 0 OID 59699)
-- Dependencies: 221
-- Data for Name: misiones_diarias; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.misiones_diarias (id, usuario_id, titulo, descripcion, area, xp_recompensa, es_mision_minima, completada, fecha, created_at, habito_usuario_id) FROM stdin;
53a3ad0a-cd8b-47d5-b037-0311730655e1	87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7	Caminar 5 minutos	Una caminata pequeña ya es movimiento.	movimiento	10	t	f	2026-05-06	2026-05-06 20:16:06.143694	\N
\.


--
-- TOC entry 5230 (class 0 OID 60091)
-- Dependencies: 232
-- Data for Name: momentos_memorables; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.momentos_memorables (id, usuario_id, titulo, descripcion, tipo_momento, imagen_referencia, reflexion, relacionado_con, importante, fecha_momento, created_at, encuentro_animal_id, mensaje_usuario_id, evento_usuario_id) FROM stdin;
0bc7d514-94cb-4e98-b2f7-94fd2465a670	87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7	El día que volviste	Regresaste incluso después de sentirte agotada.	regreso	lobo_regreso	Volver también es una forma de avanzar.	lobo	t	2026-05-06 21:48:57.701385	2026-05-06 21:48:57.701385	\N	\N	\N
\.


--
-- TOC entry 5225 (class 0 OID 59981)
-- Dependencies: 227
-- Data for Name: onboarding_respuestas; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.onboarding_respuestas (id, usuario_id, estado_actual, dificultad_principal, necesidad_principal, ritmo_vida, autoexigencia, acto_cuidado_inicial, created_at) FROM stdin;
\.


--
-- TOC entry 5221 (class 0 OID 59738)
-- Dependencies: 223
-- Data for Name: progreso_fisico; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.progreso_fisico (id, usuario_id, fecha, peso, cintura, cadera, pecho, brazo, muslo, pantorrilla, created_at, evento_usuario_id) FROM stdin;
a54cc285-a6b6-4edb-bb79-76f53958c422	87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7	2026-05-06	81.50	92.0	104.0	98.0	32.0	58.0	37.0	2026-05-06 20:26:00.320681	\N
\.


--
-- TOC entry 5229 (class 0 OID 60066)
-- Dependencies: 231
-- Data for Name: rachas; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rachas (id, usuario_id, habito_usuario_id, dias_actuales, record_historico, inicio_racha, ultima_actividad, activa, protegida, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 5228 (class 0 OID 60041)
-- Dependencies: 230
-- Data for Name: registros_habito; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.registros_habito (id, usuario_id, habito_usuario_id, fecha, completado, valor_registrado, xp_ganada, created_at) FROM stdin;
\.


--
-- TOC entry 5236 (class 0 OID 60229)
-- Dependencies: 238
-- Data for Name: tipos_atributo_kai; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tipos_atributo_kai (id, nombre, descripcion) FROM stdin;
84eda4c2-31a8-463b-b272-75a7fbf6421d	Vitalidad	Representa energía y cuidado físico
0d90fffa-fb4e-4905-91ce-51869671bbcb	Fuerza	Representa capacidad de acción y movimiento
42a645e9-7064-41b6-93ca-2a14266aee6d	Resistencia	Representa perseverancia frente a dificultades
da217a67-56c5-4a19-aa76-f97663a2a065	Disciplina	Representa constancia y hábitos sostenidos
c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	Sabiduria	Representa reflexión y aprendizaje
d0962734-7507-44d7-b276-2374be80c9f0	Equilibrio	Representa estabilidad emocional
a9c2eeb4-e786-46ba-bfc5-16e978db22d9	Conciencia	Representa autoconocimiento y presencia
e7319970-d7f0-4abc-a069-63e430c75a95	Vinculo	Representa la relación entre el usuario y Kai
\.


--
-- TOC entry 5218 (class 0 OID 59606)
-- Dependencies: 220
-- Data for Name: usuarios; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.usuarios (id, nombre, email, password_hash, fecha_registro, perfil_base, etapa_kai, racha_global, dias_inactivo) FROM stdin;
87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7	Roma	roma@gmail.com	123456	2026-05-06 19:12:14.793005	retorno	cachorro	0	0
\.


--
-- TOC entry 5240 (class 0 OID 60300)
-- Dependencies: 242
-- Data for Name: xp_atributos; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.xp_atributos (id, categoria_xp_id, atributo_kai_id, multiplicador) FROM stdin;
3fffad7d-0e6f-4801-8910-f8e567f66a89	08f90c37-7224-4da5-b796-d83b323be9ae	84eda4c2-31a8-463b-b272-75a7fbf6421d	1.00
b1260adf-ae2a-4594-94a2-91b19d06edd5	320a374a-6c0a-49d3-ba51-e646ca27d906	0d90fffa-fb4e-4905-91ce-51869671bbcb	1.00
34f13b17-1b81-4207-b6d5-10161e97da2f	41d43ad1-0f76-41a4-8298-b41f374f1f11	da217a67-56c5-4a19-aa76-f97663a2a065	1.00
6dc17a0c-8262-41e1-88eb-c40d675606b1	c02beabc-14c5-47c1-8278-35786865367f	c5b82b1a-f747-43c9-8e01-dd3889e8fe2e	1.00
e28ffbd5-8d03-4ea0-b661-4a876b05a815	2b87c63a-cfa7-4a30-a73d-d98b80eb0a6c	d0962734-7507-44d7-b276-2374be80c9f0	1.00
0a0537fb-e69a-43c1-abcc-62e4249cad1a	17dcb0db-37dc-42bb-8865-394e669e9572	42a645e9-7064-41b6-93ca-2a14266aee6d	1.00
2bc14b4e-4231-4a9f-b835-11b4ddf9557b	5dd6a8b4-18b3-4cea-aef8-c0d8cf7fe9e5	a9c2eeb4-e786-46ba-bfc5-16e978db22d9	1.00
880c957f-2b9e-45f0-b00e-0457b2d9d1f7	c770a520-57aa-4143-9cb1-098d425584bd	e7319970-d7f0-4abc-a069-63e430c75a95	1.00
\.


--
-- TOC entry 5241 (class 0 OID 60331)
-- Dependencies: 243
-- Data for Name: xp_usuario; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.xp_usuario (id, usuario_id, categoria_xp_id, valor, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 4984 (class 2606 OID 59859)
-- Name: animales_catalogo animales_catalogo_nombre_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.animales_catalogo
    ADD CONSTRAINT animales_catalogo_nombre_key UNIQUE (nombre);


--
-- TOC entry 4986 (class 2606 OID 59857)
-- Name: animales_catalogo animales_catalogo_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.animales_catalogo
    ADD CONSTRAINT animales_catalogo_pkey PRIMARY KEY (id);


--
-- TOC entry 5026 (class 2606 OID 60289)
-- Name: atributos_kai atributos_kai_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.atributos_kai
    ADD CONSTRAINT atributos_kai_pkey PRIMARY KEY (id);


--
-- TOC entry 5014 (class 2606 OID 60228)
-- Name: categorias_xp categorias_xp_nombre_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categorias_xp
    ADD CONSTRAINT categorias_xp_nombre_key UNIQUE (nombre);


--
-- TOC entry 5016 (class 2606 OID 60226)
-- Name: categorias_xp categorias_xp_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categorias_xp
    ADD CONSTRAINT categorias_xp_pkey PRIMARY KEY (id);


--
-- TOC entry 5006 (class 2606 OID 60140)
-- Name: configuracion_usuario configuracion_usuario_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.configuracion_usuario
    ADD CONSTRAINT configuracion_usuario_pkey PRIMARY KEY (id);


--
-- TOC entry 5008 (class 2606 OID 60142)
-- Name: configuracion_usuario configuracion_usuario_usuario_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.configuracion_usuario
    ADD CONSTRAINT configuracion_usuario_usuario_id_key UNIQUE (usuario_id);


--
-- TOC entry 5024 (class 2606 OID 60272)
-- Name: desbloqueos_kai desbloqueos_kai_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.desbloqueos_kai
    ADD CONSTRAINT desbloqueos_kai_pkey PRIMARY KEY (id);


--
-- TOC entry 4976 (class 2606 OID 59732)
-- Name: diario_estoico diario_estoico_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.diario_estoico
    ADD CONSTRAINT diario_estoico_pkey PRIMARY KEY (id);


--
-- TOC entry 4988 (class 2606 OID 59874)
-- Name: encuentros_animales encuentros_animales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.encuentros_animales
    ADD CONSTRAINT encuentros_animales_pkey PRIMARY KEY (id);


--
-- TOC entry 4980 (class 2606 OID 59769)
-- Name: estado_kai estado_kai_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_pkey PRIMARY KEY (id);


--
-- TOC entry 4982 (class 2606 OID 59771)
-- Name: estado_kai estado_kai_usuario_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_usuario_id_key UNIQUE (usuario_id);


--
-- TOC entry 5010 (class 2606 OID 60159)
-- Name: eventos_usuario eventos_usuario_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eventos_usuario
    ADD CONSTRAINT eventos_usuario_pkey PRIMARY KEY (id);


--
-- TOC entry 4994 (class 2606 OID 60017)
-- Name: habitos_catalogo habitos_catalogo_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.habitos_catalogo
    ADD CONSTRAINT habitos_catalogo_pkey PRIMARY KEY (id);


--
-- TOC entry 4996 (class 2606 OID 60030)
-- Name: habitos_usuario habitos_usuario_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.habitos_usuario
    ADD CONSTRAINT habitos_usuario_pkey PRIMARY KEY (id);


--
-- TOC entry 5022 (class 2606 OID 60251)
-- Name: mensaje_atributos mensaje_atributos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mensaje_atributos
    ADD CONSTRAINT mensaje_atributos_pkey PRIMARY KEY (id);


--
-- TOC entry 5004 (class 2606 OID 60123)
-- Name: mensajes_kai mensajes_kai_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mensajes_kai
    ADD CONSTRAINT mensajes_kai_pkey PRIMARY KEY (id);


--
-- TOC entry 5012 (class 2606 OID 60176)
-- Name: mensajes_usuario mensajes_usuario_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mensajes_usuario
    ADD CONSTRAINT mensajes_usuario_pkey PRIMARY KEY (id);


--
-- TOC entry 4974 (class 2606 OID 59714)
-- Name: misiones_diarias misiones_diarias_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.misiones_diarias
    ADD CONSTRAINT misiones_diarias_pkey PRIMARY KEY (id);


--
-- TOC entry 5002 (class 2606 OID 60104)
-- Name: momentos_memorables momentos_memorables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_pkey PRIMARY KEY (id);


--
-- TOC entry 4990 (class 2606 OID 59995)
-- Name: onboarding_respuestas onboarding_respuestas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.onboarding_respuestas
    ADD CONSTRAINT onboarding_respuestas_pkey PRIMARY KEY (id);


--
-- TOC entry 4992 (class 2606 OID 59997)
-- Name: onboarding_respuestas onboarding_respuestas_usuario_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.onboarding_respuestas
    ADD CONSTRAINT onboarding_respuestas_usuario_id_key UNIQUE (usuario_id);


--
-- TOC entry 4978 (class 2606 OID 59747)
-- Name: progreso_fisico progreso_fisico_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progreso_fisico
    ADD CONSTRAINT progreso_fisico_pkey PRIMARY KEY (id);


--
-- TOC entry 5000 (class 2606 OID 60080)
-- Name: rachas rachas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rachas
    ADD CONSTRAINT rachas_pkey PRIMARY KEY (id);


--
-- TOC entry 4998 (class 2606 OID 60055)
-- Name: registros_habito registros_habito_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.registros_habito
    ADD CONSTRAINT registros_habito_pkey PRIMARY KEY (id);


--
-- TOC entry 5018 (class 2606 OID 60240)
-- Name: tipos_atributo_kai tipos_atributo_kai_nombre_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tipos_atributo_kai
    ADD CONSTRAINT tipos_atributo_kai_nombre_key UNIQUE (nombre);


--
-- TOC entry 5020 (class 2606 OID 60238)
-- Name: tipos_atributo_kai tipos_atributo_kai_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tipos_atributo_kai
    ADD CONSTRAINT tipos_atributo_kai_pkey PRIMARY KEY (id);


--
-- TOC entry 5028 (class 2606 OID 68574)
-- Name: atributos_kai uq_atributos_kai_usuario_atributo; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.atributos_kai
    ADD CONSTRAINT uq_atributos_kai_usuario_atributo UNIQUE (usuario_id, atributo_kai_id);


--
-- TOC entry 5032 (class 2606 OID 68576)
-- Name: xp_usuario uq_xp_usuario_categoria; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xp_usuario
    ADD CONSTRAINT uq_xp_usuario_categoria UNIQUE (usuario_id, categoria_xp_id);


--
-- TOC entry 4970 (class 2606 OID 59623)
-- Name: usuarios usuarios_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_email_key UNIQUE (email);


--
-- TOC entry 4972 (class 2606 OID 59621)
-- Name: usuarios usuarios_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_pkey PRIMARY KEY (id);


--
-- TOC entry 5030 (class 2606 OID 60309)
-- Name: xp_atributos xp_atributos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xp_atributos
    ADD CONSTRAINT xp_atributos_pkey PRIMARY KEY (id);


--
-- TOC entry 5034 (class 2606 OID 60342)
-- Name: xp_usuario xp_usuario_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xp_usuario
    ADD CONSTRAINT xp_usuario_pkey PRIMARY KEY (id);


--
-- TOC entry 5065 (class 2606 OID 60295)
-- Name: atributos_kai atributos_kai_atributo_kai_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.atributos_kai
    ADD CONSTRAINT atributos_kai_atributo_kai_id_fkey FOREIGN KEY (atributo_kai_id) REFERENCES public.tipos_atributo_kai(id);


--
-- TOC entry 5066 (class 2606 OID 60290)
-- Name: atributos_kai atributos_kai_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.atributos_kai
    ADD CONSTRAINT atributos_kai_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5057 (class 2606 OID 60143)
-- Name: configuracion_usuario configuracion_usuario_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.configuracion_usuario
    ADD CONSTRAINT configuracion_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5063 (class 2606 OID 60363)
-- Name: desbloqueos_kai desbloqueos_kai_evento_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.desbloqueos_kai
    ADD CONSTRAINT desbloqueos_kai_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
-- TOC entry 5064 (class 2606 OID 60273)
-- Name: desbloqueos_kai desbloqueos_kai_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.desbloqueos_kai
    ADD CONSTRAINT desbloqueos_kai_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5037 (class 2606 OID 60353)
-- Name: diario_estoico diario_estoico_evento_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.diario_estoico
    ADD CONSTRAINT diario_estoico_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
-- TOC entry 5038 (class 2606 OID 59733)
-- Name: diario_estoico diario_estoico_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.diario_estoico
    ADD CONSTRAINT diario_estoico_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5043 (class 2606 OID 59880)
-- Name: encuentros_animales encuentros_animales_animal_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.encuentros_animales
    ADD CONSTRAINT encuentros_animales_animal_id_fkey FOREIGN KEY (animal_id) REFERENCES public.animales_catalogo(id);


--
-- TOC entry 5044 (class 2606 OID 59875)
-- Name: encuentros_animales encuentros_animales_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.encuentros_animales
    ADD CONSTRAINT encuentros_animales_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5041 (class 2606 OID 60325)
-- Name: estado_kai estado_kai_atributo_dominante_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_atributo_dominante_id_fkey FOREIGN KEY (atributo_dominante_id) REFERENCES public.tipos_atributo_kai(id);


--
-- TOC entry 5042 (class 2606 OID 59772)
-- Name: estado_kai estado_kai_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5058 (class 2606 OID 60160)
-- Name: eventos_usuario eventos_usuario_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.eventos_usuario
    ADD CONSTRAINT eventos_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5046 (class 2606 OID 60320)
-- Name: habitos_catalogo habitos_catalogo_categoria_xp_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.habitos_catalogo
    ADD CONSTRAINT habitos_catalogo_categoria_xp_id_fkey FOREIGN KEY (categoria_xp_id) REFERENCES public.categorias_xp(id);


--
-- TOC entry 5047 (class 2606 OID 60036)
-- Name: habitos_usuario habitos_usuario_habito_catalogo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.habitos_usuario
    ADD CONSTRAINT habitos_usuario_habito_catalogo_id_fkey FOREIGN KEY (habito_catalogo_id) REFERENCES public.habitos_catalogo(id);


--
-- TOC entry 5048 (class 2606 OID 60031)
-- Name: habitos_usuario habitos_usuario_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.habitos_usuario
    ADD CONSTRAINT habitos_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5061 (class 2606 OID 60257)
-- Name: mensaje_atributos mensaje_atributos_atributo_kai_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mensaje_atributos
    ADD CONSTRAINT mensaje_atributos_atributo_kai_id_fkey FOREIGN KEY (atributo_kai_id) REFERENCES public.tipos_atributo_kai(id);


--
-- TOC entry 5062 (class 2606 OID 60252)
-- Name: mensaje_atributos mensaje_atributos_mensaje_kai_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mensaje_atributos
    ADD CONSTRAINT mensaje_atributos_mensaje_kai_id_fkey FOREIGN KEY (mensaje_kai_id) REFERENCES public.mensajes_kai(id);


--
-- TOC entry 5059 (class 2606 OID 60182)
-- Name: mensajes_usuario mensajes_usuario_mensaje_kai_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mensajes_usuario
    ADD CONSTRAINT mensajes_usuario_mensaje_kai_id_fkey FOREIGN KEY (mensaje_kai_id) REFERENCES public.mensajes_kai(id);


--
-- TOC entry 5060 (class 2606 OID 60177)
-- Name: mensajes_usuario mensajes_usuario_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mensajes_usuario
    ADD CONSTRAINT mensajes_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5035 (class 2606 OID 60187)
-- Name: misiones_diarias misiones_diarias_habito_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.misiones_diarias
    ADD CONSTRAINT misiones_diarias_habito_usuario_id_fkey FOREIGN KEY (habito_usuario_id) REFERENCES public.habitos_usuario(id);


--
-- TOC entry 5036 (class 2606 OID 59715)
-- Name: misiones_diarias misiones_diarias_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.misiones_diarias
    ADD CONSTRAINT misiones_diarias_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5053 (class 2606 OID 60192)
-- Name: momentos_memorables momentos_memorables_encuentro_animal_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_encuentro_animal_id_fkey FOREIGN KEY (encuentro_animal_id) REFERENCES public.encuentros_animales(id);


--
-- TOC entry 5054 (class 2606 OID 60202)
-- Name: momentos_memorables momentos_memorables_evento_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
-- TOC entry 5055 (class 2606 OID 60197)
-- Name: momentos_memorables momentos_memorables_mensaje_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_mensaje_usuario_id_fkey FOREIGN KEY (mensaje_usuario_id) REFERENCES public.mensajes_usuario(id);


--
-- TOC entry 5056 (class 2606 OID 60105)
-- Name: momentos_memorables momentos_memorables_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5045 (class 2606 OID 59998)
-- Name: onboarding_respuestas onboarding_respuestas_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.onboarding_respuestas
    ADD CONSTRAINT onboarding_respuestas_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5039 (class 2606 OID 60358)
-- Name: progreso_fisico progreso_fisico_evento_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progreso_fisico
    ADD CONSTRAINT progreso_fisico_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
-- TOC entry 5040 (class 2606 OID 59748)
-- Name: progreso_fisico progreso_fisico_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progreso_fisico
    ADD CONSTRAINT progreso_fisico_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5051 (class 2606 OID 60086)
-- Name: rachas rachas_habito_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rachas
    ADD CONSTRAINT rachas_habito_usuario_id_fkey FOREIGN KEY (habito_usuario_id) REFERENCES public.habitos_usuario(id);


--
-- TOC entry 5052 (class 2606 OID 60081)
-- Name: rachas rachas_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rachas
    ADD CONSTRAINT rachas_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5049 (class 2606 OID 60061)
-- Name: registros_habito registros_habito_habito_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.registros_habito
    ADD CONSTRAINT registros_habito_habito_usuario_id_fkey FOREIGN KEY (habito_usuario_id) REFERENCES public.habitos_usuario(id);


--
-- TOC entry 5050 (class 2606 OID 60056)
-- Name: registros_habito registros_habito_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.registros_habito
    ADD CONSTRAINT registros_habito_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
-- TOC entry 5067 (class 2606 OID 60315)
-- Name: xp_atributos xp_atributos_atributo_kai_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xp_atributos
    ADD CONSTRAINT xp_atributos_atributo_kai_id_fkey FOREIGN KEY (atributo_kai_id) REFERENCES public.tipos_atributo_kai(id);


--
-- TOC entry 5068 (class 2606 OID 60310)
-- Name: xp_atributos xp_atributos_categoria_xp_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xp_atributos
    ADD CONSTRAINT xp_atributos_categoria_xp_id_fkey FOREIGN KEY (categoria_xp_id) REFERENCES public.categorias_xp(id);


--
-- TOC entry 5069 (class 2606 OID 60348)
-- Name: xp_usuario xp_usuario_categoria_xp_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xp_usuario
    ADD CONSTRAINT xp_usuario_categoria_xp_id_fkey FOREIGN KEY (categoria_xp_id) REFERENCES public.categorias_xp(id);


--
-- TOC entry 5070 (class 2606 OID 60343)
-- Name: xp_usuario xp_usuario_usuario_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xp_usuario
    ADD CONSTRAINT xp_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


-- Completed on 2026-06-24 20:51:14

--
-- PostgreSQL database dump complete
--

\unrestrict oJfded7zOp2pVAWm1W9Uv6x5WEk5GFZ781NlmWPEISPM6LMbJpHEYOZRWuaGQ2O

