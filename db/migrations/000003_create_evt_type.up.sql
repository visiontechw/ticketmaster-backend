-- 1. Preparação do ambiente
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Criar primeiro a tabela independente (Pai)
CREATE TABLE eventtype (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- 3. Criar a coluna na tabela dependente (Filho)
-- Nota: Adicionamos sem o NOT NULL primeiro caso a tabela já tenha dados
ALTER TABLE events ADD COLUMN IF NOT EXISTS eventtypeId UUID;

-- 4. Criar os Índices
-- Índice para soft delete na tabela de tipos
CREATE INDEX idx_eventtype_deleted_at ON eventtype(deleted_at);

-- Índice na tabela events (Removi o CONCURRENTLY para scripts de migration padrão, 
-- pois ele não pode rodar dentro de transações BEGIN/COMMIT)
CREATE INDEX idx_events_eventtypeId ON events(eventtypeId);

-- 5. Estabelecer a Constraint de Chave Estrangeira
ALTER TABLE events 
ADD CONSTRAINT fk_events_event_type 
FOREIGN KEY (eventtypeId) REFERENCES eventtype(id)
ON DELETE RESTRICT;

-- 6. Inserir os dados iniciais (Seed)
INSERT INTO eventtype (id, name) 
VALUES 
    (uuid_generate_v4(), 'Show'),
    (uuid_generate_v4(), 'Palestras'),
    (uuid_generate_v4(), 'Workshop'),
    (uuid_generate_v4(), 'Religioso')
ON CONFLICT (name) DO NOTHING;