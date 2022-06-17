<template>
  <div class="mainContainer">
    <terminal
        v-for="t in store.logTerminals"
        :cmd="t"
        :class="t.alias"
    ></terminal>
    <terminal
        :cmd="store.commonLog"
        class="general-log"
    ></terminal>
    <steps class="steps-container"></steps>
  </div>
</template>

<script setup>
import Terminal from "./components/Terminal.vue";
import Steps from "./components/Steps.vue";
import Cmd, {Status} from "./models/cmd";
import {onBeforeMount, reactive} from 'vue'
import {useCommandsStore} from "./stores/commands";

const store = useCommandsStore()

const counterLog = reactive(new Cmd(
    0,
    'Counter',
    'counter',
    './counter',
    [],
    Status.INACTIVE,
    true))

onBeforeMount(()=> { store.initialize() })

</script>
<style>
#app {
  font-family: 'Courier New', monospace, Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.mainContainer {
  display: grid;
  grid-auto-rows: minmax(0, 1fr);
  grid-auto-columns: minmax(0, 1fr);
  grid-template-areas:
    'scep gl'
    'scep st'
    'ngr1 st'
    'nano st'
    'nano st'
    'ngr2 st';
  gap: 10px;
  /*width: 98vw;*/
  height: 96vh;
}
.scepserver {
  grid-area: scep;
}
.ngrok_scep {
  grid-area: ngr1;
}
.nanomdm {
  grid-area: nano;
}
.ngrok_nanomdm {
  grid-area: ngr2;
}
.general-log {
  grid-area: gl;
}
.steps-container {
  grid-area: st;
  height: available;
}
.green {
  color: green;
}
.btns {
  display: flex;
  justify-content: space-evenly;
}

.btn {
  border: 1px solid black;
  border-radius: 4px;
  background-color: white;
  color: black;
  cursor: pointer;
  width: 60px;
  transition: 0.3s;
}
.btn:hover {
  box-shadow: 0 2px 4px 0 rgba(0,0,0,0.2);
}
</style>
