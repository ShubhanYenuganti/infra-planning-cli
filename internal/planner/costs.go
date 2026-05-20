package planner

// CostHintsForDomain returns seed monthly cost estimates for the most common
// resources in a domain. All figures are estimates; users must verify against
// current AWS pricing before committing to a budget.
func CostHintsForDomain(domain Domain) []CostHint {
	notes := "estimate; verify against current AWS pricing"
	switch domain {
	case DomainNetworking:
		return []CostHint{
			{ResourceType: "aws_nat_gateway", MonthlyUSDEstimate: 32.40, Notes: notes},
			{ResourceType: "aws_vpc_endpoint (Interface)", MonthlyUSDEstimate: 7.20, Notes: notes},
		}
	case DomainCompute:
		return []CostHint{
			{ResourceType: "aws_lambda_function (1M req/mo)", MonthlyUSDEstimate: 0.20, Notes: notes},
			{ResourceType: "aws_ecs_fargate_task (0.25 vCPU 24/7)", MonthlyUSDEstimate: 8.90, Notes: notes},
		}
	case DomainDatabase:
		return []CostHint{
			{ResourceType: "aws_db_instance (db.t3.micro)", MonthlyUSDEstimate: 12.40, Notes: notes},
			{ResourceType: "aws_dynamodb_table (on-demand baseline)", MonthlyUSDEstimate: 0.00, Notes: notes},
		}
	default:
		return []CostHint{}
	}
}
