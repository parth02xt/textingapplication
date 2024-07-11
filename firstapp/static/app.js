let currentRoom = '';

document.addEventListener('DOMContentLoaded', (event) => {
    document.getElementById('message').disabled = true;
    document.getElementById('sendButton').disabled = true;

    // Check for new messages every 5 seconds
    setInterval(() => {
        if (currentRoom) {
            loadMessages();
        }
    }, 5000);
});

function loadRoomList() {
    fetch('/rooms')
        .then(response => response.json())
        .then(data => {
            let roomList = document.getElementById('roomList');
            roomList.innerHTML = '';
            data.forEach(room => {
                let roomElement = document.createElement('li');
                roomElement.textContent = room;
                roomElement.onclick = () => joinRoom(room);
                roomList.appendChild(roomElement);
            });
        })
        .catch(error => console.error('Error loading rooms:', error));
}

function createRoom() {
    let roomName = document.getElementById('roomName').value;
    if (!roomName) {
        alert('Room name is required.');
        return;
    }

    fetch('/create-room', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ room: roomName }),
    })
    .then(response => {
        if (response.ok) {
            document.getElementById('roomName').value = '';
            loadRoomList();
        } else {
            throw new Error('Failed to create room');
        }
    })
    .catch(error => console.error('Error creating room:', error));
}

function joinRoom(room) {
    currentRoom = room;
    document.getElementById('currentRoom').textContent = room;
    loadMessages();
}

function loadMessages() {
    if (!currentRoom) {
        return;
    }

    fetch(`/receive?room=${currentRoom}`)
        .then(response => response.json())
        .then(data => {
            let messagesDiv = document.getElementById('messages');
            messagesDiv.innerHTML = '';
            data.forEach(msg => {
                let messageElement = document.createElement('div');
                messageElement.textContent = `${msg.username}: ${msg.content}`;
                messageElement.className = 'message ' + (msg.username === username ? 'username' : 'other');
                messagesDiv.appendChild(messageElement);
            });
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        })
        .catch(error => console.error('Error loading messages:', error));
}

function sendMessage() {
    let message = document.getElementById('message').value;

    if (!message) {
        alert('Message is required.');
        return;
    }

    fetch('/send', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, content: message, room: currentRoom }),
    })
    .then(response => {
        if (response.ok) {
            document.getElementById('message').value = '';
            loadMessages();
        } else {
            throw new Error('Failed to send message');
        }
    })
    .catch(error => console.error('Error sending message:', error));
}
