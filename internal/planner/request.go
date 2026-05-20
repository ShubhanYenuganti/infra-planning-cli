package planner

import "strings"

type Domain string

const (
	DomainNetworking Domain = "networking"
	DomainCompute    Domain = "compute"
	DomainDatabase   Domain = "database"
	DomainUnknown    Domain = "unknown"
)

type Request struct {
	Raw           string
	PrimaryDomain Domain
}

func ClassifyRequest(raw string) Request {
	text := strings.ToLower(raw)

	for _, token := range []string{"vpc", "subnet", "network", "networking", "route table", "security group", "nat gateway"} {
		if strings.Contains(text, token) {
			return Request{Raw: raw, PrimaryDomain: DomainNetworking}
		}
	}
	for _, token := range []string{"ec2", "ecs", "eks", "lambda", "compute", "container", "worker", "service"} {
		if strings.Contains(text, token) {
			return Request{Raw: raw, PrimaryDomain: DomainCompute}
		}
	}
	for _, token := range []string{"rds", "postgres", "postgresql", "mysql", "dynamodb", "database", "db"} {
		if strings.Contains(text, token) {
			return Request{Raw: raw, PrimaryDomain: DomainDatabase}
		}
	}
	return Request{Raw: raw, PrimaryDomain: DomainUnknown}
}
