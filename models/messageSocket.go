package models

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"

type MessageSocket struct {
	Action string `json:"action"`
	// Target *RoomSocket `json:"target"`
	// Sender *Client     `json:"sender"`

	Message
}
