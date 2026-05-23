-- SE AGREGAN INDICES ÚNICOS COMPUESTO PARA EVITAR DUPLICADOS
-- Índice único en xp_usuario
ALTER TABLE xp_usuario
ADD CONSTRAINT xp_usuario_unique UNIQUE (usuario_id, categoria_xp_id);

-- Índice único en atributos_kai
ALTER TABLE atributos_kai
ADD CONSTRAINT atributos_kai_unique UNIQUE (usuario_id, atributo_kai_id);
