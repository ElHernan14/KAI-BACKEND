# KAI Backend API

## Descripción del Proyecto

KAI Backend es una API REST desarrollada en Go que impulsa el ecosistema principal de la plataforma KAI: un sistema gamificado de transformación personal basado en hábitos, continuidad y evolución simbólica.

La API centraliza toda la lógica de negocio del sistema, incluyendo:

- autenticación de usuarios,
- gestión de hábitos,
- sistema de experiencia (XP),
- evolución dinámica de Kai,
- motor narrativo,
- sistema de rachas,
- eventos del usuario,
- mensajes personalizados,
- y persistencia de progresión emocional.

El backend fue diseñado bajo una arquitectura modular y escalable, enfocada en separación de responsabilidades, mantenibilidad y crecimiento futuro del sistema.

---

# Objetivo del Proyecto

El objetivo principal del backend es proveer una infraestructura robusta y escalable para soportar el MVP académico de KAI, permitiendo:

- registro y autenticación de usuarios,
- seguimiento de hábitos diarios,
- cálculo de XP,
- evolución de atributos de Kai,
- generación de estados emocionales,
- persistencia narrativa,
- y comunicación eficiente con el frontend mobile Android.

El proyecto también busca aplicar conceptos modernos de:

- arquitectura backend,
- APIs REST,
- autenticación JWT,
- middlewares,
- patrones de diseño,
- validaciones,
- manejo de errores,
- y diseño modular.

---

# Tecnologías Utilizadas

## Lenguaje Principal

- Go (Golang)

## Framework Backend

- Gin Gonic

## ORM y Base de Datos

- GORM
- PostgreSQL
- Supabase

## Autenticación y Seguridad

- JWT (JSON Web Token)
- bcrypt

## Validaciones

- go-playground/validator

## Variables de Entorno

- godotenv

## Logging

- logrus / slog (según implementación final)

## Desarrollo

- Air (hot reload)
- Docker (opcional futuro)

## Arquitectura y Patrones

- REST API
- Layered Architecture
- Repository Pattern
- Middleware Pattern
- Dependency Injection
- Modular Domain Architecture

---

# Arquitectura del Proyecto

El backend utiliza una arquitectura modular basada en dominios funcionales.

Cada módulo encapsula:

- handlers,
- services,
- repositories,
- modelos,
- DTOs,
- y lógica de negocio.

El flujo principal sigue la estructura:

HTTP Request
↓
Controller
↓
Service
↓
Repository
↓
Database

## Principios Aplicados

- Separación de responsabilidades
- Bajo acoplamiento
- Escalabilidad modular
- Reutilización de lógica
- Centralización de reglas de negocio
- Manejo global de errores
- Middleware desacoplado

---

# Filosofía Técnica de KAI

KAI no funciona únicamente como un sistema CRUD de hábitos.

El backend implementa un motor evolutivo donde:

acciones del usuario
→ generan eventos
→ otorgan XP
→ modifican atributos de Kai
→ alteran estados emocionales
→ desbloquean mensajes y contenido narrativo

Esto convierte al backend en el núcleo lógico y evolutivo del sistema.

---

# Características Principales

## Sistema de Usuarios

- Registro
- Login
- JWT Authentication
- Roles
- Sesiones protegidas

## Sistema de Hábitos

- Catálogo global de hábitos
- Hábitos personalizados por usuario
- Registro diario de hábitos
- Seguimiento de continuidad

## Sistema de XP

- XP por categorías
- XP por hábitos
- XP por eventos
- XP por continuidad

## Sistema Evolutivo de Kai

- Atributos dinámicos
- Estados emocionales
- Evolución simbólica
- Reacciones contextuales

## Sistema de Rachas

- Continuidad diaria
- Reinicio controlado
- Persistencia emocional

## Sistema Narrativo

- Mensajes personalizados
- Eventos importantes
- Historial de interacciones
- Desbloqueos

## Eventos del Usuario

- Registro de acciones importantes
- Timeline evolutivo
- Base del sistema narrativo

---

# Flujo Principal del Sistema

Usuario completa hábito
↓
Se registra acción
↓
Se genera evento
↓
Se otorga XP
↓
Se actualizan atributos Kai
↓
Se recalcula estado Kai
↓
Se generan mensajes
↓
Se verifican desbloqueos

---

# Arquitectura de Capas

## Controller Layer

Responsable de:

- recibir requests,
- validar entrada,
- devolver responses HTTP.

## Service Layer

Responsable de:

- lógica de negocio,
- motor evolutivo,
- reglas del sistema,
- cálculos y validaciones internas.

## Repository Layer

Responsable de:

- acceso a datos,
- consultas SQL,
- persistencia,
- comunicación con PostgreSQL.

---

# Estructura del Proyecto

```bash
main.go

internal/

├── config/
├── database/
├── middleware/
├── routes/

├── modules/

│   ├── auth/
│   ├── users/
│   ├── habits/
│   ├── xp/
│   ├── kai/
│   ├── messages/
│   └── events/

├── shared/
│   ├── dto/
│   ├── errors/
│   ├── response/
│   └── utils/