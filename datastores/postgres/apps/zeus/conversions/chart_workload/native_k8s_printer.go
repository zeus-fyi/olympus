package chart_workload

import (
	"strings"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/writers"
)

func (nk *NativeK8s) PrintWorkload(p structs.Path) error {
	if nk.Deployment != nil {
		name := addPrefixAndYamlSuffixIfNotExists("dep", nk.Deployment.Name)
		err := nk.printYaml(&p, name, nk.Deployment)
		if err != nil {
			return err
		}
	}
	if nk.StatefulSet != nil {
		name := addPrefixAndYamlSuffixIfNotExists("sts", nk.StatefulSet.Name)
		err := nk.printYaml(&p, name, nk.StatefulSet)
		if err != nil {
			return err
		}
	}
	if nk.Service != nil {
		name := addPrefixAndYamlSuffixIfNotExists("svc", nk.Service.Name)
		err := nk.printYaml(&p, name, nk.Service)
		if err != nil {
			return err
		}
	}
	if nk.ConfigMap != nil {
		name := addPrefixAndYamlSuffixIfNotExists("cm", nk.ConfigMap.Name)
		err := nk.printYaml(&p, name, nk.ConfigMap)
		if err != nil {
			return err
		}
	}
	if nk.Ingress != nil {
		name := addPrefixAndYamlSuffixIfNotExists("ing", nk.Ingress.Name)
		err := nk.printYaml(&p, name, nk.Ingress)
		if err != nil {
			return err
		}
	}
	return nil
}

func (nk *NativeK8s) printYaml(p *structs.Path, name string, workload interface{}) error {
	b, err := yaml.Marshal(workload)
	if err != nil {
		log.Err(err).Msgf("NativeK8s: printYaml json.Marshall  %s", name)
		return err
	}
	p.FnOut = name
	err = nk.WriteYamlConfig(*p, b)
	if err != nil {
		return err
	}
	return err
}

func (nk *NativeK8s) WriteYamlConfig(p structs.Path, jsonBytes []byte) error {
	w := writers.WriterLib{}
	err := w.CreateV2FileOut(p, jsonBytes)
	if err != nil {
		log.Err(err).Msgf("NativeK8s: WriteYamlConfig %s", p.FnOut)
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
