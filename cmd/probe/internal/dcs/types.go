package dcs

type Cluster struct {
	ClusterCompName string
	Replicas        int32
	HaConfig        *HaConfig
	Leader          *Leader
	OpTime          int64
	Members         []Member
	Switchover      *Switchover
	Extra           map[string]string
	resource        interface{}
}

func (c *Cluster) HasMember(memberName string) bool {
	for _, member := range c.Members {
		if memberName == member.Name {
			return true
		}
	}
	return false
}

func (c *Cluster) GetMemberWithName(name string) *Member {
	for _, m := range c.Members {
		if m.Name == name {
			return &m
		}
	}

	return nil
}

func (c *Cluster) GetMemberName() []string {
	var memberList []string
	for _, member := range c.Members {
		memberList = append(memberList, member.Name)
	}

	return memberList
}

func (c *Cluster) IsLocked() bool {
	return c.Leader != nil && c.Leader.Name != ""
}

func (c *Cluster) GetOpTime() int64 {
	return c.OpTime
}

type HaConfig struct {
	index              string
	ttl                int
	maxLagOnSwitchover int64
}

func (c *HaConfig) GetTtl() int {
	return c.ttl
}

func (c *HaConfig) GetMaxLagOnSwitchover() int64 {
	return c.maxLagOnSwitchover
}

type Leader struct {
	Index       string
	Name        string
	AcquireTime int64
	RenewTime   int64
	Ttl         int
	Resource    interface{}
}

type Member struct {
	Index          string
	Name           string
	Role           string
	PodIP          string
	DBPort         string
	SQLChannelPort string
}

func (m *Member) GetName() string {
	return m.Name
}

func newMember(index string, name string, role string, url string) *Member {
	return &Member{
		Index: index,
		Name:  name,
		Role:  role,
	}
}

type Switchover struct {
	Index       string
	Leader      string
	Candidate   string
	ScheduledAt int64
}

func newSwitchover(index string, leader string, candidate string, scheduledAt int64) *Switchover {
	return &Switchover{
		Index:       index,
		Leader:      leader,
		Candidate:   candidate,
		ScheduledAt: scheduledAt,
	}
}

func (s *Switchover) GetLeader() string {
	return s.Leader
}

func (s *Switchover) GetCandidate() string {
	return s.Candidate
}