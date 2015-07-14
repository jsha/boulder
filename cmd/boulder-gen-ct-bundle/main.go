// Copyright 2015 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/letsencrypt/boulder/cmd"
	"github.com/letsencrypt/boulder/core"
)

func prettyChain(chain []*x509.Certificate) string {
	var res string
	for i := range chain {
		res += chain[i].Subject.CommonName
		if i != len(chain)-1 {
			res += " -> "
		}
	}

	return res
}

func downloadCTRootList(client *http.Client, uri string) ([]*x509.Certificate, error) {
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = fmt.Sprintf("%s%s", "https://", uri)
	}
	uri = fmt.Sprintf("%s/ct/v1/get-roots", uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("Creating request failed, %s", err)
	}

	fmt.Printf("\tDownloading roots from %s\n", uri)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed, %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body, %s", err)
	}

	var roots struct {
		Certificates []string `json:"certificates"`
	}
	err = json.Unmarshal(body, &roots)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal root list, %s", err)
	}
	fmt.Printf("\tDownloaded %d certificates\n", len(roots.Certificates))
	var rootCerts []*x509.Certificate
	for _, root := range roots.Certificates {
		certBytes, err := base64.StdEncoding.DecodeString(root)
		if err != nil {
			// log but continue
			continue
		}
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			// log but continue
			continue
		}
		rootCerts = append(rootCerts, cert)
	}

	return rootCerts, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "gen-ct-bundle"
	app.Version = fmt.Sprintf("0.1.0 [%s]", core.GetBuildID())

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "config.json",
			EnvVar: "BOULDER_CONFIG",
			Usage:  "Path to Boulder JSON configuration file",
		},
		cli.StringSliceFlag{
			Name:  "cert",
			Value: &cli.StringSlice{},
			Usage: "Paths to PEM certificates in order from issuer -> root",
		},
		cli.StringFlag{
			Name:  "out",
			Usage: "Path to write out the DER certificate bundle",
		},
	}

	app.Action = func(c *cli.Context) {
		bundleFilename := c.GlobalString("out")
		if bundleFilename == "" {
			fmt.Fprintf(os.Stderr, "-out flag is required\n")
			os.Exit(1)
		}
		fmt.Println(bundleFilename)

		configFileName := c.GlobalString("config")
		configJSON, err := ioutil.ReadFile(configFileName)
		cmd.FailOnError(err, "Unable to read config file")

		var config cmd.Config
		err = json.Unmarshal(configJSON, &config)
		cmd.FailOnError(err, "Failed to read configuration")

		if config.Publisher.CT == nil {
			fmt.Fprintf(os.Stderr, "Publisher CT configuration required to assemble correct root pool\n")
			os.Exit(1)
		}

		certFilenames := c.GlobalStringSlice("cert")
		if len(certFilenames) == 0 {
			//
		}
		var chain []*x509.Certificate
		for _, certFilename := range certFilenames {
			cert, err := core.LoadCert(certFilename)
			cmd.FailOnError(err, fmt.Sprintf("Failed to load certificate (%s)", certFilename))
			chain = append(chain, cert)
		}

		roots := x509.NewCertPool()
		roots.AddCert(chain[len(chain)-1])
		opts := x509.VerifyOptions{Roots: roots}

		if len(chain) > 2 {
			inters := x509.NewCertPool()
			for _, inter := range chain[1 : len(chain)-1] {
				inters.AddCert(inter)
			}
			opts.Intermediates = inters
		}
		_, err = chain[0].Verify(opts)
		cmd.FailOnError(err, "Failed to load chain")

		client := http.Client{}
		maxLen := 0
		for _, ctLog := range config.Publisher.CT.Logs {
			fmt.Printf("# %s\n", ctLog.URI)
			logRoots, err := downloadCTRootList(&client, ctLog.URI)
			cmd.FailOnError(err, "Failed to retrieve root certificates")
			ctPool := x509.NewCertPool()
			for _, root := range logRoots {
				ctPool.AddCert(root)
			}
			ctOpts := x509.VerifyOptions{Roots: ctPool}
			var lastValidChain []*x509.Certificate
			fmt.Println("\tTesting chain validity with downloaded log root pool")
			for i := range chain {
				if len(chain)-i > 1 {
					ctOpts.Intermediates = x509.NewCertPool()
					for _, inter := range chain[1 : len(chain)-i] {
						ctOpts.Intermediates.AddCert(inter)
					}
				} else {
					ctOpts.Intermediates = nil
				}
				fmt.Printf("\t\t%s -> constructed root pool", prettyChain(chain[:len(chain)-i]))
				if _, err := chain[0].Verify(ctOpts); err != nil {
					fmt.Println(": Invalid!")
					break
				}
				fmt.Println(": Valid!")
				lastValidChain = chain[:len(chain)-i]
			}
			if len(lastValidChain) == 0 {
				fmt.Println("\n\t!! Couldn't construct any valid chains, this may mean you haven't   !!")
				fmt.Println("\t!! provided the full chain or that this CT log doesn't contain a    !!")
				fmt.Println("\t!! root certificates that chain those provided. In the case of the  !!")
				fmt.Println("\t!! latter you should remove this log from your configuration since  !!")
				fmt.Println("\t!! your submissions will fail and be discarded.                     !!")
				continue
			}
			fmt.Printf("\n\tBundle size for %s: %d\n", ctLog.URI, len(lastValidChain))
			if len(lastValidChain) > maxLen {
				maxLen = len(lastValidChain)
			}
		}

		if maxLen == 0 {
			fmt.Println("\n!! Couldn't find any valid chains for configured logs, this may     !!")
			fmt.Println("!! mean you haven't provided the full chain or that this CT log     !!")
			fmt.Println("!! doesn't contain a root certificates that chain those provided.   !!")
			fmt.Println("!! The bundle will still be written out but you should only use     !!")
			fmt.Println("!! this if you really know what you are doing!                      !!")
			maxLen = len(chain)
		}
		fmt.Printf("\n# Shared bundle size: %d certificates, %s\n", maxLen, prettyChain(chain[:maxLen]))

		// Write bundle out
		f, err := os.Create(bundleFilename)
		cmd.FailOnError(err, fmt.Sprintf("Failed to create submission bundle (%s)", bundleFilename))
		defer f.Close()

		for _, cert := range chain[:maxLen] {
			f.Write(cert.Raw)
		}
		fmt.Printf("# CT submission bundle has been written to %s\n", bundleFilename)
	}

	err := app.Run(os.Args)
	cmd.FailOnError(err, "Failed to run application")
}
