<template>
  <div>
    <div
        v-for="cmd in store.allCmd"
        class="step"
        :class="{'step-active': cmd.isActive}"
    >
      <div>
        <div :class="cmd.isInProgress ? 'loader':'circle'">{{ cmd.isActive ? '':cmd.order }}</div>
      </div>
      <div>
        <div class="title">
          <button
              v-if="cmd.showStart"
              class="btn"
              @click="cmd.startWS(wsEndpoint)"
          ><span class="green">▶</span> start
          </button>
          <button
              v-if="cmd.showStop"
              class="btn"
              @click="cmd.stopWS()"
          >🟥 stop
          </button>
          {{cmd.title}}
        </div>
        <div class="caption">{{ cmd.text }}</div>
      </div>
    </div>

    <!--<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAABQVBMVEUAAAAA/wCA/wCA/4BVqgBAv0BV1VVJ2yQuuRdd0UZq6lVi2E5b2zdRyTZZ2UBo3FFAvxVi4k5SzjFk4E1a1j5m301Nyipl4U1Avxdj4ktl30pX1DpIyCZm4U1k4ktQyzBHwyJV0ThFwx9U0DVd10FBvhpJxiRl4E5Z1TtRzDE+uxRk3kw2sQY1sQZb10Fg20NOySti3Uhh3Uhg2kM5tgxNyite2kRRzDBh3Eg5tg1k30xb10Bf20RV0DZZ1Ds+uxRl4E5V0TZl301c10BCvxpl4E01sQdRzTBl4E1Z1Ds1sgdNyitl4E1GxCBNyitV0DZZ1Dtl4E1KxiVRzS9KxiVl4E1GwyBRzTBV0DZCvxpc10Bl4E05tg0+uxRCvxpGwyBKxiVNyitRzTBV0DZZ1Dtc10Bf2kRh3Ehk30xl4E3///+oZ9JEAAAAXHRSTlMAAQICAwQGBwsLDA0OExQWGBofISUoKyssLDA1PDw9QEROUVJSYmJiZ250dXZ9foWKjY6Qk6ClpqmvuLvExcjKzNba2+Dg4+Tl5+nt7vDx9PT19vf6+/z8/P7+/uCrIccAAAABYktHRGolYpUOAAAA2klEQVQYGd3BiToCUQAG0F9TQ6stlCJLISRFKEtRthRCCpW2mdz3fwHdO98dqpkeoHMw4qYtGGq5cD2HIdw5QnJL0DV7R7oKC9AxeUWo4iK0WVOE+tqENjHxw+yAEyJRESrhpMNEoArL8pkD3IHMHAng/J+SJKWdUOxJzKkIzvbcpu49oLbazIUZKtN5i3laBbBeblGXdvwzftxkSkH4XpvU7Qx6GPYbithjg3qYR7+Nt/qfFy8GreS/ufc1aHFlaoqPALRNJavMNvRMxCtdh9BnDGVvdscwin4B+iFCWM0L0y4AAAAASUVORK5CYII=">-->
  </div>
</template>

<script setup>
import {useCommandsStore } from "../stores/commands";

defineProps({
  allCmd: {},
})

const store = useCommandsStore()

const wsEndpoint = window.__WEBSOCKET_ENDPOINT__

</script>

<style scoped>
.step {
  position: relative;
  min-height: 1em;
  color: gray;
}
.step + .step {
  margin-top: 1.5em
}
.step > div:first-child {
  position: static;
  height: 0;
}
.step > div:not(:first-child) {
  margin-left: 1.5em;
  padding-left: 1em;
}
.step.step-active {
  color: #4285f4
}
/*.step.step-active .circle {
  background-color: #4285f4;
}*/

.step.step-active .circle {
  background-color: white;
  background-image: url("../assets/check-32.png") ;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}

/* Circle */
.circle {
  background: gray;
  position: relative;
  width: 1.5em;
  height: 1.5em;
  line-height: 1.5em;
  border-radius: 100%;
  color: #fff;
  text-align: center;
  box-shadow: 0 0 0 3px #fff;
}

/* Vertical Line */
.circle:after {
  content: ' ';
  position: absolute;
  display: block;
  top: 1px;
  right: 50%;
  bottom: 1px;
  left: 50%;
  height: 100%;
  width: 1px;
  transform: scale(1, 2);
  transform-origin: 50% -100%;
  background-color: rgba(0, 0, 0, 0.25);
  z-index: -1;
}
.step:last-child .circle:after {
  display: none
}

/* Stepper Titles */
.title {
  line-height: 1.5em;
  font-weight: bold;
}
.caption {
  margin-left: 65px;
  font-size: 0.8em;
  /*text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;*/
}

.loader {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  width: 18px;
  height: 18px;
  animation: spin 2s linear infinite;
  background-color: white;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>