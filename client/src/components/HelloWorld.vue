<template>
  <v-container fluid class="h-screen">
    <v-row>
      <v-col cols="6">
        <v-row no-gutters>
          <v-col v-for="t in logTerminals" cols="12">
          <terminal v-bind="t">
            <v-btn
                :disabled="t.active || t.inProgress"
                size="x-small"  variant="plain" @click="start(t)">
              START
            </v-btn>
            <v-btn
                :disabled="!t.active || t.inProgress"
                size="x-small"  variant="plain" @click="stop(t)">
              STOP
            </v-btn>
          </terminal>
      </v-col>
    </v-row>


      </v-col>

      <v-col cols="6">
<!--        <v-col cols="6" class="d-flex flex-column">-->
        <terminal  v-bind="counterLog">
          <v-btn
              :disabled="counterLog.active || counterLog.inProgress"
              size="x-small"  variant="plain" @click="start(counterLog)">
            START
          </v-btn>
          <v-btn
              :disabled="!counterLog.active || counterLog.inProgress"
              size="x-small"  variant="plain" @click="stop(counterLog)">
            STOP
          </v-btn>
        </terminal>

      </v-col>

    </v-row>
  </v-container>
</template>

<script>
import logo from '../assets/logo.svg'
import terminal from "./Terminal.vue";

export default {
  name: 'HelloWorld',
  components: {terminal},
  data: () => ({
    logTerminals: [
      {
        title: 'SCEP',
        alias: 'scepserver',
        height: window.innerHeight * 0.33,
        log: [],
        inProgress: false,
        active: false,
      },
      {
        title: 'NGROK SCEP',
        alias: 'ngrok_scep',
        height: window.innerHeight * 0.14,
        log: [],
        inProgress: false,
        active: false,
      },
      {
        title: 'Nanomdm',
        alias: 'nanomdm',
        height: window.innerHeight * 0.33,
        log: [],
        inProgress: false,
        active: false,
      },
      {
        title: 'NGROK Nanomdm',
        alias: 'ngrok_nanomdm',
        height: window.innerHeight * 0.14,
        log: [],
        inProgress: false,
        active: false,
      },
    ],
    counterLog:{
      title: 'Counter',
      alias: 'counter',
      height: window.innerHeight * 0.25,
      log: [],
      inProgress: false,
      active: false,
      conn: null,
    }
  }),
  methods: {
    start(t){
      console.log("start:", t.title)
      t.active = true
      t.conn = new WebSocket(`ws://localhost:8085/${t.alias}`)
      t.conn.onclose = function(evt) {
        t.log.push('Connection closed')
        this.scrollToBottom(t.alias)
      }

      t.conn.onmessage = (event) => {
        t.log.push(event.data)
        this.scrollToBottom(t.alias)
      }
    },
    stop(t){
      console.log("stop:", t.title)
      if (t.conn != null){
        t.conn.send("stop")
      }
      t.active = false
    },
    scrollToBottom(alias){
      const el = this.$el.querySelector(`#terminal-${alias}`);
      el.scroll({ top: el.scrollHeight, behavior: 'smooth' });
    }
  }
}
</script>
