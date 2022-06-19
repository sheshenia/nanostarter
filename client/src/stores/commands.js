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
        enrollmentIDs: new Set(),
        saveLogsToLocalstorage: false,

    }),
    getters: {
        logTerminals: (state) => state.allCmd.filter(cmd => cmd.ifTerminal),
    },
    actions: {
        initialize() {
            const sltl = localStorage.getItem("saveLogsToLocalstorage")
            if(sltl !== null && sltl==="true"){
                this.saveLogsToLocalstorage = true
                // TODO add improvised watcher to save logs

            }

            if(this.saveLogsToLocalstorage){
                const commonLogsStr = localStorage.getItem("commonLog")
                if (commonLogsStr !== null){
                    const cl = JSON.parse(commonLogsStr)
                    if(Array.isArray(cl)){
                        this.commonLog.log = cl
                    }
                }
            }

            //TODO try to load configs from localstorage
            let allCmdArrayObj =  allCmdDefault

            // try to get all commands from localstorage and deep clone to store allCmd
            // if smth. wrong populate allCmd with default commands
            try {
                const allCmdString = localStorage.getItem("allCmd")
                if(allCmdString !== null){
                    allCmdArrayObj = JSON.parse(allCmdString)
                }
                this.allCmdDeepClone(allCmdArrayObj)
            }catch (e) {
                console.log(e)
                this.allCmdDeepClone(allCmdDefault)
            }
        },
        saveToLocalStorage(){
            localStorage.setItem("saveLogsToLocalstorage", this.saveLogsToLocalstorage || "")

            //only save title, text and logs (if saveLogsToLocalstorage)
            const localAllCmd = []
            for(let i=0;i<this.allCmd.length;i++) {
                localAllCmd.push({
                    "title": this.allCmd[i].title,
                    "text": this.allCmd[i].text,
                    "log": this.saveLogsToLocalstorage ? this.allCmd[i].log : []
                })
            }
            localStorage.setItem("allCmd", JSON.stringify(localAllCmd))
            localStorage.setItem("commonLog", JSON.stringify(this.saveLogsToLocalstorage ? this.commonLog.log : []))
        },
        // allCmdDeepClone deep clone with class assign
        allCmdDeepClone(allCmdArrayObj){
            this.allCmd = []
            // need deep cloning with class
            // JSON.parse(JSON.stringify(allCmdDefault)); doesn't work we are loosing class Cmd main functionality
            //this.allCmd = JSON.parse(JSON.stringify(allCmdDefault));
            for(let i=0;i<allCmdArrayObj.length;i++) {
                this.allCmd.push(new Cmd(
                    allCmdDefault[i].order,
                    allCmdArrayObj[i].title,
                    allCmdDefault[i].alias,
                    allCmdArrayObj[i].text,
                    allCmdArrayObj[i].log,
                    allCmdDefault[i].status,
                    allCmdDefault[i].ifTerminal,
                    allCmdDefault[i].stepAction, //function
                    null
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
            this.notifications = {}
            this.enrollmentIDs.clear()
        },
        clearLogs(){
            for(const cmd of this.allCmd){
                cmd.clearLog()
            }
            this.commonLog.clearLog()
        },
        activateEnrollDeviceStep(){
            this.allCmd.find(cmd=>cmd.alias==="enroll_device").activate()
        }
    }
})

function delay(time) {
    return new Promise(resolve => setTimeout(resolve, time));
}