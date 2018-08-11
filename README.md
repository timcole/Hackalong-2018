# Hackalong 2018
[Hackathon Information](https://hackalong.devpost.com/)

---

## Team
 - [Timothy Cole](https://www.twitch.tv/modesttim) *Back-end*
 - [Jamie Pine](https://www.twitch.tv/jamiepinelive) *Front-end*
 - [Chase Aaron](https://www.twitch.tv/doubleayeeron) *Topic curator*

---

## Idea
Create a random chat room website where users can either create a room with a set topic or join a random room where they'll be thrown into a chat with an odd number of people to talk about the set topic.

Rooms will start with 3 people and as the conversation duration progress more people will be thrown into the chat to keep it going. There will always be an odd number of people so the room can't be 50/50.

When a new user gets thrown into the chat it will feed them the previous messages so they can quickly catch up on the talk.

Members of the chat will have buttons that will allow to pick their side at any given time so everyone can see the diversity on the topic.

---

## Backend Documenation

### WebSocket Server
Pathname: `/ws`

To the server
 - `SET_USERNAME` Requires `data.username` - `^[a-zA-Z0-9_]+$`
	- Setings Your Username
	- e.g. `{"type":"SET_USERNAME", "data": { "username": "Testing321" }}`
 - `CREATE_CHANNEL` Requires `data.topic`
	- Creates a new channel with a topic
	- e.g. `{"type":"CREATE_CHANNEL", "data": { "topic": "Windows vs Mac LUL" }}`
 - `JOIN_CHANNEL`
	- Joins a random channel
	- e.g. `{"type":"JOIN_CHANNEL" }`
 - `LEAVE_CHANNEL`
	- Leaves an active channel
	- e.g. `{"type":"LEAVE_CHANNEL" }`
 - `SEND_MESSAGE` Requires `data.message`
	- Sends a message in the active channel
	- e.g. `{"type":"SEND_MESSAGE", "data": { "message": "Kappa 123" }}`

From the server
 - `SET_USERNAME`
	- Response from settings a username can either contain `data.username` or `error`
	- e.g. `{"type":"SET_USERNAME", "data": { "username": "Testing321" }}`
	- e.g. `{"type":"SET_USERNAME","error":"ERR_BADUSERNAME"}`
 - `CREATE_CHANNEL`
	- Sent when there was an error creating a new channel. If successful will send `JOIN_CHANNEL` instead
	- e.g. `{"type":"CREATE_CHANNEL","error":"ERR_UNAUTHORIZED"}`
 - `JOIN_CHANNEL`
	- Sent when you've successfully joined a channel from with sending `JOIN_CHANNEL` or `CREATE_CHANNEL`
	- e.g. `{"type":"JOIN_CHANNEL"}`
 - `NEW_MESSAGE`
	- Sent when there is a new message in the chat from either you or other members
	- e.g. `{"type":"NEW_MESSAGE","data":{"message":"Kappa 123","username":"Testing321"}}`