import Cmd, {Status} from "./cmd";
import {generateProfile} from "./enroll.mobileconfig";

export const allCmdDefault = [
    new Cmd(
        1,
        'SCEP',
        'scepserver',
        './scep/scepserver-linux-amd64 -allowrenew 0 -challenge nanomdm -debug',
        [],
        Status.INACTIVE,
        true,
        ),
    new Cmd(
        2,
        'NGROK SCEP',
        'ngrok_scep',
        'ngrok http 8080 --log=stdout',
        [],
        Status.INACTIVE,
        true,
        ),
    new Cmd(
        3,
        'Retrieve SCEP certificate',
        'get_scep_cert',
        `curl 'https://some_address.ngrok.io/scep?operation=GetCACert' | openssl x509 -inform DER > ./nanomdm/ca.pem`,
        [],
        Status.INACTIVE,
        false,
        function(store) {
            console.log("current:", this.text)
            if (!store.notifications["ngrok_scep"]){
                console.log("no ngrok_scep binding URL, start ngrok at first")
                return
            }
            this.text = this.text.replace(/https\S+\.io/i, store.notifications["ngrok_scep"])
            console.log("updated:", this.text)
            store.execCommonLogCmd(this)
        }
        ),
    new Cmd(
        4,
        'Nanomdm',
        'nanomdm',
        './nanomdm/nanomdm-linux-amd64 -ca ./ca.pem -api nanomdm -debug',
        [],
        Status.INACTIVE,
        true),
    new Cmd(
        5,
        'NGROK Nanomdm',
        'ngrok_nanomdm',
        'ngrok http 9000 --log=stdout',
        [],
        Status.INACTIVE,
        true),
    new Cmd(
        6,
        'Push certificate',
        'push_cert',
        `cat ./nanomdm/certs/*.pem ./nanomdm/certs/*.key | curl -T - -u nanomdm:nanomdm 'http://127.0.0.1:9000/v1/pushcert'`,
        [],
        Status.INACTIVE,
        false,
        function (store) {
            store.execCommonLogCmd(this)
        }),
    new Cmd(
        7,
        'Config profile xml file',
        'profile_file',
        '',
        [],
        Status.INACTIVE,
        false,
        function (store) {
            const profileErr = "Config profile error: "
            if (!store.notifications['ngrok_scep']){
                store.pushToCommonLog(profileErr + 'empty ngrok_scep URL')
                return
            }
            if (!store.notifications['ngrok_nanomdm']){
                store.pushToCommonLog(profileErr + 'empty ngrok_nanomdm URL')
                return
            }
            if (!store.notifications['nanomdm']){
                store.pushToCommonLog(profileErr + 'empty device ID')
                return
            }
            const pf = generateProfile(
                store.notifications['ngrok_scep'],
                store.notifications['ngrok_nanomdm'],
                store.notifications['nanomdm'],
            )
            downloadProfile("enroll.mobileconfig", pf)
            this.activate()
        }),
    /*new Cmd(
        8,
        'Download or send by email profile file',
        'email_profile_file',
        '',
        [],
        Status.INACTIVE),*/
    new Cmd(
        8,
        'Enroll your device',
        'enroll_device',
        '',
        [],
        Status.INACTIVE,
        false,
        function () {}),
]

function downloadProfile(fileName, content) {
    const el = document.createElement('a');
    el.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(content));
    el.setAttribute('download', fileName);
    el.style.display = 'none';
    document.body.appendChild(el);
    el.click();
    document.body.removeChild(el);
}