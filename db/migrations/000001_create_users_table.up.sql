CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (    
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    
    name        VARCHAR(255) NOT NULL,    
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    active      BOOLEAN NOT NULL DEFAULT true,
    
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    deleted_at  TIMESTAMP WITH TIME ZONE
);


CREATE INDEX IF NOT EXISTS idx_users_email_active ON users (email) 
WHERE deleted_at IS NULL;

-- 2. Índice para Soft Delete
-- Como quase todas as suas queries terão "WHERE deleted_at IS NULL", 
-- um índice parcial economiza muito processamento.
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

-- 3. Índice Composto para Listagem de Usuários Ativos
-- Útil para Admin ou buscas rápidas: SELECT * FROM users WHERE active = true;
CREATE INDEX IF NOT EXISTS idx_users_active ON users (active) 
WHERE deleted_at IS NULL;

