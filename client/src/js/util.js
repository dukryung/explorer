import moment from 'moment';

export const txType = {
    '/nikto.bankz.v1.MsgTransfer': 'send',
    '/cosmos.bank.v1beta1.MsgSend': 'send',
    '/nikto.convert.v1.MsgConvert': 'convert_token',
    '/nikto.nft.v1.MsgTransferNFT': 'send_nft',
    '/nikto.name.v1.MsgRegisterName': 'register_name',
    '/klaatoo.bankz.v1beta1.MsgIssue': 'issue',
    '/klaatoo.bankz.v1beta1.MsgMint': 'mint',
    '/klaatoo.bankz.v1beta1.MsgBurn': 'burn',
    '/cosmos.slashing.v1beta1.MsgUnjail': 'unjail',
    '/cosmos.staking.v1beta1.MsgBeginRedelegate': 'begin_redelegate',
    '/cosmos.staking.v1beta1.MsgUndelegate': 'undelegate',
    '/cosmos.staking.v1beta1.MsgDelegate': 'delegate',
    '/cosmos.staking.v1beta1.MsgEditValidator': 'edit_validator'
};

export function convertHashToString(base64) {
    const binaryString = window.atob(base64);
    const len = binaryString.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return [...new Uint8Array(bytes.buffer)]
        .map(x => x.toString(16).padStart(2, '0'))
        .join('')
        .toUpperCase();
}

export function truncateHex(hex) {
    const start = hex.substring(0, 6);
    const end = hex.substring(hex.length - 6, hex.length);

    return start + '...' + end;
}

export function getUpdateTime(blockTime) {
    moment.relativeTimeThreshold('ss', 0);
    moment.relativeTimeThreshold('m', 60);
    return moment(new Date(blockTime)).fromNow();
}

export function getVotingPower(totalBondedToken, tokens) {
    return ((tokens / totalBondedToken) * 100).toFixed(2);
}

export function fixedCommisionRate(rates) {
    return `${parseFloat(rates).toFixed(2)}%`;
}
