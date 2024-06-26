package api

import (
	"encoding/json"
	"reflect"
	"runtime"
	"time"

	rmm "github.com/jetrmm/rmm-shared"
	// _ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
)

func Svc(logger *logrus.Logger, cfg string) {
	logger.Debugln("Starting Svc()")
	db, r, err := GetConfig(cfg)
	if err != nil {
		logger.Fatalln(err)
	}

	opts := setupNatsOptions(r.Key)
	nc, err := nats.Connect(r.RpcUrl, opts...)
	if err != nil {
		logger.Fatalln(err)
	}

	nc.Subscribe("*", func(msg *nats.Msg) {
		var mh codec.MsgpackHandle
		mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
		mh.RawToString = true
		dec := codec.NewDecoderBytes(msg.Data, &mh)

		switch msg.Reply {
		case "agent-hello":
			go func() {
				var p rmm.AgentHeaderNats
				if err := dec.Decode(&p); err == nil {
					loc, _ := time.LoadLocation("UTC")
					now := time.Now().In(loc)
					logger.Debugln("Hello", p, now)
					stmt := `UPDATE agents SET last_seen=$1, version=$2 WHERE agents.agent_id=$3;`

					_, err = db.Exec(stmt, now, p.Version, p.AgentId)
					if err != nil {
						logger.Errorln(err)
					}
				}
			}()

		case "agent-publicip":
			go func() {
				var p rmm.PublicIPNats
				if err := dec.Decode(&p); err == nil {
					logger.Debugln("Public IP", p)
					stmt := `UPDATE agents SET public_ip=$1 WHERE agents.agent_id=$2;`
					_, err = db.Exec(stmt, p.PublicIP, p.AgentId)
					if err != nil {
						logger.Errorln(err)
					}
				}
			}()

		case "agent-agentinfo":
			go func() {
				var r rmm.AgentInfoNats
				if err := dec.Decode(&r); err == nil {
					stmt := `
						UPDATE agents
						SET hostname=$1, operating_system=$2,
						plat=$3, total_ram=$4, boot_time=$5, reboot_pending=$6, logged_in_username=$7
						WHERE agents.agent_id=$8;`

					logger.Debugln("Info", r)
					_, err = db.Exec(stmt, r.Hostname, r.OS, r.Platform, r.TotalRAM, r.BootTime, r.RebootPending, r.Username, r.AgentId)
					if err != nil {
						logger.Errorln(err)
					}

					if r.Username != "None" {
						stmt = `UPDATE agents SET last_logged_in_user=$1 WHERE agents.agent_id=$2;`
						logger.Debugln("Updating last logged in user:", r.Username)
						_, err = db.Exec(stmt, r.Username, r.AgentId)
						if err != nil {
							logger.Errorln(err)
						}
					}
				}
			}()

		case "agent-disks":
			go func() {
				var r rmm.StorageNats
				if err := dec.Decode(&r); err == nil {
					logger.Debugln("Drives", r)
					b, err := json.Marshal(r.Drives)
					if err != nil {
						logger.Errorln(err)
						return
					}
					stmt := `UPDATE agents SET disks=$1 WHERE agents.agent_id=$2;`

					_, err = db.Exec(stmt, b, r.AgentId)
					if err != nil {
						logger.Errorln(err)
					}
				}
			}()

		case "agent-winsvc":
			go func() {
				var r rmm.WinSvcNats
				if err := dec.Decode(&r); err == nil {
					logger.Debugln("WinSvc", r)
					b, err := json.Marshal(r.WinSvcs)
					if err != nil {
						logger.Errorln(err)
						return
					}

					stmt := `UPDATE agents SET services=$1 WHERE agents.agent_id=$2;`

					_, err = db.Exec(stmt, b, r.AgentId)
					if err != nil {
						logger.Errorln(err)
					}
				}
			}()

		case "agent-sysinfo": // was "agent-wmi"
			go func() {
				var r rmm.WinWMINats
				if err := dec.Decode(&r); err == nil {
					logger.Debugln("WMI", r)
					b, err := json.Marshal(r.WMI)
					if err != nil {
						logger.Errorln(err)
						return
					}
					stmt := `UPDATE agents SET wmi_detail=$1 WHERE agents.agent_id=$2;`

					_, err = db.Exec(stmt, b, r.AgentId)
					if err != nil {
						logger.Errorln(err)
					}
				}
			}()
		}
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		logger.Fatalln(err)
	}
	runtime.Goexit()
}
