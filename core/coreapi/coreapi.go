package coreapi

import (
	"context"

	core "github.com/ipfs/go-ipfs/core"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
	ipfspath "github.com/ipfs/go-ipfs/path"

	ipld "gx/ipfs/QmRSU5EqqWVZSNdbU51yXmVoF1uNw3JgTNB6RaiL7DZM16/go-ipld-node"
	cid "gx/ipfs/QmcTcsTvfaeEBRFo1TkFgT8sRmgi1n1LTZpecfVP8fzpGD/go-cid"
)

type CoreAPI struct {
	node *core.IpfsNode
}

func NewCoreAPI(n *core.IpfsNode) coreiface.CoreAPI {
	api := &CoreAPI{n}
	return api
}

func (api *CoreAPI) Unixfs() coreiface.UnixfsAPI {
	return (*UnixfsAPI)(api)
}

func (api *CoreAPI) ResolveNode(ctx context.Context, p coreiface.Path) (ipld.Node, error) {
	p, err := api.ResolvePath(ctx, p)
	if err != nil {
		return nil, err
	}

	node, err := api.node.DAG.Get(ctx, p.Cid())
	if err == core.ErrNoNamesys {
		return nil, coreiface.ErrOffline
	} else if err != nil {
		return nil, err
	}
	return node, nil
}

func (api *CoreAPI) ResolvePath(ctx context.Context, p coreiface.Path) (coreiface.Path, error) {
	if p.Resolved() {
		return p, nil
	}

	c, err := core.ResolveToCid(ctx, api.node, ipfspath.FromString(p.String()))
	if err != nil {
		return nil, err
	}
	return NewResolvedPath(p.String(), c), nil
}

type path struct {
	path string
	cid  *cid.Cid
}

func ParsePath(p string) (coreiface.Path, error) {
	pp, err := ipfspath.ParsePath(p)
	if err != nil {
		return nil, err
	}
	return &path{path: pp.String()}, nil
}
func NewResolvedPath(p string, c *cid.Cid) coreiface.Path { return &path{path: p, cid: c} }
func NewPathFromCid(c *cid.Cid) coreiface.Path            { return &path{path: "/ipfs/" + c.String(), cid: c} }
func (p *path) String() string                            { return p.path }
func (p *path) Cid() *cid.Cid                             { return p.cid }
func (p *path) Resolved() bool                            { return p.cid != nil }
