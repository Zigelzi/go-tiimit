package migrations

import "embed"

//go:embed schema/*.sql
var embeddedMigrations embed.FS

func GetMigrationFS() embed.FS {
	return embeddedMigrations
}
