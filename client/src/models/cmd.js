export default class Cmd {
    constructor(order, title, alias, log, status=Status.INACTIVE, ifTerminal = false, conn = null) {
        this.order = order
        this.title = title
        this.alias = alias
        this.log = log
        this.status = status===false ? Status.INACTIVE : status
        this.ifTerminal = ifTerminal
        this.conn = conn
    }
    activate() { this.status = Status.ACTIVE }
    deactivate() { this.status = Status.INACTIVE }
    inProgress() { this.status = Status.IN_PROGRESS }
}

export const Status = Object.freeze({
    INACTIVE: Symbol("inactive"),
    IN_PROGRESS: Symbol("in_progress"),
    ACTIVE: Symbol("active"),
    ERROR: Symbol("error")
})