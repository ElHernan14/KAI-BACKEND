--
-- PostgreSQL database dump
--

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3


SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET search_path TO public;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
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



--
--

CREATE TYPE public.autoexigencia_enum AS ENUM (
    'muy_duro_conmigo',
    'a_veces_me_castigo',
    'trato_de_entenderme',
    'estoy_aprendiendo_a_cuidarme'
);



--
--

CREATE TYPE public.dificultad_principal_enum AS ENUM (
    'cuidarme',
    'sostener_habitos',
    'descansar',
    'motivarme',
    'organizarme',
    'sentirme_bien_conmigo'
);



--
--

CREATE TYPE public.estado_actual_enum AS ENUM (
    'agotado',
    'desconectado',
    'ansioso',
    'estancado',
    'con_ganas_de_empezar',
    'no_lo_se'
);



--
--

CREATE TYPE public.necesidad_principal_enum AS ENUM (
    'energia',
    'calma',
    'constancia',
    'confianza',
    'conexion_conmigo',
    'bienestar_fisico'
);



--
--

CREATE TYPE public.ritmo_vida_enum AS ENUM (
    'agotador',
    'variable',
    'tranquilo',
    'dificil_de_organizar'
);



SET default_tablespace = '';

SET default_table_access_method = heap;

--
--

