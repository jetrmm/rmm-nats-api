package api

/*type Agent struct {
	ID      int    `db:"id"`
	AgentID string `db:"agent_id"`
}*/

type WebConfig struct {
	Key    string `json:"key"`
	RpcUrl string `json:"rpcurl"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DBName string `json:"dbname"`
}
