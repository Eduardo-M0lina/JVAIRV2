package permission

// Permission representa un permiso en el sistema
type Permission struct {
	ID         int64  // bigint unsigned NOT NULL AUTO_INCREMENT
	AbilityID  int64  // bigint unsigned NOT NULL
	EntityID   int64  // bigint unsigned DEFAULT NULL
	EntityType string // varchar(191) DEFAULT NULL
	Forbidden  bool   // tinyint(1) NOT NULL DEFAULT '0'
	Scope      *int   // int DEFAULT NULL
}
