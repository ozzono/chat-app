<template>
  <div class="home">
    <h1>Chat Rooms</h1>
    <button @click="createRoom">Create Room</button>
    <button @click="fetchRooms">Refresh Rooms</button>
    <ul>
      <li v-for="room in rooms" :key="room">
        {{ room }}
        <button @click="joinRoom(room)">Join</button>
      </li>
    </ul>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      rooms: [],
    };
  },
  created() {
    this.fetchRooms();
  },
  methods: {
    fetchRooms() {
      axios.get('/api/v1/rooms')
        .then(response => {
          this.rooms = response.data; // Ensure the correct data structure
        })
        .catch(error => {
          console.error('Error fetching rooms:', error);
        });
    },
    createRoom() {
      const roomID = prompt('Enter room id:');
      if (roomID) {
        axios.post('/api/v1/rooms', { id: roomID })
          .then(() => {
            this.fetchRooms();
          })
          .catch(error => {
            console.error('Error creating room:', error);
          });
      }
    },
    joinRoom(room) {
      const nickname = prompt('Enter your nickname:');
      if (nickname) {
        window.location.href = `/rooms/${room}/bind?nickname=${nickname}`;
      }
    },
  },
};
</script>

<style scoped>
.home {
  text-align: center;
}
</style>
