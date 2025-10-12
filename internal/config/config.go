package config

import (
	"encoding/json"

	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/orm"
)

type Config struct {
	Model         string          `yaml:"Model" default:"prod"`
	Host          string          `yaml:"Host" default:"0.0.0.0"`
	Port          int             `yaml:"Port" default:"8080"`
	Debug         bool            `yaml:"Debug" default:"false"`
	TLS           TLS             `yaml:"TLS"`
	JwtAuth       JwtAuth         `yaml:"JwtAuth"`
	Logger        logger.LogConf  `yaml:"Logger"`
	MySQL         orm.Config      `yaml:"MySQL"`
	Redis         RedisConfig     `yaml:"Redis"`
	Site          SiteConfig      `yaml:"Site"`
	Node          NodeConfig      `yaml:"Node"`
	Mobile        MobileConfig    `yaml:"Mobile"`
	Email         EmailConfig     `yaml:"Email"`
	Device        DeviceConfig    `yaml:"device"`
	Verify        Verify          `yaml:"Verify"`
	VerifyCode    VerifyCode      `yaml:"VerifyCode"`
	Register      RegisterConfig  `yaml:"Register"`
	Subscribe     SubscribeConfig `yaml:"Subscribe"`
	Invite        InviteConfig    `yaml:"Invite"`
	Telegram      Telegram        `yaml:"Telegram"`
	Log           Log             `yaml:"Log"`
	Administrator struct {
		Email    string `yaml:"Email" default:"admin@ppanel.dev"`
		Password string `yaml:"Password" default:"password"`
	} `yaml:"Administrator"`
}

type RedisConfig struct {
	Host string `yaml:"Host" default:"localhost:6379"`
	Pass string `yaml:"Pass" default:""`
	DB   int    `yaml:"DB" default:"0"`
}

type JwtAuth struct {
	AccessSecret string `yaml:"AccessSecret"`
	AccessExpire int64  `yaml:"AccessExpire" default:"604800"`
}

type Verify struct {
	TurnstileSiteKey    string `yaml:"TurnstileSiteKey" default:""`
	TurnstileSecret     string `yaml:"TurnstileSecret" default:""`
	LoginVerify         bool   `yaml:"LoginVerify" default:"false"`
	RegisterVerify      bool   `yaml:"RegisterVerify" default:"false"`
	ResetPasswordVerify bool   `yaml:"ResetPasswordVerify" default:"false"`
}

type SubscribeConfig struct {
	SingleModel     bool   `yaml:"SingleModel" default:"false"`
	SubscribePath   string `yaml:"SubscribePath" default:"/v1/subscribe/config"`
	SubscribeDomain string `yaml:"SubscribeDomain" default:""`
	PanDomain       bool   `yaml:"PanDomain" default:"false"`
	UserAgentLimit  bool   `yaml:"UserAgentLimit" default:"false"`
	UserAgentList   string `yaml:"UserAgentList" default:""`
}

type RegisterConfig struct {
	StopRegister            bool   `yaml:"StopRegister" default:"false"`
	EnableTrial             bool   `yaml:"EnableTrial" default:"false"`
	TrialSubscribe          int64  `yaml:"TrialSubscribe" default:"0"`
	TrialTime               int64  `yaml:"TrialTime" default:"0"`
	TrialTimeUnit           string `yaml:"TrialTimeUnit" default:""`
	IpRegisterLimit         int64  `yaml:"IpRegisterLimit" default:"0"`
	IpRegisterLimitDuration int64  `yaml:"IpRegisterLimitDuration" default:"0"`
	EnableIpRegisterLimit   bool   `yaml:"EnableIpRegisterLimit" default:"false"`
}

type EmailConfig struct {
	Enable                     bool   `yaml:"Enable" default:"true"`
	Platform                   string `yaml:"platform"`
	PlatformConfig             string `yaml:"platform_config"`
	EnableVerify               bool   `yaml:"enable_verify"`
	EnableNotify               bool   `yaml:"enable_notify"`
	EnableDomainSuffix         bool   `yaml:"enable_domain_suffix"`
	DomainSuffixList           string `yaml:"domain_suffix_list"`
	VerifyEmailTemplate        string `yaml:"verify_email_template"`
	ExpirationEmailTemplate    string `yaml:"expiration_email_template"`
	MaintenanceEmailTemplate   string `yaml:"maintenance_email_template"`
	TrafficExceedEmailTemplate string `yaml:"traffic_exceed_email_template"`
}

type MobileConfig struct {
	Enable          bool     `yaml:"Enable" default:"true"`
	Platform        string   `yaml:"platform"`
	PlatformConfig  string   `yaml:"platform_config"`
	EnableVerify    bool     `yaml:"enable_verify"`
	EnableWhitelist bool     `yaml:"enable_whitelist"`
	Whitelist       []string `yaml:"whitelist"`
}

