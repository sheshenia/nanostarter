import { defineStore } from 'pinia';
import Cmd, {Status} from "../models/cmd";

export const useCommandsStore = defineStore('commands', {
    state: () => ({
        commonLog: new Cmd(
            0,
            'Common Log',
            '',
            '',
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
        ]
    })
})