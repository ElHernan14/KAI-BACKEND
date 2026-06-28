ALTER TABLE public.usuarios
ADD COLUMN IF NOT EXISTS username character varying(30);

CREATE UNIQUE INDEX IF NOT EXISTS usuarios_username_unique_idx
ON public.usuarios (username)
WHERE username IS NOT NULL;

ALTER TABLE usuarios
ADD COLUMN habit_records_synced_at TIMESTAMP,
ADD COLUMN activity_sync_at TIMESTAMP;
