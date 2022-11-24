package topology_workloads

import (
	"strings"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (t *TopologyBaseInfraWorkload) PrintWorkload(p filepaths.Path) error {
	if t.Deployment != nil {
		name := addPrefixAndYamlSuffixIfNotExists("dep", t.Deployment.Name)
		err := t.printYaml(&p, name, t.Deployment)
		if err != nil {
			return err
		}
	}
	if t.StatefulSet != nil {
		name := addPrefixAndYamlSuffixIfNotExists("sts", t.StatefulSet.Name)
		err := t.printYaml(&p, name, t.StatefulSet)
		if err != nil {
			return err
		}
	}
	if t.Service != nil {
		name := addPrefixAndYamlSuffixIfNotExists("svc", t.Service.Name)
		err := t.printYaml(&p, name, t.Service)
		if err != nil {
			return err
		}
	}
	if t.ConfigMap != nil {
		name := addPrefixAndYamlSuffixIfNotExists("cm", t.ConfigMap.Name)
		err := t.printYaml(&p, name, t.ConfigMap)
		if err != nil {
			return err
		}
	}
	if t.Ingress != nil {
		name := addPrefixAndYamlSuffixIfNotExists("ing", t.Ingress.Name)
		err := t.printYaml(&p, name, t.Ingress)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TopologyBaseInfraWorkload) printYaml(p *filepaths.Path, name string, workload interface{}) error {
	b, err := yaml.Marshal(workload)
	if err != nil {
		log.Err(err).Msgf("TopologyBaseInfraWorkload: printYaml json.Marshall  %s", name)
		return err
	}
	p.FnOut = name
	err = t.WriteYamlConfig(*p, b)
	if err != nil {
		return err
	}
	return err
}

func (t *TopologyBaseInfraWorkload) WriteYamlConfig(p filepaths.Path, jsonBytes []byte) error {
	err := p.WriteToFileOutPath(jsonBytes)
	if err != nil {
		log.Err(err).Msgf("TopologyBaseInfraWorkload: WriteYamlConfig %s", p.FnOut)
		return err
	}
	return err
}

func addPrefixAndYamlSuffixIfNotExists(prefix, name string) string {
	if !strings.HasPrefix(name, prefix) {
		name = prefix + "-" + name
	}
	if !strings.HasSuffix(name, ".yaml") {
		name = name + ".yaml"
	}
	return name
}
