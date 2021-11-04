module github.com/max-gui/consolver

go 1.15

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gohouse/gorose/v2 v2.1.12
	github.com/gomodule/redigo v1.8.5
	github.com/hashicorp/consul/api v1.11.0 // indirect
	github.com/max-gui/consulagent v0.0.0-00010101000000-000000000000
	github.com/max-gui/fileconvagt v0.0.0-00010101000000-000000000000
	github.com/max-gui/logagent v0.0.0-00010101000000-000000000000
	github.com/max-gui/redisagent v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/zsais/go-gin-prometheus v0.1.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/max-gui/logagent => ../logagent

replace github.com/max-gui/consulagent => ../consulagent

replace github.com/max-gui/redisagent => ../redisagent

replace github.com/max-gui/fileconvagt => ../fileconvagt
