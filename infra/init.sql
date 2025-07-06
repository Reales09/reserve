-- =============================================================================
-- INICIALIZACIÓN DE BASE DE DATOS - RESERVE APP
-- =============================================================================
-- Este script se ejecuta automáticamente al crear la base de datos por primera vez

-- Crear extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Configurar zona horaria
SET timezone TO 'America/Bogota';

-- Crear esquema para logs si es necesario
CREATE SCHEMA IF NOT EXISTS logs;

-- Mensaje de confirmación
SELECT 'Base de datos inicializada correctamente' AS status; 