<!--<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  title: String,
  alias: String,
  height: String,
  log: Array,
})

// declare a ref to hold the element reference
// the name must match template ref value
const [alias] = ref(null)

onMounted(() => {
  [alias].scroll({ top: [alias].scrollHeight, behavior: 'smooth' });
})


</script>-->

<template>
  <v-card
      class="mb-2 overflow-y-auto"
      theme="dark"
      :height="height"
      :id="`terminal-${alias}`"
  >
    <v-banner sticky class="bg-black pt-1" style="height: 25px;">
      <v-card-subtitle>{{ title }}</v-card-subtitle>
      <v-spacer></v-spacer>
      <v-card-subtitle class="float-right">
        <slot></slot>
      </v-card-subtitle>
    </v-banner>
<!--    <v-card-header sticky  class="bg-black pt-4" style="height: 25px;">
      <v-card-subtitle>{{ title }}</v-card-subtitle>
      <v-spacer></v-spacer>
      <v-card-subtitle class="float-right">
        <slot></slot>
      </v-card-subtitle>
    </v-card-header>-->

    <v-card-text class="pt-1">
      <p v-for="l in log" v-text="l"></p>
    </v-card-text>

  </v-card>
</template>

<script>
export default {
  props: {
    title: String,
    alias: String,
    height: Number,
    log: Array,
  }
}
</script>