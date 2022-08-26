package migrator

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/anton7r/mgx/config"
	"github.com/jackc/pgx/v4"
)

type Migrator struct {
	dir        fs.FS
	connection *pgx.Conn
}

func NewMigrator(config config.Config) *Migrator {
	return &Migrator{
		dir: os.DirFS(baseDir),
	}
}

func connect(str string) (*pgx.Conn, error) {
	ctx := context.Background()

	return pgx.Connect(ctx, str)
}

func ConnectDSN(dsn string) (*pgx.Conn, error) {
	return connect(dsn)
}

func ConnectURL(url string) (*pgx.Conn, error) {
	return connect(url)
}

func Migrate(connection *pgx.Conn, version string) error {
	currVersion, err := getMigrationDBVersion(connection)
	if err != nil {
		return fmt.Errorf("error while trying to fetch migration database version: %w", err)
	}

	fmt.Println("version of currently installed migration: ", currVersion)

	if IsVerEqualThan(version, currVersion) {
		fmt.Println("Versions are equal, quitting...")
		return nil
	}

	isVersionNewer := IsVerNewerThan(version, currVersion)
	if isVersionNewer {
		MigrateUp(connection, version)
	}

	return nil
}

func MigrateDown(connection *pgx.Conn, version string) {

}

func MigrateUp(connection *pgx.Conn, version string) {

}

func MigrateLatest(connection *pgx.Conn, version string) {

}

const magicCommentStr = "/* UP MIGRATION ABOVE / DOWN MIGRATION BELOW */"

//TODO make this customizable
const baseDir = "./migrations/"

//TODO REFACTOR LOGIC, THIS IS GETTING OUT OF HAND
func assureDirExists(dirpath string) error {
	f, err := os.Open(dirpath)
	if err == nil {
		info, err := f.Stat()
		if err != nil {
			return err
		}
		if !info.IsDir() {
			newErr := errors.New("expected '" + dirpath + "' to be a folder, but instead it is a file")
			return newErr
		}
	}
	f.Close()

	err = os.Mkdir(dirpath, os.ModePerm)

	if err == nil {
		return nil
	}

	pathErr := err.(*os.PathError)

	if pathErr.Op != "mkdir" {
		return errors.New("Unexpected error while trying to create directory: " + err.Error())
	}

	return nil
}

func assureDirExistsRecursive(dirpath string) error {
	dirs := strings.Split(dirpath, "/")

	for i := range dirs {
		err := assureDirExists(strings.Join(dirs[0:i+1], "/"))
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateNewMigration(filepath string) {
	// By default adding the creation time of the migration as it version/id
	// you could just get a way with using it although it would limit the migrations to only be able to be
	// made through the command line interface

	// if we then dont version each file we will never really be able to have a fast
	// migration library as this is a bottleneck in java based libraries as you have to pretty much run the
	// entire migration again when adding new columns etc.

	// although most migration libraries are so basic that they really do not have any complex functionalities,
	// this design sort of works but we could just use database migrations which only save the current state of the table
	// and then try to match the migratable database to have equal model
	// there may be some bottlenecks in such design but with that you could ensure that
	// you would not need to use a cli for the migration

	// using a cli to create new migration files is not really a problem, but id rather
	// have the migrations be more editable by hand.

	// Assures that files cannot be written outside of the directory that they were intended to be saved in.
	if strings.Contains(filepath, "..") {
		err := errors.New("illegal substring '..' found")

		fmt.Println(err.Error())
		return
	}

	dir, ogFileName := path.Split(filepath)

	versionId := PrintTime(time.Now())
	// Gets rid of suffixes, for example abc.efg.hjk -> abc and 352.asdw.g -> 352
	migrationName := strings.Split(ogFileName, ".")[0]
	fileName := versionId + "_" + migrationName + ".sql"

	relPath := path.Join(baseDir, dir)

	err := assureDirExistsRecursive(relPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	newFilePath := path.Join(relPath, fileName)

	f, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println("Could not create file '" + newFilePath + "'")
		return
	}

	f.WriteString("\n\n" + magicCommentStr + "\n\n")

	err = f.Close()
	if err != nil {
		fmt.Println("Could not close file '" + newFilePath + "'")
		return
	}

}

type VersionRow struct {
	Ver string `db:"ver"`
}

func getMigrationDBVersion(connection *pgx.Conn) (string, error) {
	ctx := context.Background()

	row := connection.QueryRow(ctx, "SELECT ver FROM mgr LIMIT 1")

	versionRow := &VersionRow{
		Ver: "0",
	}

	err := row.Scan(versionRow)
	if err != nil {
		return "0", err
	}

	return versionRow.Ver, err
}

func PrintTime(stamp time.Time) string {
	return strconv.FormatInt(stamp.UnixMilli(), 36)
}

// VERY FAST AND IT SHOULD BE USED
func IsVerNewerThan(a string, b string) bool {
	aLen := len(a)
	bLen := len(b)

	if aLen > bLen {
		return true
	} else if aLen < bLen {
		return false
	}

	return a > b
}

func IsVerEqualThan(a string, b string) bool {
	return a == b
}
