package config

// **********************************************************
// This is where you add the parameters you want to retrieve

const (
	PlatformEventsSnsTopicArnParameterName       ConfigParameterName = "worker-tasks-topic-arn"
	EmailSubscriptionsTableNameParameterName     ConfigParameterName = "email-subscriptions-table-name"
	RelationalDatabaseDSNParameterName           ConfigParameterName = "planetscale-core-dsn"
	WebAppDomainNameParameterName                ConfigParameterName = "web-app-domain-name"
	MainTransactionalSendingAddressParameterName ConfigParameterName = "main-transactional-sending-address"
)

var paramDefinitions = []ConfigParameterDefinition{
	{
		ParameterName: PlatformEventsSnsTopicArnParameterName,
		ParameterType: StandardParameter,
	},
	{
		ParameterName: EmailSubscriptionsTableNameParameterName,
		ParameterType: StandardParameter,
	},
	{
		ParameterName: RelationalDatabaseDSNParameterName,
		ParameterType: SecretParameter,
	},
	{
		ParameterName: WebAppDomainNameParameterName,
		ParameterType: StandardParameter,
	},
	{
		ParameterName: MainTransactionalSendingAddressParameterName,
		ParameterType: StandardParameter,
	},
}

func (s *Config) PlatformEventsTopicArn() string {
	return s.parameters[PlatformEventsSnsTopicArnParameterName]
}

func (s *Config) EmailSubscriptionsTableName() string {
	return s.parameters[EmailSubscriptionsTableNameParameterName]
}

func (s *Config) RelationalDatabaseConnectionString() string {
	return s.parameters[RelationalDatabaseDSNParameterName]
}

func (s *Config) WebAppDomainName() string {
	return s.parameters[WebAppDomainNameParameterName]
}

func (s *Config) MainTransactionalSendingAddress() string {
	return s.parameters[MainTransactionalSendingAddressParameterName]
}

// **********************************************************
