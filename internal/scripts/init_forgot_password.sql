CREATE TABLE password_reset_tokens (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    usuario_id uuid NOT NULL REFERENCES usuarios(id),
    token_hash text NOT NULL,
    expires_at timestamp NOT NULL,
    used_at timestamp,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_password_reset_tokens_token_hash
ON password_reset_tokens(token_hash);

CREATE INDEX idx_password_reset_tokens_usuario_id
ON password_reset_tokens(usuario_id);

CREATE TABLE password_reset_codes (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    usuario_id uuid NOT NULL REFERENCES usuarios(id),
    code_hash text NOT NULL,
    expires_at timestamp NOT NULL,
    used_at timestamp,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_password_reset_codes_usuario_id
ON password_reset_codes(usuario_id);

CREATE INDEX idx_password_reset_codes_code_hash
ON password_reset_codes(code_hash);