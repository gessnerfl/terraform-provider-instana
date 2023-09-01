package restapi

// AlertingChannelType type of the alerting channel
type AlertingChannelType string

const (
	//EmailChannelType constant value for alerting channel type EMAIL
	EmailChannelType = AlertingChannelType("EMAIL")
	//GoogleChatChannelType constant value for alerting channel type GOOGLE_CHAT
	GoogleChatChannelType = AlertingChannelType("GOOGLE_CHAT")
	//Office365ChannelType constant value for alerting channel type OFFICE_365
	Office365ChannelType = AlertingChannelType("OFFICE_365")
	//OpsGenieChannelType constant value for alerting channel type OPS_GENIE
	OpsGenieChannelType = AlertingChannelType("OPS_GENIE")
	//PagerDutyChannelType constant value for alerting channel type PAGER_DUTY
	PagerDutyChannelType = AlertingChannelType("PAGER_DUTY")
	//SlackChannelType constant value for alerting channel type SLACK
	SlackChannelType = AlertingChannelType("SLACK")
	//SplunkChannelType constant value for alerting channel type SPLUNK
	SplunkChannelType = AlertingChannelType("SPLUNK")
	//VictorOpsChannelType constant value for alerting channel type VICTOR_OPS
	VictorOpsChannelType = AlertingChannelType("VICTOR_OPS")
	//WebhookChannelType constant value for alerting channel type WEB_HOOK
	WebhookChannelType = AlertingChannelType("WEB_HOOK")
)
