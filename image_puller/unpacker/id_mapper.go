package unpacker // import "code.cloudfoundry.org/grootfs/image_puller/unpacker"

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"

	"code.cloudfoundry.org/commandrunner"

	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/lager"
)

type CommandIDMapper struct {
	cmdRunner commandrunner.CommandRunner
}

func NewIDMapper(cmdRunner commandrunner.CommandRunner) *CommandIDMapper {
	return &CommandIDMapper{
		cmdRunner: cmdRunner,
	}
}

func (im *CommandIDMapper) MapUIDs(logger lager.Logger, pid int, mappings []groot.IDMappingSpec) error {
	logger = logger.Session("mapUID")
	return im.execute(logger, "newuidmap", pid, mappings)
}

func (im *CommandIDMapper) MapGIDs(logger lager.Logger, pid int, mappings []groot.IDMappingSpec) error {
	logger = logger.Session("mapGID")
	return im.execute(logger, "newgidmap", pid, mappings)
}

func (im *CommandIDMapper) execute(logger lager.Logger, command string, pid int, mappings []groot.IDMappingSpec) error {
	logger.Info("start")
	defer logger.Info("end")

	mappingArgs := append([]string{strconv.Itoa(pid)}, im.idMappingsToArgs(mappings)...)
	mapCmd := exec.Command(command, mappingArgs...)

	buffer := bytes.NewBuffer([]byte{})
	mapCmd.Stdout = buffer
	mapCmd.Stderr = buffer

	logger.Debug("starting-id-map", lager.Data{"path": mapCmd.Path, "args": mapCmd.Args})
	if err := im.cmdRunner.Run(mapCmd); err != nil {
		return fmt.Errorf("%s %s: %s", command, err, buffer.String())
	}

	return nil
}

func (im *CommandIDMapper) idMappingsToArgs(idMappings []groot.IDMappingSpec) (args []string) {
	for _, idMapping := range idMappings {
		args = append(args, im.idMappingToArgs(idMapping)...)
	}

	return args
}

func (im *CommandIDMapper) idMappingToArgs(idMapping groot.IDMappingSpec) []string {
	return []string{
		strconv.Itoa(idMapping.NamespaceID),
		strconv.Itoa(idMapping.HostID),
		strconv.Itoa(idMapping.Size),
	}
}
