<script setup lang="ts">
import { ref, reactive, watchEffect } from 'vue';
import { useRouter } from 'vue-router';
import { useRoute } from 'vue-router';
import { useStore } from '@/stores';

type Message = {
  content: string
  clientId: string
  roomId: string
  userName: string
  type: 'recv' | 'self'
}

type Subscriber = {
  id: string
  userName: string
}
const users = reactive<Subscriber[]>([])

const router = useRouter()
const route = useRoute()
const roomId = route.params.id as string
console.log(roomId)
const { conn } = useStore()
const msgInput = ref("")

watchEffect(async () => {
  if (conn === null) {
    router.push("/")
    return
  }
  try {
    const response = await fetch(`http://localhost:5000/getSubscriber/${roomId}`, {
      method: "GET",
      credentials: "include",
    });

    if (response.status === 200) {
      const data: Subscriber[] = await response.json()
      console.log(data)
      Object.assign(users, data)
      console.log(users)
    }
  } catch (e) {
    console.log(e)
  }
})

watchEffect(() => {
  if (conn === null) {
    router.push("/")
    return
  }

  conn.onmessage = (message) => {
    const m: Message = JSON.parse(message.data)
    if (m.content == "A new user has joined the room") {
      users.push({
        id: m.clientId,
        userName: m.userName
      })
    }

    if (m.content == "user left the chat") {
      const deleteUser = users.filter((user) => user.userName != m.userName)
      Object.assign(users, deleteUser)
      messages.push(m)
      return
    }

    const currentUser = sessionStorage.getItem("username") as string
    currentUser == m.userName ? (m.type = 'self') : (m.type = 'recv')
    messages.push(m)
    console.log(users)
  }
  conn.onclose = () => { }
  conn.onerror = () => { }
  conn.onopen = () => { }
})

const sendMessage = async () => {
  if (!msgInput.value) return
  if (conn === null) {
    router.push("/")
    return
  }

  conn.send(msgInput.value)
  msgInput.value = ""
};



const messages = reactive<Message[]>([])
</script>

<template>
  <div class="container mx-auto max-w-screen-lg ">
    <div class="flex flex-col h-[calc(100vh-74px)]">
      <div class="w-full flex-grow overflow-auto my-2 p-4">
        <div v-for="(msg, index) in   messages  " :key="index">
          <div :class="[msg.type === 'self' ? ['flex', 'flex-col', 'mt-2', 'w-full', 'text-right', 'justify-end'] : ['mt-2']]">
            <div class="text-sm">{{ msg.userName }}</div>
            <div>
              <div :class="[msg.type === 'self' ? ['bg-blue-500', 'text-white', 'px-4', 'py-1', 'rounded-md', 'inline-block', 'mt-1'] : ['bg-gray-200', 'text-black', 'px-4', 'py-1', 'rounded-md', 'inline-block', 'mt-1']]">{{ msg.content }}</div>
            </div>
          </div>
        </div>
      </div>
      <div class="mt-4">
        <div class="flex px-4 py-2 bg-gray-100 rounded-md">
          <div class="flex w-full mr-4 rounded-md border border-blue-500">
            <textarea v-model="msgInput" class="w-full h-10 p-2 rounded-md focus:outline-none" placeholder="type your message here" style="resize: none;"></textarea>
          </div>
          <div class="flex items-center">
            <button @click="sendMessage()" class="p-2 rounded-md bg-blue-500 text-white">Send</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>