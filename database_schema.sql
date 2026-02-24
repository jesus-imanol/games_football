-- Script SQL para crear las tablas necesarias
-- Base de datos: games_football

CREATE DATABASE IF NOT EXISTS games_football;
USE games_football;

-- Tabla de retas (partidos de fútbol)
CREATE TABLE IF NOT EXISTS retas (
    id VARCHAR(36) PRIMARY KEY,
    zona_id VARCHAR(50) NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    fecha_hora DATETIME NOT NULL,
    max_jugadores INT NOT NULL DEFAULT 14,
    jugadores_actuales INT NOT NULL DEFAULT 0,
    creador_id VARCHAR(50) NOT NULL,
    creador_nombre VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_zona_id (zona_id),
    INDEX idx_fecha_hora (fecha_hora)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de jugadores en retas
CREATE TABLE IF NOT EXISTS reta_jugadores (
    id VARCHAR(36) PRIMARY KEY,
    reta_id VARCHAR(36) NOT NULL,
    usuario_id VARCHAR(50) NOT NULL,
    nombre_jugador VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reta_id) REFERENCES retas(id) ON DELETE CASCADE,
    UNIQUE KEY unique_usuario_reta (reta_id, usuario_id),
    INDEX idx_reta_id (reta_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
