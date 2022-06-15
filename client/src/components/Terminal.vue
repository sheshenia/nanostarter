<template>
  <div class="terminal" >
    <div class="header">
      <div>{{ cmd.title }}</div>
      <div class="btns">
        <button
            v-if="cmd.showStart"
            class="terminal-btn"
            @click="cmd.startWS(wsEndpoint)"
        ><span class="green">▶</span> start
        </button>
        <button
            v-if="cmd.showStop"
            class="terminal-btn"
            @click="cmd.stopWS()"
        >🟥 stop
        </button>
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
    }
  },
  computed: {
    // creating computed value for watching changes, and scroll down
    logLen(){ return this.cmd.log.length},
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
    console.log("Cmd prop:", this.cmd)
    console.log("Cmd alias:", this.cmd.alias)

    // dynamically adding watchers for some terminals like ngrok
    switch (this.cmd.alias) {
      case 'ngrok_scep':
        console.log("added watcher for ngrok_scep")
        this.$watch('logLen', ()=>{
          const ngUrl = /https\S+\.io/i.exec(this.cmd.log[this.cmd.log.length-1])
          if (ngUrl != null){
            console.log("Founded ngrok URL:", ngUrl[0])
          }

        })
        break
    }
    if (this.cmd.alias == "counter"){
      console.log("init watcher for counter")
      this.$watch('logLen', ()=>{
        console.log(this.cmd.log[this.cmd.log.length-1])
        if (this.cmd.log[this.cmd.log.length-1].includes("5")){
          console.log("Number 5 triggered")
        }

      })
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

}

.header {
  height: 20px;
  background-color: black;
  padding: .2rem .5rem;
  display: flex;
  justify-content: space-between;
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