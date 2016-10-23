package slack

type rtmStartResponse struct {
	Ok           bool         `json:"ok"`
	Error        string       `json:"error"`
	WebsocketURL string       `json:"url"`
	Self         rtmStartSelf `json:"self"`
	Channels     []RTMChannel `json:"channels"`
	Groups       []RTMChannel `json:"groups"`
	Users        []RTMUser    `json:"users"`
}

type rtmStartSelf struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RTMChannel represents an channel or group from the Team account
type RTMChannel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsChannel bool   `json:"is_channel"`
	IsMember  bool   `json:"is_member"`
}

// RTMUser represents an User that belongs to the Team account
type RTMUser struct {
	ID       string         `json:"id"`
	Username string         `json:"name"`
	Deleted  bool           `json:"deleted"`
	Profile  RTMUserProfile `json:"profile"`
}

// RTMUserProfile represents the Profile section of a given RTMUser
type RTMUserProfile struct {
	RealName       string `json:"real_name"`
	ProfilePicture string `json:"image_192"`
	Title          string `json:"title"`
}

// RTMMessage is the common structure used to unmarshal incoming common data
// from the remote RTM server
type RTMMessage struct {
	Type           string      `json:"type"`
	SubType        string      `json:"subtype"`
	Message        string      `json:"txt"`
	DeletionTarget string      `json:"deleted_ts"`
	EventTimestamp string      `json:"event_ts"`
	ChannelID      string      `json:"channel"`
	UserID         string      `json:"user"`
	Timestamp      string      `json:"ts"`
	Reaction       string      `json:"reaction"`
	Item           *RTMMessage `json:"item"`
	Client         Client
	Channel        RTMChannel
	User           RTMUser
}

type rtmMessageUpdate struct {
	Message rtmMessageUpdateInner `json:"message"`
}

type rtmMessageUpdateInner struct {
	Message   string `json:"text"`
	Timestamp string `json:"ts"`
	UserID    string `json:"user"`
}

type rtmTextMessage struct {
	Message string `json:"text"`
}
