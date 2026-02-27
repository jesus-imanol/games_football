-- Script SQL para crear las tablas necesarias
-- Base de datos: games_football

CREATE DATABASE IF NOT EXISTS gamefotball;
USE gamefotball;

-- ============================================================
-- Eliminar tablas en orden correcto (hijos antes que padres)
-- ============================================================
DROP TABLE IF EXISTS mensajes_reta;
DROP TABLE IF EXISTS reta_jugadores;
DROP TABLE IF EXISTS retas;
DROP TABLE IF EXISTS zonas;
DROP TABLE IF EXISTS usuarios;

-- ============================================================
-- Tabla de usuarios (Login)
-- ============================================================
CREATE TABLE usuarios (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nombre VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- Tabla de zonas geográficas
-- ============================================================
CREATE TABLE zonas (
    id VARCHAR(50) PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- Tabla de retas (partidos de fútbol)
-- ============================================================
CREATE TABLE retas (
    id VARCHAR(36) PRIMARY KEY,
    zona_id VARCHAR(50) NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    fecha_hora DATETIME NOT NULL,
    max_jugadores INT NOT NULL DEFAULT 14,
    jugadores_actuales INT NOT NULL DEFAULT 0,
    creador_id VARCHAR(36) NOT NULL,
    creador_nombre VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_zona_id (zona_id),
    INDEX idx_fecha_hora (fecha_hora),
    FOREIGN KEY (zona_id) REFERENCES zonas(id) ON DELETE CASCADE,
    FOREIGN KEY (creador_id) REFERENCES usuarios(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- Tabla de jugadores en retas
-- ============================================================
CREATE TABLE reta_jugadores (
    id VARCHAR(36) PRIMARY KEY,
    reta_id VARCHAR(36) NOT NULL,
    usuario_id VARCHAR(36) NOT NULL,
    nombre_jugador VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reta_id) REFERENCES retas(id) ON DELETE CASCADE,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    UNIQUE KEY unique_usuario_reta (reta_id, usuario_id),
    INDEX idx_reta_id (reta_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- Tabla de mensajes del chat en vivo de cada reta
-- ============================================================
CREATE TABLE mensajes_reta (
    id VARCHAR(36) PRIMARY KEY,
    reta_id VARCHAR(36) NOT NULL,
    usuario_id VARCHAR(36) NOT NULL,
    texto VARCHAR(500) NOT NULL,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reta_id) REFERENCES retas(id) ON DELETE CASCADE,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    INDEX idx_mensajes_reta_id (reta_id),
    INDEX idx_mensajes_creado_en (reta_id, creado_en ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- Datos de prueba
-- ============================================================

-- Zonas de Suchiapa
INSERT INTO zonas (id, nombre) VALUES
('suchiapa_centro', 'Suchiapa Centro'),
('suchiapa_norte',  'Suchiapa Norte'),
('suchiapa_sur',    'Suchiapa Sur');


INSERT INTO usuarios (id, username, password, nombre) VALUES
('u-001', 'jesus-imanol', '$2a$10$DaW5YJlrFdh4cyVg/p1De./Dl10IUjDMfZDXzeADqKVq4kuipJrDu', 'Jesús Imanol'),
('u-002', 'carlos-dev',   '$2a$10$DaW5YJlrFdh4cyVg/p1De./Dl10IUjDMfZDXzeADqKVq4kuipJrDu', 'Carlos Dev');
