<template>
  <div class="terminal" >
    <div class="header">
      <div>{{ title }}</div>
      <div class="btns">
        <button
            :disabled="!showStart"
            class="btn"
            @click="$emit('startClick')"
        >start
        </button>
        <button
            :disabled="!showStop"
            class="btn"
            @click="$emit('stopClick')"
        >stop
        </button>
<!--        <button class="btn" @click="scrollToBottom">scroll</button>-->
      </div>
    </div>
    <div class="content" :id="`terminal-${alias}`">
      <div class="log-el" v-for="l in log" v-text="l"></div>
    </div>
  </div>
</template>

<script>
import {Status} from "../models/cmd";

export default {
  props: {
    title: String,
    alias: String,
    log: [],
    status: Symbol,
  },
  data(){
    return{
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
    logLen(){ return this.log.length},
    showStart(){ return this.status != Status.ACTIVE && this.status != Status.IN_PROGRESS },
    showStop(){ return this.status == Status.ACTIVE }
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
    this.elToScroll = window.document.querySelector(`#terminal-${this.alias}`);
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

.btns {
  display: flex;
  justify-content: space-evenly;
}

.btn {
  margin-right: 5px;
}

.content {
  overflow: auto;
  padding-left: 5px;
}

</style>