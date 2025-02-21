package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cmosh/interlock/config"
	"github.com/cmosh/interlock/server"
	"github.com/cmosh/interlock/version"
	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/tlsconfig"
	"github.com/docker/libkv"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
)

const (
	defaultConfig = `ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
EnableMetrics = true
`
	kvConfigKey = "interlock/v1/config"
)

var cmdRun = cli.Command{
	Name:   "run",
	Usage:  "run interlock",
	Action: runAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config file",
			Value: "",
		},
		cli.StringFlag{
			Name:  "discovery, k",
			Usage: "discovery address",
			Value: "",
		},
		cli.StringFlag{
			Name:  "discovery-tls-ca-cert",
			Usage: "discovery tls ca certificate",
			Value: "",
		},
		cli.StringFlag{
			Name:  "discovery-tls-cert",
			Usage: "discovery tls certificate",
			Value: "",
		},
		cli.StringFlag{
			Name:  "discovery-tls-key",
			Usage: "discovery tls key",
			Value: "",
		},
	},
}

func init() {
	consul.Register()
	etcd.Register()
}

func getKVStore(addr string, options *kvstore.Config) (kvstore.Store, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err

	}

	kvType := strings.ToLower(u.Scheme)
	kvHost := u.Host
	var backend kvstore.Backend

	switch kvType {
	case "consul":
		backend = kvstore.CONSUL
	case "etcd":
		backend = kvstore.ETCD
	}

	kv, err := libkv.NewStore(
		backend,
		[]string{kvHost},
		options,
	)

	if err != nil {
		return nil, err
	}

	return kv, nil
}

func runAction(c *cli.Context) {
	log.Infof("interlock %s", version.FullVersion())

	var data string

	if envCfg := os.Getenv("INTERLOCK_CONFIG"); envCfg != "" {
		log.Debug("loading config from environment")

		data = envCfg
	}

	if dURL := c.String("discovery"); dURL != "" {
		log.Debugf("loading config from key value store: addr=%s", dURL)

		// init kv
		kvOpts := &kvstore.Config{
			ConnectionTimeout: time.Second * 10,
		}

		dTLSCACert := c.String("discovery-tls-ca-cert")
		dTLSCert := c.String("discovery-tls-cert")
		dTLSKey := c.String("discovery-tls-key")

		if dTLSCACert != "" && dTLSCert != "" && dTLSKey != "" {
			tlsConfig, err := tlsconfig.Client(tlsconfig.Options{
				CAFile:   dTLSCACert,
				CertFile: dTLSCert,
				KeyFile:  dTLSKey,
			})
			if err != nil {
				log.Fatal(err)
			}

			log.Debug("configuring TLS for KV")
			kvOpts.TLS = tlsConfig
		}

		kv, err := getKVStore(dURL, kvOpts)
		if err != nil {
			log.Fatal(err)
		}

		// get config from kv
		exists, err := kv.Exists(kvConfigKey)
		if err != nil {
			log.Fatal(err)
		}

		if !exists {
			log.Warnf("unable to find config in key %s; using default config", kvConfigKey)
			data = defaultConfig
		} else {
			kvPair, err := kv.Get(kvConfigKey)
			if err != nil {
				log.Fatalf("error getting configuration from key value store: %s", err)
			}

			data = string(kvPair.Value)

			if data == "" {
				data = defaultConfig
			}
		}
	}

	if configPath := c.String("config"); configPath != "" && data == "" {
		log.Debugf("loading config from: file=%s", configPath)

		d, err := ioutil.ReadFile(configPath)
		switch {
		case os.IsNotExist(err):
			log.Errorf("Missing Interlock configuration: file=%s", configPath)
			log.Error("Use the run --config option to set a custom location for the configuration file")
			log.Error("Examples of an Interlock configuration file: url=https://github.com/cmosh/interlock/tree/master/docs/examples")
			log.Fatalf("config not found: file=%s", configPath)
		case err == nil:
			data = string(d)
		default:
			log.Fatal(err)
		}
	}

	if data == "" {
		log.Error("Examples of Interlock configuration: url=https://github.com/cmosh/interlock/blob/master/docs/configuration.md")
		log.Fatal("You must specify a config from a file, environment variable, or key value store")
	}

	config, err := config.ParseConfig(data)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
