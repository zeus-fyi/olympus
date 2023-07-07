package athena_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	"github.com/zeus-fyi/olympus/pkg/athena"
	athena_workloads "github.com/zeus-fyi/olympus/pkg/athena/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

const bootnodes = "enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303,enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303,enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@65.108.70.101:30303,enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303"

func WorkloadStartup(ctx context.Context, w athena_workloads.WorkloadInfo) {
	log.Info().Interface("w", w).Msg("starting workload")
	switch w.WorkloadType {
	case "p2pCrawler":
		selectedNodes, serr := artemis_mev_models.SelectP2PNodes(ctx, 0)
		if serr != nil {
			log.Fatal().Msg("failed to select p2p nodes")
			misc.DelayedPanic(serr)
		}
		pin := filepaths.Path{
			DirOut: "/data",
			FnOut:  "all-nodes.json",
		}
		b, berr := json.Marshal(selectedNodes)
		if berr != nil {
			log.Fatal().Msg("failed to marshal p2p nodes")
			misc.DelayedPanic(berr)
		}
		werr := pin.WriteToFileOutPath(b)
		if werr != nil {
			log.Fatal().Msg("failed to write p2p nodes")
			misc.DelayedPanic(werr)
		}
		log.Info().Msg("starting address generator")
		go func() {
			err := dynamic_secrets.SaveAddress(ctx, 100000000000000000000, athena.AthenaS3Manager, age)
			if err != nil {
				log.Err(err).Msg("failed to save address")
			}
		}()
		log.Info().Msg("starting p2pCrawler")
		go func() {
			for {
				log.Info().Msg("p2pCrawler loop start")
				cmd := exec.Command("devp2p", "discv4", "crawl", "-timeout", "30m", "--extaddr=127.0.0.1:30303", fmt.Sprintf("--bootnodes=%s", bootnodes), "/data/all-nodes.json")
				err := cmd.Run()
				if err != nil {
					log.Fatal().Msg("failed to start p2pCrawler")
					misc.DelayedPanic(err)
				}
				log.Info().Msg("p2pCrawler main loop complete")
				p := filepaths.Path{
					DirIn:  "/data",
					FnIn:   "all-nodes.json",
					DirOut: "/data",
					FnOut:  "mainnet-nodes.json",
				}
				b = p.ReadFileInPath()
				err = artemis_mev_models.InsertP2PNodes(ctx, artemis_autogen_bases.EthP2PNodes{
					ID:                0,
					ProtocolNetworkID: 0,
					Nodes:             string(b),
				})
				if err != nil {
					log.Fatal().Msg("failed to insert p2pCrawler node results")
					misc.DelayedPanic(err)
				}
				log.Info().Msg("p2pCrawler mainnet filter start")
				cmd = exec.Command("devp2p", "nodeset", "filter", p.FileInPath(), "-eth-network", "mainnet", "-limit", "2000")
				outFile, err := os.Create(p.FileOutPath())
				if err != nil {
					log.Fatal().Msg("failed to filter p2pCrawler mainnet node results")
					misc.DelayedPanic(err)
				}
				log.Info().Msg("p2pCrawler mainnet filter start")
				cmd.Stdout = outFile
				err = cmd.Run()
				if err != nil {
					log.Fatal().Msg("failed to filter p2pCrawler nodes")
					misc.DelayedPanic(err)
				}
				outFile.Close()
				log.Info().Msg("p2pCrawler mainnet filter done")
				p = filepaths.Path{
					DirIn: "/data",
					FnIn:  "mainnet-nodes.json",
				}
				b = p.ReadFileInPath()
				if b == nil {
					log.Fatal().Msg("failed to read p2pCrawler mainnet node results")
					misc.DelayedPanic(err)
				}
				var nodes artemis_mev_models.P2PNodes
				err = json.Unmarshal(b, &nodes)
				if err != nil {
					log.Fatal().Msg("failed to Unmarshal p2pCrawler nodes")
					misc.DelayedPanic(err)
				}
				log.Info().Int("nodeLen", len(nodes)).Msg("p2pCrawler nodes")
				err = artemis_mev_models.InsertP2PNodes(ctx, artemis_autogen_bases.EthP2PNodes{
					ID:                hestia_req_types.EthereumMainnetProtocolNetworkID,
					ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
					Nodes:             string(b),
				})
				if err != nil {
					log.Fatal().Msg("failed to insert p2pCrawler nodes")
					misc.DelayedPanic(err)
				}
				log.Info().Msg("p2pCrawler loop end")

				//p = filepaths.Path{
				//	DirOut: "/data",
				//	FnOut:  "goerli-nodes.json",
				//}
				//cmd = exec.Command("devp2p", "nodeset", "filter", p.FileInPath(), "-eth-network", "goerli")
				//outFile, err = os.Create(p.FileOutPath())
				//if err != nil {
				//	log.Fatal().Msg("failed to filter p2pCrawler goerli node results")
				//	misc.DelayedPanic(err)
				//}
				//cmd.Stdout = outFile
				//err = cmd.Run()
				//if err != nil {
				//	log.Fatal().Msg("failed to filter p2pCrawler goerli nodes")
				//	misc.DelayedPanic(err)
				//}
				//outFile.Close()
				//p = filepaths.Path{
				//	DirIn: "/data",
				//	FnIn:  "goerli-nodes.json",
				//}
				//b = p.ReadFileInPath()
				//var goerliNodes artemis_validator_service_groups_models.P2PNodes
				//err = json.Unmarshal(b, &goerliNodes)
				//if err != nil {
				//	log.Fatal().Msg("failed to Unmarshal goerli p2pCrawler nodes")
				//	misc.DelayedPanic(err)
				//}
				//err = artemis_validator_service_groups_models.InsertP2PNodes(ctx, artemis_autogen_bases.EthP2PNodes{
				//	ID:                hestia_req_types.EthereumGoerliProtocolNetworkID,
				//	ProtocolNetworkID: hestia_req_types.EthereumGoerliProtocolNetworkID,
				//	Nodes:             string(b),
				//})
				//if err != nil {
				//	log.Fatal().Msg("failed to insert goerli p2pCrawler nodes")
				//	misc.DelayedPanic(err)
				//}
			}
		}()
	}
}
