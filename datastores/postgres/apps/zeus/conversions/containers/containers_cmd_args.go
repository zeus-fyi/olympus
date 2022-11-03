package containers

import (
	"strings"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	v1 "k8s.io/api/core/v1"
)

func ConvertContainerCmdArgsToContainerDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	cmd := cs.Command
	args := cs.Args
	ca := autogen_bases.ContainerCommandArgs{
		CommandArgsID: 0,
		CommandValues: strings.Join(cmd, ","),
		ArgsValues:    strings.Join(args, ","),
	}
	dbContainer.CmdArgs = ca
	return dbContainer
}
