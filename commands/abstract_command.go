package commands

import (
	"github.com/malyg1n/sql-migrator/services"
)

const (
	exitStatusSuccess = 0
	exitStatusError   = 1
)

// AbstractCommand
type AbstractCommand struct {
	srv services.ServiceContract
}

//func (c *AbstractCommand) checkMigrationsTable() error {
//	row, err := c.db.Query(fmt.Sprintf("SELECT COUNT(id) from %s", migrationTableName))
//	if err != nil || row == nil {
//		if err := c.createMigrationsTable(); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (c *AbstractCommand) createMigrationsTable() error {
//	sqlQuery := fmt.Sprintf(`CREATE TABLE %s (
//    id bigserial not null primary key,
//    migration varchar not null unique,
//    version int not null,
//    created_at timestamp DEFAULT CURRENT_TIMESTAMP
//)`, migrationTableName)
//	_, err := c.db.Query(sqlQuery)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (c *AbstractCommand) applyMigrationsUp(migrationsFiles []string, version uint) error {
//	for _, file := range migrationsFiles {
//		data, err := ioutil.ReadFile(file)
//		if err != nil {
//			return err
//		}
//		migrationQuery := fmt.Sprintf("INSERT INTO %s (migration, version) VALUES ($1, $2)", migrationTableName)
//		_, err = c.db.Query(migrationQuery, file, version)
//		if err != nil {
//			return err
//		}
//
//		query := string(data)
//		_, err = c.db.Query(query)
//		if err != nil {
//			return err
//		}
//		output.ShowMessage(fmt.Sprintf("migrated: %s", file))
//	}
//	return nil
//}
//
//func (c *AbstractCommand) getMigrationsFromBD() ([]*MigrationEntity, error) {
//	migrations := make([]*MigrationEntity, 0)
//	query := fmt.Sprintf("SELECT * from %s ORDER BY created_at DESC", migrationTableName)
//	rows, err := c.db.Query(query)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		me := &MigrationEntity{}
//		if err := rows.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt); err != nil {
//			return nil, err
//		}
//		migrations = append(migrations, me)
//	}
//
//	return migrations, nil
//}
//
//func (c *AbstractCommand) getMigrationsByVersion(version uint) ([]*MigrationEntity, error) {
//	query := fmt.Sprintf("SELECT * FROM %s WHERE version=%d ORDER BY created_at DESC", migrationTableName, version)
//	rows, err := c.db.Query(query)
//	if err != nil {
//		return nil, err
//	}
//	migrations := make([]*MigrationEntity, 0)
//	for rows.Next() {
//		me := &MigrationEntity{}
//		if err := rows.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt); err != nil {
//			return nil, err
//		}
//		migrations = append(migrations, me)
//	}
//
//	return migrations, nil
//}
//
//func (c *AbstractCommand) getMigrationUpFiles(folder string) ([]string, error) {
//	files := make([]string, 0)
//	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
//		if strings.Contains(path, "-up.sql") {
//			files = append(files, path)
//		}
//		return nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//	return files, nil
//}
//
//func (c *AbstractCommand) filterMigrations(dbMigrations []*MigrationEntity, files []string) []string {
//	newFiles := make([]string, 0)
//	for _, file := range files {
//		found := false
//		for _, m := range dbMigrations {
//			if m.Migration == file {
//				found = true
//				break
//			}
//		}
//		if found == false {
//			newFiles = append(newFiles, file)
//		}
//	}
//	return newFiles
//}
//
//func (c *AbstractCommand) getLatestVersionNumber() (uint, error) {
//	var versionNumber uint
//	me := &MigrationEntity{}
//	query := fmt.Sprintf("SELECT * FROM %s ORDER BY version DESC limit 1", migrationTableName)
//	row := c.db.QueryRow(query)
//	if row == nil {
//		versionNumber = 1
//	} else {
//		row.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt)
//		if me.Version == 0 {
//			return 1, nil
//		}
//		versionNumber = me.Version
//	}
//	return versionNumber, nil
//}
//
//func (c *AbstractCommand) rollbackMigration(entity *MigrationEntity) error {
//	filePath := strings.Replace(entity.Migration, "-up.sql", "-down.sql", 1)
//	data, err := ioutil.ReadFile(filePath)
//	if err != nil {
//		return err
//	}
//
//	downQuery := string(data)
//	_, err = c.db.Query(downQuery)
//	if err != nil {
//		return err
//	}
//
//	query := fmt.Sprintf("DELETE FROM %s WHERE id=%d", migrationTableName, entity.Id)
//	_, err = c.db.Query(query)
//	if err != nil {
//		return err
//	}
//
//	output.ShowWarning(fmt.Sprintf("rolled back: %s", entity.Migration))
//	return nil
//}
//
//func (c *AbstractCommand) checkFolder(dir string) error {
//	_, err := os.Stat(dir)
//	if !os.IsNotExist(err) {
//		return nil
//	}
//	return errors.New(fmt.Sprintf(" no such file or directory: %s", dir))
//}
