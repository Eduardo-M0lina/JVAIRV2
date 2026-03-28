-- Script para verificar datos disponibles en el dashboard

-- 1. Jobs por estado (abiertos)
SELECT
    'Jobs Abiertos' as tipo,
    COUNT(*) as total,
    MIN(date_received) as fecha_mas_antigua,
    MAX(date_received) as fecha_mas_reciente
FROM jobs
WHERE deleted_at IS NULL AND closed = 0;

-- 2. Jobs cerrados por período
SELECT
    'Jobs Cerrados Este Año' as tipo,
    COUNT(*) as total,
    MIN(updated_at) as fecha_mas_antigua,
    MAX(updated_at) as fecha_mas_reciente
FROM jobs
WHERE deleted_at IS NULL
    AND closed = 1
    AND updated_at >= DATE_FORMAT(NOW(), '%Y-01-01');

-- 3. Jobs cerrados últimos 30 días
SELECT
    'Jobs Cerrados Últimos 30 Días' as tipo,
    COUNT(*) as total
FROM jobs
WHERE deleted_at IS NULL
    AND closed = 1
    AND updated_at >= DATE_SUB(NOW(), INTERVAL 30 DAY);

-- 4. Invoices por período
SELECT
    'Invoices Este Año' as tipo,
    COUNT(*) as total,
    SUM(CASE WHEN status IN ('draft', 'sent') THEN 1 ELSE 0 END) as pendientes,
    SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END) as pagadas
FROM invoices
WHERE deleted_at IS NULL
    AND created_at >= DATE_FORMAT(NOW(), '%Y-01-01');

-- 5. Quotes por período
SELECT
    'Quotes Este Año' as tipo,
    COUNT(*) as total,
    SUM(CASE WHEN quote_status_id = 1 THEN 1 ELSE 0 END) as pendientes,
    SUM(CASE WHEN quote_status_id = 2 THEN 1 ELSE 0 END) as aprobadas
FROM quotes
WHERE deleted_at IS NULL
    AND created_at >= DATE_FORMAT(NOW(), '%Y-01-01');

-- 6. Warranty Claims
SELECT
    'Warranty Claims' as tipo,
    COUNT(*) as total,
    SUM(CASE WHEN warranty_claim_status_id IN (1, 2) THEN 1 ELSE 0 END) as activas
FROM warranty_claims
WHERE deleted_at IS NULL;

-- 7. Alerts no leídas
SELECT
    'Alerts No Leídas' as tipo,
    COUNT(*) as total
FROM alerts
WHERE deleted_at IS NULL AND `read` = 0;

-- 8. Tasks pendientes
SELECT
    'Tasks Pendientes' as tipo,
    COUNT(*) as total,
    SUM(CASE WHEN due_date < NOW() THEN 1 ELSE 0 END) as vencidas
FROM job_tasks
WHERE deleted_at IS NULL AND task_status_id IN (1, 2);

-- 9. Jobs por categoría (abiertos)
SELECT
    jc.label as categoria,
    COUNT(*) as total
FROM jobs j
JOIN job_categories jc ON jc.id = j.job_category_id
WHERE j.deleted_at IS NULL AND j.closed = 0
GROUP BY jc.id, jc.label
ORDER BY total DESC;

-- 10. Jobs por status (abiertos)
SELECT
    js.label as status,
    COUNT(*) as total
FROM jobs j
JOIN job_statuses js ON js.id = j.job_status_id
WHERE j.deleted_at IS NULL AND j.closed = 0
GROUP BY js.id, js.label
ORDER BY total DESC;

-- 11. Actividad reciente (últimos 7 días)
SELECT
    'Actividad Últimos 7 Días' as tipo,
    COUNT(*) as total
FROM job_activity_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);

-- 12. Técnicos con carga de trabajo
SELECT
    u.name as tecnico,
    COUNT(*) as jobs_asignados
FROM jobs j
JOIN users u ON u.id = j.user_id
WHERE j.deleted_at IS NULL
    AND j.closed = 0
    AND u.role_id = 2  -- Asumiendo que 2 es el role_id de técnicos
GROUP BY u.id, u.name
ORDER BY jobs_asignados DESC
LIMIT 10;