CREATE TABLE public.animales_catalogo (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.atributos_kai (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    usuario_id uuid NOT NULL,
    atributo_kai_id uuid NOT NULL,
    valor integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);



--
--

CREATE TABLE public.categorias_xp (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    nombre character varying(30) NOT NULL,
    descripcion text
);



--
--

CREATE TABLE public.configuracion_usuario (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.desbloqueos_kai (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    usuario_id uuid NOT NULL,
    tipo_desbloqueo character varying(50),
    referencia_id uuid,
    descripcion text,
    desbloqueado_en timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    evento_usuario_id uuid
);



--
--

CREATE TABLE public.diario_estoico (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.encuentros_animales (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.estado_kai (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.eventos_usuario (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    usuario_id uuid NOT NULL,
    tipo_evento character varying(50) NOT NULL,
    referencia_id uuid,
    referencia_tabla character varying(50),
    descripcion text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    xp_otorgada integer DEFAULT 0
);



--
--

CREATE TABLE public.habitos_catalogo (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.habitos_usuario (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT habitos_usuario_id_not_null1 NOT NULL,
    usuario_id uuid NOT NULL,
    habito_catalogo_id uuid NOT NULL,
    personalizado boolean DEFAULT false,
    activo boolean DEFAULT true,
    fecha_inicio date DEFAULT CURRENT_DATE,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);



--
--

CREATE TABLE public.mensaje_atributos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    mensaje_kai_id uuid NOT NULL,
    atributo_kai_id uuid NOT NULL,
    nivel_minimo integer DEFAULT 1
);



--
--

CREATE TABLE public.mensajes_kai (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tipo character varying(50) NOT NULL,
    contexto character varying(50),
    tono character varying(30),
    mensaje text NOT NULL,
    rareza character varying(20) DEFAULT 'comun'::character varying,
    desbloqueado_por character varying(100),
    activo boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);



--
--

CREATE TABLE public.mensajes_usuario (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    usuario_id uuid NOT NULL,
    mensaje_kai_id uuid NOT NULL,
    leido boolean DEFAULT false,
    mostrado_en timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);



--
--

CREATE TABLE public.misiones_diarias (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.momentos_memorables (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.onboarding_respuestas (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    usuario_id uuid NOT NULL,
    estado_actual public.estado_actual_enum NOT NULL,
    dificultad_principal public.dificultad_principal_enum NOT NULL,
    necesidad_principal public.necesidad_principal_enum NOT NULL,
    ritmo_vida public.ritmo_vida_enum NOT NULL,
    autoexigencia public.autoexigencia_enum NOT NULL,
    acto_cuidado_inicial public.acto_cuidado_enum NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);



--
--

CREATE TABLE public.progreso_fisico (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.rachas (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
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



--
--

CREATE TABLE public.registros_habito (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    usuario_id uuid NOT NULL,
    habito_usuario_id uuid NOT NULL,
    fecha date NOT NULL,
    completado boolean DEFAULT false,
    valor_registrado text,
    xp_ganada integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);



--
--

CREATE TABLE public.tipos_atributo_kai (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    nombre character varying(30) NOT NULL,
    descripcion text
);



--
--

CREATE TABLE public.usuarios (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    nombre character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    password_hash text NOT NULL,
    fecha_registro timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    perfil_base character varying(30),
    etapa_kai character varying(30) DEFAULT 'cachorro'::character varying,
    racha_global integer DEFAULT 0,
    dias_inactivo integer DEFAULT 0
);



--
--

CREATE TABLE public.xp_atributos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    categoria_xp_id uuid NOT NULL,
    atributo_kai_id uuid NOT NULL,
    multiplicador numeric(4,2) DEFAULT 1.0
);



--
--

CREATE TABLE public.xp_usuario (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    usuario_id uuid NOT NULL,
    categoria_xp_id uuid NOT NULL,
    valor integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);



--
--
--
--

ALTER TABLE ONLY public.animales_catalogo
    ADD CONSTRAINT animales_catalogo_nombre_key UNIQUE (nombre);


--
--

ALTER TABLE ONLY public.animales_catalogo
    ADD CONSTRAINT animales_catalogo_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.atributos_kai
    ADD CONSTRAINT atributos_kai_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.categorias_xp
    ADD CONSTRAINT categorias_xp_nombre_key UNIQUE (nombre);


--
--

ALTER TABLE ONLY public.categorias_xp
    ADD CONSTRAINT categorias_xp_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.configuracion_usuario
    ADD CONSTRAINT configuracion_usuario_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.configuracion_usuario
    ADD CONSTRAINT configuracion_usuario_usuario_id_key UNIQUE (usuario_id);


--
--

ALTER TABLE ONLY public.desbloqueos_kai
    ADD CONSTRAINT desbloqueos_kai_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.diario_estoico
    ADD CONSTRAINT diario_estoico_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.encuentros_animales
    ADD CONSTRAINT encuentros_animales_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_usuario_id_key UNIQUE (usuario_id);


--
--

ALTER TABLE ONLY public.eventos_usuario
    ADD CONSTRAINT eventos_usuario_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.habitos_catalogo
    ADD CONSTRAINT habitos_catalogo_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.habitos_usuario
    ADD CONSTRAINT habitos_usuario_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.mensaje_atributos
    ADD CONSTRAINT mensaje_atributos_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.mensajes_kai
    ADD CONSTRAINT mensajes_kai_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.mensajes_usuario
    ADD CONSTRAINT mensajes_usuario_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.misiones_diarias
    ADD CONSTRAINT misiones_diarias_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.onboarding_respuestas
    ADD CONSTRAINT onboarding_respuestas_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.onboarding_respuestas
    ADD CONSTRAINT onboarding_respuestas_usuario_id_key UNIQUE (usuario_id);


--
--

ALTER TABLE ONLY public.progreso_fisico
    ADD CONSTRAINT progreso_fisico_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.rachas
    ADD CONSTRAINT rachas_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.registros_habito
    ADD CONSTRAINT registros_habito_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.tipos_atributo_kai
    ADD CONSTRAINT tipos_atributo_kai_nombre_key UNIQUE (nombre);


--
--

ALTER TABLE ONLY public.tipos_atributo_kai
    ADD CONSTRAINT tipos_atributo_kai_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_email_key UNIQUE (email);


--
--

ALTER TABLE ONLY public.usuarios
    ADD CONSTRAINT usuarios_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.xp_atributos
    ADD CONSTRAINT xp_atributos_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.xp_usuario
    ADD CONSTRAINT xp_usuario_pkey PRIMARY KEY (id);


--
--

ALTER TABLE ONLY public.atributos_kai
    ADD CONSTRAINT atributos_kai_atributo_kai_id_fkey FOREIGN KEY (atributo_kai_id) REFERENCES public.tipos_atributo_kai(id);


--
--

ALTER TABLE ONLY public.atributos_kai
    ADD CONSTRAINT atributos_kai_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.configuracion_usuario
    ADD CONSTRAINT configuracion_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.desbloqueos_kai
    ADD CONSTRAINT desbloqueos_kai_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
--

ALTER TABLE ONLY public.desbloqueos_kai
    ADD CONSTRAINT desbloqueos_kai_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.diario_estoico
    ADD CONSTRAINT diario_estoico_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
--

ALTER TABLE ONLY public.diario_estoico
    ADD CONSTRAINT diario_estoico_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.encuentros_animales
    ADD CONSTRAINT encuentros_animales_animal_id_fkey FOREIGN KEY (animal_id) REFERENCES public.animales_catalogo(id);


--
--

ALTER TABLE ONLY public.encuentros_animales
    ADD CONSTRAINT encuentros_animales_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_atributo_dominante_id_fkey FOREIGN KEY (atributo_dominante_id) REFERENCES public.tipos_atributo_kai(id);


--
--

ALTER TABLE ONLY public.estado_kai
    ADD CONSTRAINT estado_kai_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.eventos_usuario
    ADD CONSTRAINT eventos_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.habitos_catalogo
    ADD CONSTRAINT habitos_catalogo_categoria_xp_id_fkey FOREIGN KEY (categoria_xp_id) REFERENCES public.categorias_xp(id);


--
--

ALTER TABLE ONLY public.habitos_usuario
    ADD CONSTRAINT habitos_usuario_habito_catalogo_id_fkey FOREIGN KEY (habito_catalogo_id) REFERENCES public.habitos_catalogo(id);


--
--

ALTER TABLE ONLY public.habitos_usuario
    ADD CONSTRAINT habitos_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.mensaje_atributos
    ADD CONSTRAINT mensaje_atributos_atributo_kai_id_fkey FOREIGN KEY (atributo_kai_id) REFERENCES public.tipos_atributo_kai(id);


--
--

ALTER TABLE ONLY public.mensaje_atributos
    ADD CONSTRAINT mensaje_atributos_mensaje_kai_id_fkey FOREIGN KEY (mensaje_kai_id) REFERENCES public.mensajes_kai(id);


--
--

ALTER TABLE ONLY public.mensajes_usuario
    ADD CONSTRAINT mensajes_usuario_mensaje_kai_id_fkey FOREIGN KEY (mensaje_kai_id) REFERENCES public.mensajes_kai(id);


--
--

ALTER TABLE ONLY public.mensajes_usuario
    ADD CONSTRAINT mensajes_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.misiones_diarias
    ADD CONSTRAINT misiones_diarias_habito_usuario_id_fkey FOREIGN KEY (habito_usuario_id) REFERENCES public.habitos_usuario(id);


--
--

ALTER TABLE ONLY public.misiones_diarias
    ADD CONSTRAINT misiones_diarias_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_encuentro_animal_id_fkey FOREIGN KEY (encuentro_animal_id) REFERENCES public.encuentros_animales(id);


--
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_mensaje_usuario_id_fkey FOREIGN KEY (mensaje_usuario_id) REFERENCES public.mensajes_usuario(id);


--
--

ALTER TABLE ONLY public.momentos_memorables
    ADD CONSTRAINT momentos_memorables_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.onboarding_respuestas
    ADD CONSTRAINT onboarding_respuestas_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.progreso_fisico
    ADD CONSTRAINT progreso_fisico_evento_usuario_id_fkey FOREIGN KEY (evento_usuario_id) REFERENCES public.eventos_usuario(id);


--
--

ALTER TABLE ONLY public.progreso_fisico
    ADD CONSTRAINT progreso_fisico_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.rachas
    ADD CONSTRAINT rachas_habito_usuario_id_fkey FOREIGN KEY (habito_usuario_id) REFERENCES public.habitos_usuario(id);


--
--

ALTER TABLE ONLY public.rachas
    ADD CONSTRAINT rachas_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.registros_habito
    ADD CONSTRAINT registros_habito_habito_usuario_id_fkey FOREIGN KEY (habito_usuario_id) REFERENCES public.habitos_usuario(id);


--
--

ALTER TABLE ONLY public.registros_habito
    ADD CONSTRAINT registros_habito_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);


--
--

ALTER TABLE ONLY public.xp_atributos
    ADD CONSTRAINT xp_atributos_atributo_kai_id_fkey FOREIGN KEY (atributo_kai_id) REFERENCES public.tipos_atributo_kai(id);


--
--

ALTER TABLE ONLY public.xp_atributos
    ADD CONSTRAINT xp_atributos_categoria_xp_id_fkey FOREIGN KEY (categoria_xp_id) REFERENCES public.categorias_xp(id);


--
--

ALTER TABLE ONLY public.xp_usuario
    ADD CONSTRAINT xp_usuario_categoria_xp_id_fkey FOREIGN KEY (categoria_xp_id) REFERENCES public.categorias_xp(id);


--
--

ALTER TABLE ONLY public.xp_usuario
    ADD CONSTRAINT xp_usuario_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES public.usuarios(id);



--
--


