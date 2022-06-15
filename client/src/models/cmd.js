export default class Cmd {
    constructor(order, title, alias, text, log, status=Status.INACTIVE, ifTerminal = false, conn = null) {
        this.order = order
        this.title = title
        this.alias = alias
        this.text = text
        this.log = log
        this.status = status===false ? Status.INACTIVE : status
        this.ifTerminal = ifTerminal
        this.conn = conn
    }

    get showStart() { return this.status != Status.ACTIVE && this.status != Status.IN_PROGRESS }
    get showStop() { return this.status == Status.ACTIVE }
    get isActive() { return this.status == Status.ACTIVE }
    get isInProgress() { return this.status == Status.IN_PROGRESS }

    activate() { this.status = Status.ACTIVE }
    deactivate() { this.status = Status.INACTIVE }
    inProgress() { this.status = Status.IN_PROGRESS }

    // startWS starts websocket connection
    startWS(websocketEndpoint) {
        console.log("start:", this.title)
        this.inProgress()
        this.conn = new WebSocket(websocketEndpoint + this.alias)
        this.conn.onclose = (_) => {
            this.log.push('Connection closed')
            this.deactivate()
        }
        this.conn.onmessage = (event) => {
            this.log.push(event.data)
            this.activate()
        }
    }
    // stopWS try to close websocket connection, or send "stop" cmd to the server
    stopWS(){
        console.log("stop:", this.title)
        if (this.conn == null){
            this.deactivate()
            return
        }
        // if nothing to send yet closing the connection
        if (this.conn.bufferedAmount==0){
            this.conn.close()
            return
        }
        // otherwise send "stop" command to the server
        this.conn.send("stop")
    }
}

export const Status = Object.freeze({
    INACTIVE: Symbol("inactive"),
    IN_PROGRESS: Symbol("in_progress"),
    ACTIVE: Symbol("active"),
    ERROR: Symbol("error")
})