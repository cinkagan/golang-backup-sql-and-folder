# Sql And File Backup With Golang
#### 

### **Description**
#### You should change variables' values on main.go.

#### **If you want ignore db name use that**
var ingoreSqlTableList = []string{"information_schema"}

#### **If you want ignore folder name use that**
var ingoreFolderNameList = []string{"www"}

#### **Sql Configs**
Variable Name | Description
|---|---|
hostname | hostname for sql connection
port | port for sql connection
dbname | database name for sql connection
username | username for sql connection
password | password for sql connection


#### **Backup Configs**
Variable Name | Description
|---|---|
sourceFolderDir | name of the folder to be backed up
backupSqlDir |  the name of the folder where the sql backups will be installed
backupFolderDir | the name of the folder where the file backups will be installed

#
### Usage
``` go run main.go ```