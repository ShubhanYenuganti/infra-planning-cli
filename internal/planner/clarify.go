package planner

func ClarificationQuestion(req Request) string {
	if req.PrimaryDomain == DomainUnknown {
		return "Which AWS area should this plan focus on first: networking, compute, or database?"
	}
	return ""
}
