package solusvm

import (
	log "github.com/Sirupsen/logrus"
  "fmt"

	"github.com/megamsys/megdcui/automation"
	"github.com/megamsys/megdcui/migration"
	"github.com/megamsys/libgo/action"
	"gopkg.in/yaml.v2"
)

func init() {
	migration.Register("solusvm", solusvmManager{})
}

type solusvmManager struct{}

type VirtualServer struct {
	Vserverid   	*string `json:"vserverid"`
	Ctid_xid    	*string `json:"ctid-xid"`
	Clientid      *string `json:"clientid"`
	Ipaddress     *string `json:"ipaddress"`
	Hostname      *string `json:"hostname"`
	Template      *string `json:"template"`
	Hdd    				*string `json:"hdd"`
	Memory     		*string `json:"memory"`
	Swap_burst 		*string `json:"swap-burst"`
	Type    			*string `json:"type"`
	Mac      			*string `json:"mac"`
}

type VServers struct {
	Status     *string `json:"status"`
	Statusmsg  *string `json:"statusmsg"`
	VirtualServers *[]VirtualServer `json:"virtualservers"`
	Org_id     string
}

func (b *VirtualServer) String() string {
	if d, err := yaml.Marshal(b); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func (m solusvmManager) MigratablePrepare(h *automation.HostInfo)  error {

	actions := []*action.Action{
		&VertifyMigratableCredentials,
		//&VerfiyMigrationComplete,
	}
	pipeline := action.NewPipeline(actions...)

	args := runActionsArgs{
    h:        h,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute status pipeline for github %s - %s", h.SolusMaster, err)
		return err
	}

	return nil

}
func (m solusvmManager) MigrateHost(h *automation.HostInfo) (*automation.Result, error) {

	actions := []*action.Action{
  // &ListClientsInMigratable,
	 &OnboardClientsInVertice,
	// &ListVMsinMigratable,
	 &TagMigratableInVertice,
	}

	pipeline := action.NewPipeline(actions...)

	args := runActionsArgs{
    h: h,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute status pipeline for node %s - %s", h.SolusNode, err)
		r := &automation.Result{
			Status: "error",
		 	Statusmsg: fmt.Sprintf("%s",err),
		 	VirtualServers: "",
		}
		return r,err
	}
	r := &automation.Result{
		Status: "success",
	 	Statusmsg: "virtual machines migrated successfully",
	 	VirtualServers: "",
	}
	return r, nil

}
