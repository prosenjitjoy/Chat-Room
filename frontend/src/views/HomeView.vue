<script setup lang="ts">
import { ref, reactive, onMounted, watchEffect } from 'vue';
import { v4 as uuid } from 'uuid';
import { useRouter } from 'vue-router';
import { useStore } from '@/stores';

type Room = {
  id: string
  name: string
}

const router = useRouter()

const username = sessionStorage.getItem("username")
const roomName = ref("")

onMounted(() => {
  getRooms()
})


const rooms = reactive<Room[]>([])

const submit = async () => {
  try {
    const newRoom: Room = {
      id: uuid(),
      name: roomName.value
    }

    const response = await fetch("http://localhost:5000/createRoom", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify(newRoom),
    });

    if (response.status === 200) {
      getRooms()
    }
  } catch (e) {
    console.log(e)
  }
};

const getRooms = async () => {
  try {
    const response = await fetch("http://localhost:5000/getRooms", {
      method: "GET",
      credentials: "include",
    });

    if (response.status === 200) {
      const data: Room[] = await response.json()
      console.log(data)
      Object.assign(rooms, data)
      console.log(rooms)
    }
  } catch (e) {
    console.log(e)
  }
};

const joinRoom = async (roomID: string) => {
  const userID = sessionStorage.getItem("userid") as string
  const userName = sessionStorage.getItem("username") as string
  const ws = new WebSocket(`ws://localhost:5000/joinRoom/${roomID}?userId=${userID}&userName=${userName}`)

  if (ws.OPEN) {
    useStore().conn = ws
    router.push(`/room/${roomID}`)
    return
  }
}

</script>

<template>
  <div class="container mx-auto max-w-screen-xl">
    <div class="my-8 px-4 md:mx-32">
      <div class="flex justify-center mt-3 p-5">
        <input v-model="roomName" type="text" class="border border-gray-300 p-2 rounded-md focus:outline-none focus:border-blue-500" placeholder="room name">
        <button @click="submit" type="submit" class="bg-blue-500 border text-white rounded-md p-2 md:ml-4 hover:bg-blue-600 active:bg-blue-500">Create Room</button>
      </div>

      <div class="mt-6">
        <div class="font-bold">Available Rooms</div>
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mt-6">
          <div v-for="(room, index) in rooms" :key="index" class="border border-blue-500 p-4 flex items-center rounded-md w-full">
            <div :id=room.id class="w-full">
              <div class="text-sm">room</div>
              <div class="text-blue font-bold text-lg break-words">
                {{ room.name.slice(0, 10) }}
              </div>
            </div>
            <div>
              <button @click="joinRoom(room.id)" class="px-4 text-white bg-blue-500 rounded-md">join</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
