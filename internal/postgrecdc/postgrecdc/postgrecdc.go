package postgrecdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// https://debezium.io/documentation/reference/stable/connectors/postgresql.html

// After all values ​​are set, you can register.
type PostgreCDC struct {
	// debezium address:port
	// ex) localhost:8083
	DebeziumAddress string

	// database ip or container network name
	// ex) postgres
	DatabaseHost string

	// database port
	DatabasePort string

	// database user
	DatbaseUser string

	// database password
	DatabasePassword string

	// database name
	DatabaseName string

	// database table list
	DatabaseWatchTable []string
}

// p.DatabaseHost + ":" + p.DatabasePort + "." + p.DatabaseName
func (p PostgreCDC) ServerName() string {
	return p.DatabaseHost + "." + p.DatabasePort + "." + p.DatabaseName
}

// = PostgreCDC.ConnectorName()
func (p PostgreCDC) ConnectorName() string {
	return p.ServerName()
}

// // Servername.public.tablename
func (p PostgreCDC) Topics() []string {
	ret := []string{}
	for _, v := range p.DatabaseWatchTable {
		ret = append(ret, p.ServerName()+".public."+v)
	}

	return ret
}

func (p PostgreCDC) IsRegisted() (bool, error) {
	ret, err := p.ConnectorList()
	if err != nil {
		return false, err
	}

	return findKey(p.ConnectorName(), ret), nil
}

func (p PostgreCDC) ConnectorList() ([]string, error) {
	ret := []string{}
	responseGet, err := http.Get("http://" + p.DebeziumAddress + "/connectors/")
	if err != nil {
		return ret, fmt.Errorf("postgrecdc error : %s", err)
	}
	defer responseGet.Body.Close()

	connectionsData, err := ioutil.ReadAll(responseGet.Body)
	if err != nil {
		return ret, fmt.Errorf("postgrecdc get readall error : %s", err)
	}

	if err := json.Unmarshal(connectionsData, &ret); err != nil {
		return ret, fmt.Errorf("postgrecdc json parse error : %s", err)
	}

	return ret, nil
}

func (p PostgreCDC) Regist() error {

	// checking registed
	isRegisted, err := p.IsRegisted()
	if isRegisted {
		return nil
	}
	if err != nil {
		return err
	}

	// create sendPacket
	sendPacket := packet{
		Name:   p.ConnectorName(),
		Config: make(map[string]string),
	}
	sendPacket.Config["connector.class"] = "io.debezium.connector.postgresql.PostgresConnector"
	sendPacket.Config["tasks.max"] = "1"
	sendPacket.Config["database.hostname"] = p.DatabaseHost
	sendPacket.Config["database.port"] = p.DatabasePort
	sendPacket.Config["database.user"] = p.DatbaseUser
	sendPacket.Config["database.password"] = p.DatabasePassword
	sendPacket.Config["database.dbname"] = p.DatabaseName
	sendPacket.Config["database.server.name"] = p.ServerName()

	includeList := ""
	for _, v := range p.DatabaseWatchTable {
		if len(includeList) != 0 {
			includeList += ","
		}
		includeList += "public." + v
	}

	sendPacket.Config["table.include.list"] = includeList

	sendByte, err := json.Marshal(sendPacket)
	if err != nil {
		return fmt.Errorf("kafka cdc json marshal error : %s", err)
	}

	responsePost, err := http.Post("http://"+p.DebeziumAddress+"/connectors/", "application/json", bytes.NewBuffer(sendByte))
	if err != nil {
		return fmt.Errorf("kafka cdc post error : %s", err)
	}
	defer responsePost.Body.Close()
	respBody, err := ioutil.ReadAll(responsePost.Body)
	if err != nil {
		return fmt.Errorf("kafka cdc post read all error : %s", err)
	}

	switch responsePost.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return nil
	default:
		return fmt.Errorf("kafka cdc post error : response (%s) %s", responsePost.Status, string(respBody))
	}
}

func (p PostgreCDC) RemoteConfig() (string, error) {
	responseGet, err := http.Get("http://" + p.DebeziumAddress + "/connectors/" + p.ConnectorName())
	if err != nil {
		return "", fmt.Errorf("kafka cdc error : %s", err)
	}
	defer responseGet.Body.Close()

	connectionsData, err := ioutil.ReadAll(responseGet.Body)
	if err != nil {
		return "", fmt.Errorf("kafka cdc get readall error : %s", err)
	}

	return string(connectionsData), nil
}

func (p PostgreCDC) Unregist() error {
	req, err := http.NewRequest(http.MethodDelete, "http://"+p.DebeziumAddress+"/connectors/"+p.ConnectorName(), nil)
	if err != nil {
		return fmt.Errorf("kafka cdc error : %s", err)
	}
	_, err = http.DefaultClient.Do(req)

	return err
}

func (p PostgreCDC) UnregistOther(connectorName string) error {
	req, err := http.NewRequest(http.MethodDelete, "http://"+p.DebeziumAddress+"/connectors/"+connectorName, nil)
	if err != nil {
		return fmt.Errorf("kafka cdc error : %s", err)
	}
	_, err = http.DefaultClient.Do(req)

	return err
}

type packet struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
}

func findKey(k string, s []string) bool {
	for _, v := range s {
		if k == v {
			return true
		}
	}

	return false
}
