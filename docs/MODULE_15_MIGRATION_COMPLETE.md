# Module 15: Alerts — Migration Complete

## Overview

This document summarizes the migration of the Alerts module from Laravel (JVAIR) to Go (JVAIRV2).

## Entity: Alert

### Database Table Structure (from `scripts/db_structure.sql`)

```sql
CREATE TABLE `alerts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `alert_type` varchar(191) NOT NULL,
  `entity_id` bigint unsigned NOT NULL,
  `entity_type` varchar(191) NOT NULL,
  `message_level` varchar(191) NOT NULL,
  `message` varchar(191) NOT NULL,
  `is_read` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `alerts_user_id_foreign` (`user_id`),
  KEY `alert_entity_index` (`entity_id`,`entity_type`),
  KEY `alerts_alert_type_index` (`alert_type`),
  KEY `alerts_is_read_index` (`is_read`),
  CONSTRAINT `alerts_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
```

### Go Entity (`pkg/domain/alert/entity.go`)

```go
type Alert struct {
    ID           int64      `json:"id"`
    UserID       *int64     `json:"userId,omitempty"`     // nullable FK → users
    AlertType    string     `json:"alertType"`            // e.g., 'system'
    EntityID     int64      `json:"entityId"`             // polymorphic ID
    EntityType   string     `json:"entityType"`           // e.g., 'job', 'call_log'
    MessageLevel string     `json:"messageLevel"`         // e.g., 'info', 'danger'
    Message      string     `json:"message"`
    IsRead       bool       `json:"isRead"`
    CreatedAt    *time.Time `json:"createdAt,omitempty"`
    UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}
```

## Endpoints Implemented

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/alerts` | List alerts with filters and pagination |
| GET | `/api/v1/alerts/{id}` | Get alert by ID |
| POST | `/api/v1/alerts` | Create a new alert |
| PUT | `/api/v1/alerts/{id}/read` | Mark a single alert as read |
| PUT | `/api/v1/alerts/read-all` | Mark all unread alerts for authenticated user as read |
| GET | `/api/v1/alerts/unread-count` | Get count of unread alerts for authenticated user |
| DELETE | `/api/v1/alerts/{id}` | Delete an alert (hard delete) |

### List Filters

- `userId` (int): Filter by user ID
- `isRead` (bool): Filter by read status (true/false)
- `alertType` (string): Filter by alert type
- `entityType` (string): Filter by entity type
- `page` (int, default: 1): Page number
- `limit` (int, default: 15): Items per page

## Design Decisions

1. **No Soft Delete**: The alerts table does not have a `deleted_at` column. All deletes are hard deletes (`DELETE FROM alerts WHERE id = ?`).

2. **Nullable user_id**: The `user_id` field is nullable in the database, so we use `*int64` in the Go entity and `sql.NullInt64` in the repository for proper handling.

3. **Polymorphic Relationship**: `entity_type` and `entity_id` form a polymorphic relationship to other entities (jobs, call_logs). We store these as simple string/int64 without resolving the related entity in this module.

4. **No Full Update**: Following the Laravel pattern, there is no full UPDATE endpoint. The only write operations after creation are `MarkAsRead` and `MarkAllRead`.

5. **User Context for Mark All/Unread Count**: The `MarkAllRead` and `UnreadCount` endpoints use the authenticated user from the JWT context, not a query parameter, for security.

6. **Route Order**: The `PUT /api/v1/alerts/read-all` route is registered before `PUT /api/v1/alerts/{id}/read` to prevent Chi from interpreting "read-all" as an ID.

## Files Created

### Domain Layer (`pkg/domain/alert/`)
- `entity.go` — Alert struct and validation
- `errors.go` — Domain errors (ErrNotFound, ErrUserNotFound)
- `repository.go` — Repository interface and ListFilters
- `usecase.go` — Service interface and constructor
- `create.go` — Create use case
- `get_by_id.go` — GetByID use case
- `list.go` — List use case
- `mark_as_read.go` — MarkAsRead use case
- `mark_all_read.go` — MarkAllRead use case
- `unread_count.go` — UnreadCount use case
- `delete.go` — Delete use case

### Repository Layer (`pkg/repository/mysql/alert/`)
- `repository.go` — Repository struct and constructor
- `create.go` — Create implementation
- `get_by_id.go` — GetByID implementation
- `list.go` — List implementation with dynamic filters
- `mark_as_read.go` — MarkAsRead implementation
- `mark_all_read.go` — MarkAllRead implementation
- `unread_count.go` — UnreadCount implementation
- `delete.go` — Delete implementation
- `adapters.go` — UserExistsChecker adapter

### Handler Layer (`pkg/rest/handler/alert/`)
- `handler.go` — HTTP handler with Swagger annotations

## Integration Points

### `cmd/api/container.go`
- Added imports for `domainAlert`, `mysqlAlert`, and `alertHandler`
- Initialized alert repository, user checker, use case, and handler
- Added `AlertHandler` to Container struct
- Passed handler to router

### `pkg/rest/router/router.go`
- Added import for `alertHandler`
- Added `alertHandler` parameter to `New()` function
- Registered alert routes via `alertHandler.RegisterRoutes(r)`

## Swagger

All endpoints have complete Swagger annotations. Regenerate with:

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs
```

## Testing

### Build Verification
```bash
go build ./...
```

### Pre-commit Hooks
```bash
pre-commit run --all-files
```

## Postman Collection

See `docs/postman_alerts_collection.json` for a ready-to-import Postman collection with all alert endpoints.
