package model

type HostInfo struct {
	Veid            string // 搬瓦工主机识别id
	ApiKey          string // 搬瓦工主机api key
	Hostname        string // 主机名称
	NodeLocation    string // 地理位置
	IpAddress       string // IP地址
	PlanMonthlyData int64  // 每月流量
	DataCounter     int64  // 已用流量
	DataNextReset   int64  // 下个周期（秒时间戳）
	LastUsedData    int64  // 昨日用量（若在同计费周期内，则为"此次已用流量-上次已用流量"；若是新的一个计费周期，则为"此次已用流量"）
	NewDataReset    bool   // 是否为新开启的计费周期
	Error           int    // 错误码
}

func (h *HostInfo) UpdateHostInfo(res BwhApiGetServiceInfoResponse) {
	h.Error = res.Error
	// 没错误才会更新，否则只更新Error码
	if res.Error == 0 {
		h.Hostname = res.Hostname
		h.NodeLocation = res.NodeLocation
		if len(res.IpAddresses) > 0 {
			h.IpAddress = res.IpAddresses[0]
		} else {
			h.IpAddress = ""
		}
		h.PlanMonthlyData = res.PlanMonthlyData
		if h.DataNextReset == 0 { // 首次记录
			h.NewDataReset = true
			h.LastUsedData = 0
		} else if res.DataNextReset > h.DataNextReset { // 新计费周期
			h.NewDataReset = true
			h.LastUsedData = res.DataCounter
		} else { // 同计费周期
			h.NewDataReset = false
			h.LastUsedData = res.DataCounter - h.DataCounter
		}
		h.DataNextReset = res.DataNextReset
		h.DataCounter = res.DataCounter
	}
}

type BwhApiGetServiceInfoResponse struct { // 搬瓦工API请求返回
	VmType                          string   `json:"vm_type"`                            // Hypervizor type (ovz or kvm)
	Hostname                        string   `json:"hostname"`                           // * Hostname of the VPS
	NodeIp                          string   `json:"node_ip"`                            // address of the physical node
	NodeAlias                       string   `json:"node_alias"`                         // Internal nickname of the physical node
	NodeLocationId                  string   `json:"node_location_id"`                   // Physical location id
	NodeLocation                    string   `json:"node_location"`                      // * Physical location (country, state)
	NodeDatacenter                  string   `json:"node_datacenter"`                    // Physical location of datacenter
	LocationIpv6Ready               bool     `json:"location_ipv6_ready"`                // Whether IPv6 is supported at the current location
	Plan                            string   `json:"plan"`                               // Name of plan
	PlanMonthlyData                 int64    `json:"plan_monthly_data"`                  // * Allowed monthly data transfer (bytes). Needs to be multiplied by monthly_data_multiplier - see below.
	MonthlyDataMultiplier           int      `json:"monthly_data_multiplier"`            // Some locations offer more expensive bandwidth; this variable contains the bandwidth accounting coefficient.
	PlanDisk                        int64    `json:"plan_disk"`                          // Disk quota (bytes)
	PlanRam                         int      `json:"plan_ram"`                           // RAM (bytes)
	PlanSwap                        int      `json:"plan_swap"`                          // SWAP (bytes)
	PlanMaxIpv6S                    int      `json:"plan_max_ipv6s"`                     // Maximum number of IPv6 addresses allowed by plan
	Os                              string   `json:"os"`                                 // Operating system
	Email                           string   `json:"email"`                              // Primary e-mail address of the account
	DataCounter                     int64    `json:"data_counter"`                       // * Data transfer used in the current billing month. Needs to be multiplied by monthly_data_multiplier - see below.
	DataNextReset                   int64    `json:"data_next_reset"`                    // * Date and time of transfer counter reset (UNIX timestamp)
	IpAddresses                     []string `json:"ip_addresses"`                       // * IPv4 and IPv6 addresses assigned to VPS (Array)
	PrivateIpAddresses              []string `json:"private_ip_addresses"`               // Private IPv4 addresses assigned to VPS (Array)
	PlanPrivateNetworkAvailable     bool     `json:"plan_private_network_available"`     // Whether or not Private Network features are available on this plan
	LocationPrivateNetworkAvailable bool     `json:"location_private_network_available"` // Whether or not Private Network features are available at this location
	RdnsApiAvailable                bool     `json:"rdns_api_available"`                 // Whether or not rDNS records can be set via API
	Suspended                       bool     `json:"suspended"`                          // Whether VPS is suspended
	PolicyViolation                 bool     `json:"policy_violation"`                   // Whether there is an active policy violation that needs attention (see getPolicyViolations)
	TotalAbusePoints                int      `json:"total_abuse_points"`                 // Total abuse points accumulated in current calendar year
	MaxAbusePoints                  int      `json:"max_abuse_points"`                   // Maximum abuse points allowed by plan in a calendar year
	Error                           int      `json:"error"`                              // Error code. 0 when no error
	Message                         string   `json:"message"`                            // Error message
}
