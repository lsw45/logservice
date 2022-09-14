package entity

type NotifyDeployMessage struct {
	UUID    string `json:"uuid"`
	Title   string `json:"title"`
	Content []struct {
		Servers    []Servers `json:"servers"`
		RegionID   int       `json:"region_id"`
		RegionName string    `json:"region_name"`
	} `json:"content"`
}

type Servers struct {
	ID    int    `json:"id"`
	Os    string `json:"os"`
	CPU   string `json:"cpu"`
	Env   int    `json:"env"`
	Mem   string `json:"mem"`
	Name  string `json:"name"`
	IPObj []struct {
		ID      int    `json:"id"`
		IP      string `json:"ip"`
		Subnet  int    `json:"subnet"`
		Private bool   `json:"private"`
		Version string `json:"version"`
	} `json:"ip_obj"`
	Status        string  `json:"status"`
	Creator       Creator `json:"creator"`
	Deleted       bool    `json:"deleted"`
	Project       int     `json:"project"`
	AddTime       string  `json:"add_time"`
	Password      string  `json:"password"`
	Username      string  `json:"username"`
	RemoteID      string  `json:"remote_id"`
	ResourceID    int     `json:"resource_id"`
	UpdateTime    string  `json:"update_time"`
	CorporationID string  `json:"corporation_id"`
}

type Creator struct {
	ID              int           `json:"id"`
	UID             int           `json:"uid"`
	Email           string        `json:"email"`
	Avatar          string        `json:"avatar"`
	Groups          []interface{} `json:"groups"`
	Mobile          interface{}   `json:"mobile"`
	Status          string        `json:"status"`
	Company         int           `json:"company"`
	Country         string        `json:"country"`
	Creator         interface{}   `json:"creator"`
	Deleted         bool          `json:"deleted"`
	AddTime         string        `json:"add_time"`
	IsStaff         bool          `json:"is_staff"`
	Nickname        interface{}   `json:"nickname"`
	Password        string        `json:"password"`
	Username        string        `json:"username"`
	IsActive        bool          `json:"is_active"`
	LastName        string        `json:"last_name"`
	FirstName       string        `json:"first_name"`
	LastLogin       interface{}   `json:"last_login"`
	DateJoined      string        `json:"date_joined"`
	UpdateTime      string        `json:"update_time"`
	AccessToken     string        `json:"access_token"`
	IsSuperuser     bool          `json:"is_superuser"`
	SessionToken    string        `json:"session_token"`
	CorporationID   string        `json:"corporation_id"`
	UserPermissions []interface{} `json:"user_permissions"`
}

type UpdateFileReq struct {
	Remote   string
	Server   string
	Preserve bool
}

type TunnelUploadFileRes struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
	} `json:"data"`
}
