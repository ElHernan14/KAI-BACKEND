-- =========================================================
-- SEED: animales_catalogo
-- =========================================================

INSERT INTO public.animales_catalogo (
    id,
    nombre,
    simbolo_principal,
    tipo_animal,
    rareza,
    descripcion,
    imagen_base,
    atributo_principal,
    activo,
    created_at
) VALUES
(
    '45bb0b0f-67a3-4456-8275-1cacdbc8b609',
    'Búho',
    'sabiduria',
    'reflexivo',
    'raro',
    'Representa la capacidad de observarse y aprender del propio camino.',
    'buho_base',
    'sabiduria',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    '247f2b44-fa44-42b4-9341-1d28e67bed0a',
    'Lobo',
    'resiliencia',
    'protector',
    'epico',
    'Representa la fuerza de seguir adelante incluso después de perderse.',
    'lobo_base',
    'resistencia',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    '364c7ebd-fafd-4cd7-9ce9-c524e4223121',
    'Zorro',
    'adaptabilidad',
    'astuto',
    'raro',
    'Representa la inteligencia para adaptarse y encontrar nuevos caminos.',
    'zorro_base',
    'conciencia',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    '3e99b7a3-0781-4c8e-b75f-6974d09e626a',
    'Ciervo',
    'calma',
    'sereno',
    'comun',
    'Representa la sensibilidad, la tranquilidad y el avance sin violencia.',
    'ciervo_base',
    'equilibrio',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    'ee9d8a80-0f1d-46e6-8d5e-b3fa27f1bce8',
    'Oso',
    'fortaleza',
    'guardian',
    'epico',
    'Representa la fuerza interior y la protección emocional.',
    'oso_base',
    'fuerza',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    '23b15958-ded3-406c-8f2d-084cbd5375cd',
    'Colibrí',
    'pequenos_pasos',
    'ligero',
    'comun',
    'Representa los pequeños avances diarios que terminan transformando una vida.',
    'colibri_base',
    'vitalidad',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    'e68eaf8e-2837-49cb-9cd2-241dc065834a',
    'Tortuga',
    'constancia',
    'persistente',
    'comun',
    'Representa el progreso lento pero constante.',
    'tortuga_base',
    'disciplina',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    'f41ae20c-e70f-4322-b4d8-718eb56d1326',
    'Cuervo',
    'transformacion',
    'mistico',
    'legendario',
    'Representa los cambios profundos que nacen después de atravesar oscuridad.',
    'cuervo_base',
    'conciencia',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    'a2fa0867-a243-4052-9dee-91a007ca3d49',
    'León',
    'coraje',
    'lider',
    'legendario',
    'Representa el valor de enfrentar el miedo y avanzar igualmente.',
    'leon_base',
    'fuerza',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    'b4cf2d23-53aa-4904-b034-e6bfe190faf0',
    'Mariposa',
    'renacimiento',
    'espiritual',
    'raro',
    'Representa la transformación personal y los nuevos comienzos.',
    'mariposa_base',
    'equilibrio',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    'c35b7a29-e728-4a54-b109-3688a07d73f5',
    'Nutria',
    'alegria',
    'social',
    'comun',
    'Representa la importancia del juego, el descanso y la conexión emocional.',
    'nutria_base',
    'vitalidad',
    true,
    '2026-05-06 21:08:55.641677'
),
(
    'f06e2460-b575-43ae-9004-e26192cf32c5',
    'Águila',
    'vision',
    'visionario',
    'legendario',
    'Representa la capacidad de ver más allá del momento presente.',
    'aguila_base',
    'sabiduria',
    true,
    '2026-05-06 21:08:55.641677'
);

-- =========================================================
-- SEED: configuracion_usuario
-- =========================================================

INSERT INTO public.configuracion_usuario (
    id,
    usuario_id,
    notificaciones_activas,
    sonidos_activos,
    mostrar_rachas,
    modo_discreto,
    intensidad_kai,
    horario_recordatorio,
    bloquear_con_pin,
    permitir_mensajes_emocionales,
    created_at,
    updated_at
) VALUES (
    'd33403fc-d72c-40fe-8736-8a6eed1f5a10',
    '87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7',
    true,
    true,
    true,
    false,
    'suave',
    NULL,
    false,
    true,
    '2026-05-07 16:47:26.504542',
    '2026-05-07 16:47:26.504542'
);

-- =========================================================
-- SEED: diario_estoico
-- =========================================================

INSERT INTO public.diario_estoico (
    id,
    usuario_id,
    fecha,
    que_hice_bien,
    que_mejorar,
    como_me_senti,
    agradecimiento_1,
    agradecimiento_2,
    agradecimiento_3,
    estado_animo,
    xp_otorgada,
    created_at,
    evento_usuario_id
) VALUES (
    'dca8924a-8329-473d-9081-605e8e78545d',
    '87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7',
    '2026-05-06',
    'Caminé hoy',
    'Dormir más temprano',
    'Me sentí tranquilo',
    'Mi familia',
    'Haber descansado',
    'La comida de hoy',
    'bien',
    15,
    '2026-05-06 20:21:40.141051',
    NULL
);

-- =========================================================
-- SEED: estado_kai
-- =========================================================

INSERT INTO public.estado_kai (
    id,
    usuario_id,
    estado_actual,
    etapa_actual,
    energia,
    imagen_kai,
    ultimo_mensaje,
    ultima_interaccion,
    dias_sin_actividad,
    modo_recuperacion,
    created_at,
    updated_at,
    modo_actual,
    atributo_dominante_id,
    nivel_vinculo,
    ultima_evolucion
) VALUES (
    '16a11ef0-38c3-4ed9-b752-71f146ccddf1',
    '87fc7cf6-8ceb-4b81-a7c5-8a1aa16796d7',
    'activo',
    'cachorro',
    100,
    'cachorro_activo',
    'Empezamos. Eso ya es algo grande para mí.',
    '2026-05-06 20:38:18.556979',
    0,
    false,
    '2026-05-06 20:38:18.556979',
    '2026-05-06 20:38:18.556979',
    NULL,
    NULL,
    1,
    NULL
);
-- =========================================================
-- SEED: mensajes de evolucion para demo
-- =========================================================

INSERT INTO public.mensajes_kai (
    id,
    tipo,
    contexto,
    tono,
    mensaje,
    rareza,
    activo
) VALUES
(
    '9e6dc15c-e7c9-469c-a0f1-f38ec3f08b01',
    'evolucion',
    'evolucion',
    'celebracion',
    '¡Kai evolucionó gracias a tus hábitos! Este cambio también refleja tu progreso.',
    'comun',
    true
),
(
    '9e6dc15c-e7c9-469c-a0f1-f38ec3f08b02',
    'evolucion',
    'evolucion',
    'celebracion',
    '¡Mirá cuánto crecieron juntos! Kai está celebrando su nueva etapa.',
    'comun',
    true
),
(
    '9e6dc15c-e7c9-469c-a0f1-f38ec3f08b03',
    'evolucion',
    'evolucion',
    'celebracion',
    'Tu constancia le dio a Kai una nueva forma. ¡Sigan avanzando!',
    'comun',
    true
)
ON CONFLICT (id) DO UPDATE SET
    tipo = EXCLUDED.tipo,
    contexto = EXCLUDED.contexto,
    tono = EXCLUDED.tono,
    mensaje = EXCLUDED.mensaje,
    rareza = EXCLUDED.rareza,
    activo = EXCLUDED.activo;
