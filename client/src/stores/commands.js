import { defineStore } from 'pinia';
import Cmd, {Status} from "../models/cmd";
import {allCmdDefault} from "../models/cmd_defaults";

export const useCommandsStore = defineStore('commands', {
    state: () => ({
        commonLog: new Cmd(
            0,
            'Common Log',
            'command',
            '',
            [],
            Status.INACTIVE,
            true),
        allCmd: [],
        notifications: {
            //'ngrok_scep': "some url"
        },

    }),
    getters: {
        logTerminals: (state) => state.allCmd.filter(cmd => cmd.ifTerminal),
    },
    actions: {
        initialize() {
            // need deep cloning with class
            // JSON.parse(JSON.stringify(allCmdDefault)); doesn't work we are loosing class Cmd main functionality
            //this.allCmd = JSON.parse(JSON.stringify(allCmdDefault));
            this.allCmd = []
            for(const cmd of allCmdDefault) {
                this.allCmd.push(new Cmd(
                    cmd.order,
                    cmd.title,
                    cmd.alias,
                    cmd.text,
                    cmd.log,
                    cmd.status,
                    cmd.ifTerminal,
                    cmd.stepAction,
                    cmd.conn
                ))
            }
        },
        async execCommonLogCmd(cmd){
            this.commonLog.text = cmd.text
            this.commonLog.log.push(cmd.title)
            try {
                await this.commonLog.startCmd(window.__API_ENDPOINT__)
            }catch (e) {
                console.log(e)
            }finally {
                cmd.status = this.commonLog.status
                this.commonLog.deactivate()
                this.commonLog.text = ""
            }
        },
        pushToCommonLog(msg){
            this.commonLog.log.push(msg)
        },
        async startAll(){
            console.log("startAll click")
            for(const cmd of this.allCmd){
                if(cmd.isActive){ continue }
                if(cmd.stepAction != null){
                    cmd.stepAction(this)
                }else{
                    cmd.startWS(window.__WEBSOCKET_ENDPOINT__)
                }
                await delay(1500);
            }
        },
        stopAll(){
            console.log("stopAll click")
            for(let i = this.allCmd.length -1;i>=0;i--){
                this.allCmd[i].stopWS()
            }
        },
        clearLogs(){
            for(const cmd of this.allCmd){
                cmd.clearLog()
            }
            this.commonLog.clearLog()
        },
    }
})

function delay(time) {
    return new Promise(resolve => setTimeout(resolve, time));
}