package commands

import (
	"errors"
	"fmt"

	cmds "github.com/jbenet/go-ipfs/commands"
	core "github.com/jbenet/go-ipfs/core"
	crypto "github.com/jbenet/go-ipfs/crypto"
	nsys "github.com/jbenet/go-ipfs/namesys"
	u "github.com/jbenet/go-ipfs/util"
)

var errNotOnline = errors.New("This command must be run in online mode. Try running 'ipfs daemon' first.")

var publishCmd = &cmds.Command{
	Description: "Publish an object to IPNS",
	Help: `IPNS is a PKI namespace, where names are the hashes of public keys, and
the private key enables publishing new (signed) values.

Examples:

Publish a <ref> to your identity name:

  > ipfs name publish QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy
  published name QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n to QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy

Publish a <ref> to another public key:

  > ipfs name publish QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy
  published name QmbCMUZw6JFeZ7Wp9jkzbye3Fzp2GGcPgC3nmeUjfVF87n to QmatmE9msSfkKxoffpHwNLNKgwZG8eT9Bud6YoPab52vpy

`,

	Arguments: []cmds.Argument{
		cmds.Argument{"name", cmds.ArgString, false, false,
			"The IPNS name to publish to. Defaults to your node's peerID"},
		cmds.Argument{"ipfs-path", cmds.ArgString, true, false,
			"IPFS path of the obejct to be published at <name>"},
	},
	Run: func(res cmds.Response, req cmds.Request) {
		n := req.Context().Node
		args := req.Arguments()

		if n.Network == nil {
			res.SetError(errNotOnline, cmds.ErrNormal)
			return
		}

		if n.Identity == nil {
			res.SetError(errors.New("Identity not loaded!"), cmds.ErrNormal)
			return
		}

		// name := ""
		ref := ""

		switch len(args) {
		case 2:
			// name = args[0]
			ref = args[1].(string)
			res.SetError(errors.New("keychains not yet implemented"), cmds.ErrNormal)
			return
		case 1:
			// name = n.Identity.ID.String()
			ref = args[0].(string)
		}

		// TODO n.Keychain.Get(name).PrivKey
		k := n.Identity.PrivKey()
		publishOutput, err := publish(n, k, ref)

		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}
		res.SetOutput(publishOutput)
	},
	Marshallers: map[cmds.EncodingType]cmds.Marshaller{
		cmds.Text: func(res cmds.Response) ([]byte, error) {
			v := res.Output().(*IpnsEntry)
			s := fmt.Sprintf("Published name %s to %s\n", v.Name, v.Value)
			return []byte(s), nil
		},
	},
	Type: &IpnsEntry{},
}

func publish(n *core.IpfsNode, k crypto.PrivKey, ref string) (*IpnsEntry, error) {
	pub := nsys.NewRoutingPublisher(n.Routing)
	err := pub.Publish(k, ref)
	if err != nil {
		return nil, err
	}

	hash, err := k.GetPublic().Hash()
	if err != nil {
		return nil, err
	}

	return &IpnsEntry{
		Name:  u.Key(hash).String(),
		Value: ref,
	}, nil
}
