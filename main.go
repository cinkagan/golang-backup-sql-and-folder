package main

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cinkagan/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ziutek/mymysql/godrv"
)

//If you want ignore db name use that
var ingoreSqlTableList = []string{"information_schema"}

// If you want ignore folder name use that
var ingoreFolderNameList = []string{"www"}

// Sql Config
var hostname string = ""
var port string = ""
var dbname string = ""
var username string = ""
var password string = ""

// Backup Config
var sourceFolderDir string = "www"
var backupSqlDir string = "backup/sql"
var backupFolderDir string = "backup/folder"

// Array contains function
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// File zip function
func zipit(source, target string) error {
	zipfile, err := os.Create(backupFolderDir + "/" + target + ".zip")
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

//Backup function of all db
func sqlBackup() {
	godrv.Register("SET NAMES utf8")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname))
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return
	}

	defer db.Close()

	rows, err := db.Query("Show databases;")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		if !contains(ingoreSqlTableList, name) {
			log.Println("*******************")
			log.Println(name + " is backuping")
			dumpFilenameFormat := fmt.Sprintf("%s-20060102150405", name)
			dumper, err := mysqldump.Register(db, backupSqlDir, dumpFilenameFormat)
			if err != nil {
				fmt.Println("Error registering databse:", err)
				return
			}

			resultFilename, err := dumper.Dump()
			if err != nil {
				fmt.Println("Error dumping:", err)
			}
			log.Println(name + " is backuped")
			log.Printf("File is saved to %s \n", resultFilename)
			log.Println("*******************")

			dumper.Close()
		}

	}
}

//Folder backup function
func folderBackup() {
	err := filepath.Walk(sourceFolderDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}
		if contains(ingoreFolderNameList, info.Name()) {
			return nil
		}
		log.Println("*******************")
		log.Println(info.Name() + " is backuping")
		zipFileName := fmt.Sprintf("%s-%s", info.Name(), time.Now().Format("20060102150405"))
		err = zipit(path, zipFileName)
		if err != nil {
			return err
		}
		resultFileName := backupFolderDir + "/" + zipFileName + ".zip"
		log.Println(info.Name() + " is backuped")
		log.Printf("File is saved to %s \n", resultFileName)
		log.Println("*******************")
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	sqlBackup()
	folderBackup()
}
