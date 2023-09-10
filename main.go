package main

import (
	"errors"
	"log"
	"net/url"
	"os/exec"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/google/go-querystring/query"
)

const (
	PrefixGKEContext = "gke"
	GKEContextSplit  = "_"
	GKEContextLength = 4
	GKETypeIdx       = 0
	GKEProjectIdx    = 1
	GKELocationIdx   = 2
	GKEClusterIdx    = 3

	// CommandOpenURL is the command to open URL, it's different between OS, currently only support Linux
	CommandOpenURL = "xdg-open"
)

var (
	errNotGKEContext = errors.New("context is not from GKE")
)

type GKELogQueryParams struct {
	ProjectID     string `url:"resource.labels.project_id"`
	Location      string `url:"resource.labels.location"`
	ClusterName   string `url:"resource.labels.cluster_name"`
	GkeContext    string `url:"-"                              arg:"-g,required"`
	ResourceType  string `url:"resource.type"                  arg:"-t,required"`
	NamespaceName string `url:"resource.labels.namespace_name" arg:"-n,required"`
	ContainerName string `url:"resource.labels.container_name" arg:"-c,required" `
}

func (qp *GKELogQueryParams) getGKELogQuery() (string, error) {
	values, err := query.Values(qp)
	if err != nil {
		return "", err
	}

	gkeQuery := strings.ReplaceAll(values.Encode(), "&", "\n")
	gkeQueryURL := &url.URL{
		Scheme: "https",
		Host:   "console.cloud.google.com",
		Path:   "/logs/query;query=" + gkeQuery,
	}

	urlQuery := url.Values{}
	urlQuery.Add("project", qp.ProjectID)
	gkeQueryURL.RawQuery = urlQuery.Encode()

	return gkeQueryURL.String(), nil
}

func NewParamsParser() *ParamsParser {
	return &ParamsParser{
		parseArgs: func(q *GKELogQueryParams) {
			arg.MustParse(q)
		},
	}
}

type ParamsParser struct {
	parseArgs func(*GKELogQueryParams)
}

func (p *ParamsParser) getQueryParams() (*GKELogQueryParams, error) {
	args := GKELogQueryParams{}
	p.parseArgs(&args)

	gkeContextParts := strings.Split(args.GkeContext, GKEContextSplit)
	if len(gkeContextParts) != GKEContextLength {
		return nil, errNotGKEContext
	}

	if gkeContextParts[GKETypeIdx] != PrefixGKEContext {
		return nil, errNotGKEContext
	}

	args.ProjectID = gkeContextParts[GKEProjectIdx]
	args.Location = gkeContextParts[GKELocationIdx]
	args.ClusterName = gkeContextParts[GKEClusterIdx]

	return &args, nil
}

func openURL(url string) error {
	return exec.Command(CommandOpenURL, url).Run()
}

func main() {
	args, err := NewParamsParser().getQueryParams()
	if err != nil {
		log.Panic(err)
	}

	gkeLogQuery, err := args.getGKELogQuery()
	if err != nil {
		log.Panic(err)
	}

	if err := openURL(gkeLogQuery); err != nil {
		log.Panic(err)
	}
}
