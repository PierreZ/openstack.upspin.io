// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The upspin-setupstorage-aws command is an external upspin subcommand that
// executes the second step in establishing an upspinserver for AWS.
// Run upspin setupstorage-aws -help for more information.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"upspin.io/subcmd"
)

type state struct {
	*subcmd.State
}

const help = `
Setupstorage-aws is the second step in establishing an upspinserver.
It sets up OpenStack storage for your Upspin installation. You may skip this step
if you wish to store Upspin data on your server's local disk.
The first step is 'setupdomain' and the final step is 'setupserver'.

Setupstorage-aws creates an Amazon Swift container. It then updates the server
configuration files in $where/$domain/ to use the specified container and region.

Before running this command, you should ensure you have an OpenStack region and that
the swift CLI command line tool works for you.

If something goes wrong during the setup process, you can run the same command
with the -clean flag. It will attempt to remove any entities previously created
with the same options provided.
`

func main() {
	const name = "setupstorage-openstack"

	log.SetFlags(0)
	log.SetPrefix("upspin setupstorage-openstack: ")

	s := &state{
		State: subcmd.NewState(name),
	}

	var (
		where  = flag.String("where", filepath.Join(os.Getenv("HOME"), "upspin", "deploy"), "`directory` to store private configuration files")
		domain = flag.String("domain", "", "domain `name` for this Upspin installation")
		region = flag.String("region", "", "region for the Swift container")
		//clean  = flag.Bool("clean", false, "deletes all artifacts that would be created using this command")
	)

	s.ParseFlags(flag.CommandLine, os.Args[1:], help,
		"setupstorage-openstack -domain=<name> [-region=<region>] [-clean] <container_name>")
	if flag.NArg() != 1 {
		s.Exitf("a single container name must be provided")
	}
	if len(*domain) == 0 {
		s.Exitf("the -domain flag must be provided")
	}

	containerName := flag.Arg(0)
	// if *clean {
	// 	s.clean(*roleName, bucketName, *region)
	// 	s.ExitNow()
	// }

	cfgPath := filepath.Join(*where, *domain)
	cfg := s.ReadServerConfig(cfgPath)

	cfg.StoreConfig = []string{
		"backend=OpenStack",
		"openstackContainer=" + containerName,
		"openstackRegion=" + *region,
	}
	s.WriteServerConfig(cfgPath, cfg)

	fmt.Fprintf(os.Stderr, "You should now deploy the upspinserver binary and run 'upspin setupserver'.\n")
	s.ExitNow()
}
