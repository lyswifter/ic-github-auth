package cmd

import (
	"github.com/lyswifter/ic-auth/db"
	"github.com/lyswifter/ic-auth/util"
	"github.com/urfave/cli"
)

var InitCmd = cli.Command{
	Name:        "init",
	Description: "Initial Internet Computer authorization repo",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "repo",
			Value: "~/.icauth",
			Usage: "Specify the location of database",
		},
	},
	Action: func(cctx *cli.Context) error {
		rdb, idb, err := db.DataStores(cctx.String("repo"))
		if err != nil {
			return err
		}

		util.Infof("rdb: %+v idb: %+v", rdb, idb)
		return nil
	},
}