type DeviceConfig struct {
	Enable         bool   `yaml:"enable" default:"true"`
	ShowAds        bool   `yaml:"show_ads"`
	EnableSecurity bool   `yaml:"enable_security"`
	OnlyRealDevice bool   `yaml:"only_real_device"`
	SecuritySecret string `yaml:"security_secret"`
}

type SiteConfig struct {
	Host       string `yaml:"Host" default:""`
	SiteName   string `yaml:"SiteName" default:""`
	SiteDesc   string `yaml:"SiteDesc" default:""`
	SiteLogo   string `yaml:"SiteLogo" default:""`
	Keywords   string `yaml:"Keywords" default:""`
	CustomHTML string `yaml:"CustomHTML" default:""`
	CustomData string `yaml:"CustomData" default:""`
}

type NodeConfig struct {
	NodeSecret             string         `yaml:"NodeSecret" default:""`
	NodePullInterval       int64          `yaml:"NodePullInterval" default:"60"`
	NodePushInterval       int64          `yaml:"NodePushInterval" default:"60"`
	TrafficReportThreshold int64          `yaml:"TrafficReportThreshold" default:"0"`
	IPStrategy             string         `yaml:"IPStrategy" default:""`
	DNS                    []NodeDNS      `yaml:"DNS"`
	Block                  []string       `yaml:"Block" `
	Outbound               []NodeOutbound `yaml:"Outbound"`
}

func (n *NodeConfig) Marshal() ([]byte, error) {
	type Alias NodeConfig
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}

func (n *NodeConfig) Unmarshal(data []byte) error {
	type Alias NodeConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	return json.Unmarshal(data, &aux)
}

type NodeDNS struct {
	Proto   string   `json:"proto"`
	Address string   `json:"address"`
	Domains []string `json:"domains"`
}

func (n *NodeDNS) Marshal() ([]byte, error) {
	type Alias NodeDNS
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}

func (n *NodeDNS) Unmarshal(data []byte) error {
	type Alias NodeDNS
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	return json.Unmarshal(data, &aux)
}

type NodeOutbound struct {
	Name     string   `json:"name"`
	Protocol string   `json:"protocol"`
	Address  string   `json:"address"`
	Port     int64    `json:"port"`
	Password string   `json:"password"`
	Rules    []string `json:"rules"`
}

func (n *NodeOutbound) Marshal() ([]byte, error) {
	type Alias NodeOutbound
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}

type File struct {
	Host    string         `yaml:"Host" default:"0.0.0.0"`
	Port    int            `yaml:"Port" default:"8080"`
	TLS     TLS            `yaml:"TLS"`
	Debug   bool           `yaml:"Debug" default:"true"`
	JwtAuth JwtAuth        `yaml:"JwtAuth"`
	Logger  logger.LogConf `yaml:"Logger"`
	MySQL   orm.Config     `yaml:"MySQL"`
	Redis   RedisConfig    `yaml:"Redis"`
}

type InviteConfig struct {
	ForcedInvite       bool  `yaml:"ForcedInvite" default:"false"`
	ReferralPercentage int64 `yaml:"ReferralPercentage" default:"0"`
	OnlyFirstPurchase  bool  `yaml:"OnlyFirstPurchase" default:"false"`
}

type Telegram struct {
	Enable        bool   `yaml:"Enable" default:"false"`
	BotID         int64  `yaml:"BotID" default:""`
	BotName       string `yaml:"BotName" default:""`
	BotToken      string `yaml:"BotToken" default:""`
	EnableNotify  bool   `yaml:"EnableNotify" default:"false"`
	WebHookDomain string `yaml:"WebHookDomain" default:""`
}

type TLS struct {
	Enable   bool   `yaml:"Enable" default:"false"`
	CertFile string `yaml:"CertFile" default:""`
	KeyFile  string `yaml:"KeyFile" default:""`
}

type VerifyCode struct {
	ExpireTime int64 `yaml:"ExpireTime" default:"300"`
	Limit      int64 `yaml:"Limit" default:"15"`
	Interval   int64 `yaml:"Interval" default:"60"`
}

type Log struct {
	AutoClear bool  `yaml:"AutoClear" default:"true"`
	ClearDays int64 `yaml:"ClearDays" default:"7"`
}

type NodeDBConfig struct {
	NodeSecret             string
	NodePullInterval       int64
	NodePushInterval       int64
	TrafficReportThreshold int64
	IPStrategy             string
	DNS                    string
	Block                  string
	Outbound               string
}
