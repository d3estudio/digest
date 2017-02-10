package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Client provides an abstration to the Slack Real-Time API. After connecting,
// it caches known Channels and Users in order to perform querying operations
// and replacing Users and Channels IDs from incoming data with their respective
// instances.
type Client struct {
	Token        string
	Users        map[string]RTMUser
	Channels     map[string]RTMChannel
	WebSocketURL string
	Channel      chan RTMMessage
}

const (
	// TypeMessage indicates that incoming message type is "message"
	TypeMessage = "message"

	// TypeMessageDeleted indicates that incoming message type is
	// "message_deleted"
	TypeMessageDeleted = "message_deleted"

	// TypeMessageChanged indicates that the incoming message type is
	// "message_changed"
	TypeMessageChanged = "message_changed"

	// TypeReactionAdded indicates that incoming message type is "reaction_added"
	TypeReactionAdded = "reaction_added"

	// TypeReactionRemoved indicates that incoming message type is
	// "reaction_removed"
	TypeReactionRemoved = "reaction_removed"

	// TypeEmojiChanged indicates that incoming message type is "emoji_changed"
	TypeEmojiChanged = "emoji_changed"
)

var watchedMessageTypes = []string{
	TypeMessage,
	TypeMessageDeleted,
	TypeMessageChanged,
	TypeReactionAdded,
	TypeReactionRemoved,
	TypeEmojiChanged,
}

// Handshake performs the initial handshake process with the Slack Real-Time API
// server using the provided Token string.
func (c *Client) handshake() error {
	c.Users = make(map[string]RTMUser)
	c.Channels = make(map[string]RTMChannel)

	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", c.Token)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var respObj rtmStartResponse
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return err
	}
	for _, i := range respObj.Users {
		if i.Deleted {
			continue
		}
		c.Users[i.ID] = i
	}

	for _, i := range respObj.Channels {
		c.Channels[i.ID] = i
	}

	for _, i := range respObj.Groups {
		i.IsMember = true
		c.Channels[i.ID] = i
	}

	c.WebSocketURL = respObj.WebsocketURL
	return nil
}

func readMessageFromWS(ws *websocket.Conn) (m RTMMessage, err error) {
	var data []byte
	err = websocket.Message.Receive(ws, &data)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}

	if m.SubType != "" && m.SubType == TypeMessageChanged {
		var messageChanged rtmMessageUpdate
		err = json.Unmarshal(data, &messageChanged)
		if err != nil {
			return
		}
		m.Message = messageChanged.Message.Message
		m.Timestamp = messageChanged.Message.Timestamp
		m.UserID = messageChanged.Message.UserID
		return
	}

	var txtMessage rtmTextMessage
	err = json.Unmarshal(data, &txtMessage)
	if err != nil {
		return
	}
	m.Message = txtMessage.Message
	return
}

// Listen initializes the WebSocket connection with the remote RTM server,
// and pushes incoming messages into the provided channel.
func (c *Client) Listen() {
	attemptNumber := 1
outerLoop:
	for {
		log.Info("Attempting handshake... (Attempt ", attemptNumber, ")")
		if attemptNumber >= 3 {
			backoff := time.Duration(attemptNumber) * time.Second
			log.Warn("Backing off by ", attemptNumber, " seconds")
			time.Sleep(backoff)
		}
		err := c.handshake()
		if err != nil {
			log.Error("Handshake failed: ", err)
			attemptNumber++
			continue
		}
		log.Info("Handshake succeeded.")
		log.Info("Dialing to ", c.WebSocketURL)
		ws, err := websocket.Dial(c.WebSocketURL, "", "https://api.slack.com/")
		if err != nil {
			log.Error("WebSocket Dial failed: ", err)
			attemptNumber++
			continue
		}
		log.Info("Connection succeeded.")
		attemptNumber = 1
		log.Info("Entering message loop...")
		for {
			msg, err := readMessageFromWS(ws)
			if err != nil {
				if err == io.EOF {
					log.Error("WebSocket Connection Interrupted. Performing Handshake... Error was: ", err)
					continue outerLoop
				} else {
					log.Error(err)
					continue
				}
			}

			skip := true
			msg.Client = *c
			if user, ok := c.Users[msg.UserID]; ok {
				msg.User = user
			}

			if channel, ok := c.Channels[msg.ChannelID]; ok {
				msg.Channel = channel
			}

			for _, t := range watchedMessageTypes {
				if t == msg.Type {
					skip = false
					break
				}
			}

			if skip {
				continue
			}

			// We actually want to move the SubType data into the Type, so
			// parsing is easier on our end, but not when we're dealing
			// with an emoji_changed message.
			if msg.SubType != "" && msg.Type != TypeEmojiChanged {
				msg.Type = msg.SubType
			}

			logMessage(msg)
			c.Channel <- msg
		}
	}
}

func logMessage(m RTMMessage) {
	fields := log.Fields{
		"type":    m.Type,
		"channel": m.ChannelID,
		"user":    m.UserID,
	}
	if m.Type == TypeReactionAdded || m.Type == TypeReactionRemoved {
		fields["item"] = m.Item.Timestamp
		fields["reaction"] = m.Reaction
	}
	log.WithFields(fields).Info("Message received")
}
