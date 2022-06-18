<template>
  <div class="container">
    <div class="btns-all">
      <button class="btn large" @click="$emit('hide')" style="background-color: rgba(255,0,19,0.1)"><span style="float: left;">↩</span>Cancel</button>
      <button class="btn large" @click="save" style="background-color: rgba(187,251,123,0.37)"><span style="float: left;">💾</span>Save</button>
      <button class="btn large" @click="loadDefaults" style="background-color: rgba(245,217,6,0.37)"><span style="float: left;">🌀</span>Defaults</button>
    </div>
    <div class="one-cmd" v-for="cmd in state.localAllCmd">
      <input style="width: 100%;background-color: lightgrey" v-model="cmd.title" placeholder="Command title" />
      <textarea style="width: 100%;border: none;" v-model="cmd.text" placeholder="command"></textarea>
    </div>

  </div>
</template>

<script setup>
import {useCommandsStore} from "../stores/commands";
import {onMounted, reactive} from "vue";
import {allCmdDefault} from "../models/cmd_defaults";

const store = useCommandsStore()

const state = reactive({localAllCmd: []})

const emit = defineEmits([ 'hide'])

function save() {
  for(let i=0;i<store.allCmd.length;i++){
    store.allCmd[i].title = state.localAllCmd[i].title
    store.allCmd[i].text = state.localAllCmd[i].text
  }
  store.saveToLocalStorage()
  emit('hide')
}
function loadDefaults() {
    for(let i=0;i<allCmdDefault.length;i++){
      state.localAllCmd[i].title = allCmdDefault[i].title
      state.localAllCmd[i].text = allCmdDefault[i].text
    }
}
onMounted(() => {
  // no need class methods, simple deep copy as object
  // but logs can be huge
  //state.localAllCmd = JSON.parse(JSON.stringify(store.allCmd))

  // this more effective
  for (const cmd of store.allCmd) {
    state.localAllCmd.push(
        {
          "title": cmd.title,
          "text": cmd.text,
        }
    )
  }
})
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
  overflow-x: auto;
}

.one-cmd {
  padding-left: 0;
  padding-right: 10px;
  box-shadow: 0 3px 6px 0 rgba(0, 0, 0, 0.4);
  transition: 0.3s;
}

.one-cmd:hover {
  box-shadow: 0 6px 12px 0 rgba(0, 0, 0, 0.4);
}

.btns-all {
  display: flex;
  justify-content: end;
}
</style>