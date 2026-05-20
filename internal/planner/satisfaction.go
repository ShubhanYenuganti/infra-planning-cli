package planner

import (
	"strings"

	"github.com/ShubhanYenuganti/infra-planning-cli/internal/discovery"
)

var satisfactionKeywords = map[Domain][]string{
	DomainNetworking: {"aws_vpc", "AWS::EC2::VPC"},
	DomainCompute: {
		"aws_lambda_function", "aws_ecs_service", "aws_instance",
		"AWS::Serverless::Function", "AWS::Lambda::Function",
	},
	DomainDatabase: {
		"aws_db_instance", "aws_rds_cluster", "aws_dynamodb_table",
		"AWS::RDS::DBInstance", "AWS::DynamoDB::Table",
	},
}

// DetectSatisfaction returns (true, evidencePath) when any existing repo
// snippet contains a keyword associated with the request domain.
func DetectSatisfaction(req Request, repo discovery.RepoContext) (bool, string) {
	keywords, ok := satisfactionKeywords[req.PrimaryDomain]
	if !ok {
		return false, ""
	}
	for _, ev := range repo.Evidence {
		for _, kw := range keywords {
			if strings.Contains(ev.Snippet, kw) {
				return true, ev.Path
			}
		}
	}
	return false, ""
}
