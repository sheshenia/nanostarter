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
        }
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
        }
    }
})

