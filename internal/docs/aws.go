package docs

type DocLink struct {
	Title string
	URL   string
}

func AWSDocsForDomain(domain string) []DocLink {
	switch domain {
	case "networking":
		return []DocLink{
			{Title: "Amazon VPC documentation", URL: "https://docs.aws.amazon.com/vpc/"},
			{Title: "VPC security groups", URL: "https://docs.aws.amazon.com/vpc/latest/userguide/vpc-security-groups.html"},
			{Title: "Route tables", URL: "https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html"},
		}
	case "compute":
		return []DocLink{
			{Title: "Amazon ECS documentation", URL: "https://docs.aws.amazon.com/ecs/"},
			{Title: "AWS Lambda documentation", URL: "https://docs.aws.amazon.com/lambda/"},
		}
	case "database":
		return []DocLink{
			{Title: "Amazon RDS documentation", URL: "https://docs.aws.amazon.com/rds/"},
			{Title: "Amazon DynamoDB documentation", URL: "https://docs.aws.amazon.com/dynamodb/"},
		}
	default:
		return nil
	}
}
