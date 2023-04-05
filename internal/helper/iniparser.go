package helper

import (
	"gopkg.in/ini.v1"
)

var localDB DBInfo
var collectPort int
var restPort int
var snapshotPath string
var cornersPort int
var seohoPort int
var tiberoDSN string

func ReadIni(filename string) error {
	config, err := ini.Load(filename)
	if err != nil {
		return err
	}

	localDB = DBInfo{
		Host:     config.Section("TMSDB").Key("Host").MustString("127.0.0.1"),
		User:     config.Section("TMSDB").Key("User").MustString("tms"),
		Password: config.Section("TMSDB").Key("Password").MustString("nvidia"),
		Database: config.Section("TMSDB").Key("Database").MustString("tms"),
		Port:     config.Section("TMSDB").Key("Port").MustInt(3306),
	}

	collectPort = config.Section("TCS").Key("CollectPort").MustInt(10200)
	restPort = config.Section("TCS").Key("RestPort").MustInt(9377)
	snapshotPath = config.Section("TCS").Key("SnapshotPath").MustString("/opt/laonpeople/tcs/share/stream")

	cornersPort = config.Section("Corners").Key("Port").MustInt(4072)

	seohoPort = config.Section("Seoho").Key("Port").MustInt(9485)

	tiberoDSN = config.Section("Tibero").Key("DSN").MustString("tibero6")

	return nil
}

func GetLocalDBInfo() DBInfo {
	return localDB
}
func GetCollectPort() int {
	return collectPort
}
func GetSnapshotPath() string {
	return snapshotPath
}
func GetCornersPort() int {
	return cornersPort
}
func GetRestPort() int {
	return restPort
}
func GetSeohoPort() int {
	return seohoPort
}
func GetTiberoDSN() string {
	return tiberoDSN
}
