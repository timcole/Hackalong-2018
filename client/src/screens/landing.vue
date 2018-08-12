<template>
    <div class="container">

                <transition name="fade">
                    <div v-if="usernameModal" @click="close" class="backdrop"></div>
                </transition>

                <transition name="slideup">
                  <div v-if="usernameModal" class="modal"> 
                    <div class="modal_content">
                    <h2>Oh hi mark!</h2>
                      <div class="dropshadow biginput">
                        <input v-model="username" type="text" placeholder="Type username here"/>
                      </div>
                      <div @click="createUsername" class="button dropshadow skew blue"><span class="noselect"> Go </span></div>
                  </div>
                  </div>
                </transition>


      <div class="logo">ChangeMyMind<span>.io</span></div>
    <div class="dropshadow biginput">
      <input v-model="question" type="text" placeholder="Type something controversial"/>
    </div>

    <div class="landing_actions">
      <div @click="createRoomButton" class="button dropshadow skew red"><span class="noselect"> Create Room </span></div>
      <h3 class="or">OR</h3>
      <div @click="joinRoomButton" class="button dropshadow skew blue"><span class="noselect"> Join Random Room </span></div>
   </div>

   <div class="landing_info">
     <p>Enter a controversial statment and get matched up with random people for a debate, as the conversation progresses, we’ll throw more randoms into your room to make things more intresting.</p>
   </div>





    </div>

</template>

<script>
export default {
  mounted() {

     // in case we were already in a channel
        console.log('Leaving channel...')
        if (this.socket) this.socket.send(JSON.stringify({
          type: "LEAVE_CHANNEL"
        }))

    // esablish connection with socket
    this.socket = new WebSocket('wss://56ce6fe7.ngrok.io/ws');
    this.socket.addEventListener('open', event => {
      console.log(`Connection established via socket`);
      this.connectionAlive = true
    })
    
    // Anything coming back from server is handled here
    this.socket.addEventListener('message', event => {
          var data = JSON.parse(event.data)
          console.log(data)

          if (data.type === 'SET_USERNAME') {
            this.joinOrCreateRoom()
          }

          if (data.type === 'JOIN_CHANNEL') {
            setTimeout(() => {
            // document.getElementById('fuckThisShit').setAttribute("value", data.data.topic);
            this.$root.$emit('topic', data.data.topic)
            }, 500)
            this.$router.push('/room')
          }

          if (data.type === 'RESPONSE') {
            alert('no free channels :(')
            this.usernameModal = false
          }

          this.$root.$emit('response', data)
    });
       

    window.socket = this.socket


  },
  data() {
    return {
      socket: null,
      connectionAlive: false,
      username: '',
      question: '',
      usernameModal: false,
      whatAreWeDoing: false, // JOIN, CREATE


    }
  },

    methods: {
      createRoomButton() {
        if (this.question.length < 3) return false
        this.whatAreWeDoing = 'CREATE'
        this.usernameModal = true
      },
      joinRoomButton() {
         this.whatAreWeDoing = 'JOIN'
          this.usernameModal = true
      },
      createUsernameModal(type) {
        this.usernameModalCreateRoom = true;
      },

      createUsername() {
          this.socket.send(JSON.stringify({
          type: "SET_USERNAME",
          data: {
            username: this.username
          }
        }))
      },

      // called after response from server on SET_USERNAME
      joinOrCreateRoom() {
        if (!this.connectionAlive) return false
        // CREATE
        if (this.whatAreWeDoing === 'CREATE') {
            if (this.username.length < 1) return false 

            this.socket.send(JSON.stringify({
              type: "CREATE_CHANNEL",
              data: {
                topic: this.question
              }
            }))
            this.$router.push('/room')
        }
        // JOIN
        if (this.whatAreWeDoing === 'JOIN') {

          this.socket.send(JSON.stringify({
            type: "JOIN_CHANNEL"
          }))

        }
      },

      close() {
        this.usernameModal = false
      }

    }
}
</script>



<style lang="scss">
.container {
    margin-top: 60px;
}

.landing_actions {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  .or {
    margin: 0 40px;
    margin-top: 34px;
    text-align: center;
    font-size: 1.5rem;
  }
}
.landing_info {
  width: 600px;
  margin: auto;
  margin-top: 70px;
}
</style>
