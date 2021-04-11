# GoWS

For now this library is made for managing multiple Rooms. Users join a room / 'game' and when they leave that room they get disconnected from the WS. However, adding functionality for switching rooms is half baked in.

### Features currently supported:
- Redis Pub/Sub
- Reading / Writing to a room

### Features that can be easily implemented:
- Switching rooms

### Planned:
- Adding Supabase for auth / user storage
- Making stuff more organized

### Files:

##### Client.go
Manages writing to user and reading from user.

##### Hub.go
Manages routing users to proper rooms.