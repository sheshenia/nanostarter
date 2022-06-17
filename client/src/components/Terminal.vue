<template>
  <div class="terminal" >
    <div class="header">
      <div>{{ cmd.title }}</div>
      <div>{{ notification }}</div>
      <div class="btns">
        <button
            v-if="cmd.showStart"
            class="terminal-btn"
            @click="cmd.startWS(wsEndpoint)"
            :title="`Start ${cmd.alias}`"
        ><span class="green">▶</span> start
        </button>
        <button
            v-if="cmd.showStop"
            class="terminal-btn"
            @click="cmd.stopWS()"
            :title="`Stop ${cmd.alias}`"
        >🟥 stop
        </button>
        <button class="terminal-btn" @click="clearLog" title="Clear Log screen">🗑</button>
<!--        <button class="btn" @click="scrollToBottom">scroll</button>-->
      </div>
    </div>
    <div class="content" :id="`terminal-${cmd.alias}`">
      <div class="log-el" v-for="l in cmd.log" v-text="l"></div>
    </div>
  </div>
</template>

<script>
import Cmd from "../models/cmd";
import {useCommandsStore} from "../stores/commands";
import {mapStores} from "pinia/dist/pinia";

export default {
  props: {
    cmd: Cmd,
  },
  data(){
    return{
      wsEndpoint: window.__WEBSOCKET_ENDPOINT__,
      elToScroll: null,
    }
  },
  methods: {
    // scrollToBottom scrolls to the bottom of div element with logs
    scrollToBottom(){
      // TODO: find out why ref doesn't work, till than use native DOM
      if (this.elToScroll != null){
        this.elToScroll.scroll({ top: this.elToScroll.scrollHeight, behavior: "smooth" });
      }
    },
    clearLog() {
      console.log("clearLog click")
      this.cmd.clearLog()
    },
    addWatchersToNgrok(){
      // watcher for detection ngrok URL address
      this.$watch('logLen', ()=>{
        const ngUrl = /https\S+\.io/i.exec(this.cmd.log[this.cmd.log.length-1])
        if (ngUrl != null){
          console.log("Founded ngrok URL:", ngUrl[0])
          this.commandsStore.pushToCommonLog(`Got ${this.cmd.alias} URL: ${ngUrl[0]}`)
          this.commandsStore.notifications[this.cmd.alias] = ngUrl[0]
        }
      })

      //on change status, clean up ngrok url, and need to stop all other dependent step processes
      this.$watch('cmd.status', ()=>{
        if (!this.cmd.isActive){
          this.commandsStore.notifications[this.cmd.alias] = ""
        }
      })
    }
  },
  computed: {
    // creating computed value for watching changes, and scroll down
    logLen(){ return this.cmd.log.length},

    // pinia uses store's: id + Store
    ...mapStores(useCommandsStore),
    notification() { return this.commandsStore.notifications[this.cmd.alias]}
  },
  watch: {
    // watching the length of the log, if changed - scroll down
    // nextTick because if not watcher fires earlier than vue creates dom element
    logLen(val){
      this.$nextTick(() => {
        this.scrollToBottom()
      });
    }
  },
  mounted() {
    // getting dom element after mounted for scrolling
    this.elToScroll = window.document.querySelector(`#terminal-${this.cmd.alias}`);
    console.log("Cmd terminal mounted:", this.cmd.alias || 'no alias')

    // dynamically adding watchers for some terminals like ngrok
    switch (this.cmd.alias) {
      case 'ngrok_scep'  :
      case 'ngrok_nanomdm':
        console.log("add watcher to", this.cmd.alias)
        this.addWatchersToNgrok()
        break
      case 'nanomdm':
        // watcher for detection device id topic=com.apple.mgmt.External.e3b8ceac-1f18-2c8e-8a63-dd17d99435d9
        // e3b8ceac-1f18-2c8e-8a63-dd17d99435d9
        this.$watch('logLen', ()=>{
          const deviceIdReg = /topic=com.apple[\w\.]+\.([a-z0-9-]+)/i
          const lastLog = this.cmd.log[this.cmd.log.length-1]
          const match = lastLog.match(deviceIdReg)
          if (match != null){
            console.log("Founded nano device ID:", match[1])
            this.commandsStore.pushToCommonLog(`Got ${this.cmd.alias} device ID: ${match[1]}`)
            this.commandsStore.notifications[this.cmd.alias] = match[1]
          }
        })
        break
    }

  }
}
</script>

<style scoped>
  .terminal {
    display: flex;
    flex-direction: column;
    background-color: #2e2d2d;
    color: whitesmoke;
    border-radius: 5px;
    box-shadow: 0 3px 6px 0 rgba(0,0,0,0.4);
    transition: 0.3s;
  }
  .terminal:hover {
    box-shadow: 0 6px 12px 0 rgba(0,0,0,0.4);
  }

  .header {
    height: 20px;
    background-color: black;
    padding: .2rem .5rem;
    display: flex;
    justify-content: space-between;
    border-radius: 5px;
  }

  .content {
    overflow: auto;
    padding-left: 5px;
  }
  .terminal-btn {
    border: none;
    border-radius: 4px;
    background-color: inherit;
    color: lightgrey;
    cursor: pointer;
    margin-right: 5px;
  }
  /* On mouse-over */
  .terminal-btn:hover {background: #322f2f;}
</style>