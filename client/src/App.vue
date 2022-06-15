<template>
  <div class="mainContainer">
    <terminal
        v-for="t in logTerminals"
        :cmd="t"
        :class="t.alias"
    ></terminal>
    <terminal
        :cmd="counterLog"
        class="general-log"
    ></terminal>
    <steps class="steps-container" v-bind:allCmd="allCmd"></steps>
  </div>
</template>

<script>
import Terminal from "./components/Terminal.vue";
import Cmd, {Status} from "./models/cmd";
import Steps from "./components/Steps.vue";

export default {
  components: {
    Steps,
    Terminal,
  },
  data() {
    return {
      counterLog: new Cmd(
          0,
          'Counter',
          'counter',
          './counter',
          [],
          Status.INACTIVE,
          true),

      allCmd: [
        new Cmd(
            1,
            'SCEP',
            'scepserver',
            './scep/scepserver-linux-amd64 -allowrenew 0 -challenge nanomdm -debug',
            [],
            Status.INACTIVE,
            true),
        new Cmd(
            2,
            'NGROK SCEP',
            'ngrok_scep',
            'ngrok http 8080 --log=stdout',
            [],
            Status.INACTIVE,
            true),
        new Cmd(
            3,
            'Retrieve SCEP certificate',
            'get_scep_cert',
            `curl 'https://d53c350ae070.eu.ngrok.io/scep?operation=GetCACert' | openssl x509 -inform DER > ca.pem`,
            [],
            Status.INACTIVE,
            false),
        new Cmd(
            4,
            'Nanomdm',
            'nanomdm',
            './nanomdm-linux-amd64  -ca ca.pem -api nanomdm -debug',
            [],
            Status.INACTIVE,
            true),
        new Cmd(
            5,
            'NGROK Nanomdm',
            'ngrok_nanomdm',
            'ngrok http 9000',
            [],
            Status.INACTIVE,
            true),
        new Cmd(
            6,
            'Push certificate',
            'push_cert',
            `cat ./certs/certs/*.pem ./certs/certs/*.key | curl -T - -u nanomdm:nanomdm 'http://127.0.0.1:9000/v1/pushcert'`,
            [],
            Status.INACTIVE,
            false),
        new Cmd(
            7,
            'Config profile xml file',
            'profile_file',
            '',
            [],
            Status.INACTIVE,
            false),
        new Cmd(
            8,
            'Download or send by email profile file',
            'email_profile_file',
            '',
            [],
            Status.INACTIVE),
        new Cmd(
            9,
            'Enroll your device',
            'enroll_device',
            '',
            [],
            Status.INACTIVE,
            false),
      ],
    }
  },
  computed:{
    logTerminals(){
      return this.allCmd.filter(cmd => cmd.ifTerminal)
    }
  },
  methods: {
  },
  mounted() {
  }
}
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
}

</style>
