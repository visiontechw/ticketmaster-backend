CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE events (
    -- Base Fields
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE, -- Soft delete

    -- Event Fields
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    occurs_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    location    VARCHAR(255),
    capacity    BIGINT NOT NULL DEFAULT 0,
    owner_id    UUID NOT NULL,

    CONSTRAINT check_capacity_positive CHECK (capacity >= 0)
);

-- Índice para Soft Delete (GORM padrão)
-- Melhora a performance de quase todas as queries, já que o GORM sempre filtra por deleted_at IS NULL
CREATE INDEX idx_events_deleted_at ON events (deleted_at) WHERE deleted_at IS NULL;

-- Índice para busca por Owner (FK)
CREATE INDEX idx_events_owner_id ON events (owner_id);

-- Índice para ordenação e busca por data do evento
-- Importante se você tiver uma lista de "Próximos Eventos"
CREATE INDEX idx_events_occurs_at ON events (occurs_at);



INSERT INTO events (name, description, occurs_at, location, capacity, owner_id)
VALUES 
(
    'Conferência de Tecnologia 2026', 
    'O maior evento de tecnologia do semestre.', 
    NOW() + interval '30 days', 
    'Centro de Convenções São Paulo', 
    500, 
    uuid_generate_v4()
),
(
    'Workshop Go Avançado', 
    'Práticas de concorrência e GORM.', 
    NOW() + interval '15 days', 
    'Auditório Tech', 
    50, 
    uuid_generate_v4()
),
(
    'Meetup de IA e Machine Learning', 
    'Discussão sobre modelos generativos.', 
    NOW() + interval '5 days', 
    'Rua das Startups, 123', 
    100, 
    uuid_generate_v4()
);

-- 2. Inserindo Eventos ENCERRADOS (Passados)
INSERT INTO events (name, description, occurs_at, location, capacity, owner_id)
VALUES 
(
    'Hackathon Interno 2025', 
    'Evento de inovação do ano passado.', 
    NOW() - interval '6 months', 
    'Escritório Sede', 
    200, 
    uuid_generate_v4()
),
(
    'Lançamento Versão Alpha', 
    'Celebração do primeiro deploy.', 
    NOW() - interval '1 year', 
    'Restaurante SkyBar', 
    30, 
    uuid_generate_v4()
);

-- 3. Inserindo um Evento DELETADO (Soft Delete)
-- Útil para testar se seu GetOpenedEvents realmente ignora registros com deleted_at preenchido
INSERT INTO events (name, description, occurs_at, location, capacity, owner_id, deleted_at)
VALUES 
(
    'Evento Cancelado', 
    'Este evento não deve aparecer em nenhuma listagem padrão.', 
    NOW() + interval '10 days', 
    'Local Remoto', 
    0, 
    uuid_generate_v4(),
    NOW()
);

