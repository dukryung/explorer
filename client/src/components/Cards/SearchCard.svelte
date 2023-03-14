<script>
    import { navigate } from 'svelte-routing';
    import { _ } from 'svelte-i18n';

    export let subTitle = '';
    export let title = '';

    let input = '';
    let replace = false;

    function onSearchEvent() {
        replace = false;

        if (input === '' || !input) {
            return;
        }

        const height = getBlockHeight(input);
        console.log('height: ', height);
        if (height !== undefined) {
            console.log('height: inside');
            if (location.pathname.startsWith('/block')) {
                location = `/block/${height}`;
            } else {
                navigate(`/block/${height}`, { replace: replace });
            }
            return;
        }

        const addr = getAddress(input);
        if (addr !== undefined) {
            if (location.pathname.startsWith('/account')) {
                location = `/account/${addr}`;
            } else {
                navigate(`/account/${addr}`, { replace: replace });
            }
            return;
        }

        const nodeAddr = getValAddress(input);
        if (nodeAddr !== undefined) {
            if (location.pathname.startsWith('/validator')) {
                location = `/validator/${nodeAddr}`;
            } else {
                navigate(`/validator/${nodeAddr}`, { replace: replace });
            }
            return;
        }

        const txHash = getTxHash(input);
        if (txHash !== undefined) {
            if (location.pathname.startsWith('/tx')) {
                location = `/tx/${txHash}`;
            } else {
                navigate(`/tx/${txHash}`, { replace: replace });
            }
            return;
        }

        const nftId = getNftId(input);
        console.log('nftId:', nftId);
        if (nftId !== undefined) {
            if (location.pathname.startsWith('/nft/')) {
                location = `/nft/${nftId}/`;
            } else {
                navigate(`/nft/${nftId}/`, {
                    replace: replace,
                    state: { id: nftId }
                });
            }
            return;
        }
    }

    function getBlockHeight(value) {
        console.log('getBlockHeight', value);
        value = value.replace(/^(b|bl|block|h|hei|height):/, '');
        console.log('getBlockHeight', value);
        const parsed = parseInt(value);
        if (isNaN(parsed)) {
            return undefined;
        }
        return parsed;
    }

    function getAddress(value) {
        value = value.replace(/^(a|addr|address|acc|accnt|account):/, '');
        console.log(value, network.bech32.account);
        if (value.startsWith(`${network.bech32.account}1`)) {
            return value;
        }
        return undefined;
    }

    function getValAddress(value) {
        value = value.replace(/^(n|node|v|val|validator)/, '');
        if (
            value.startsWith(
                `${network.bech32.account}${network.bech32.validator}${network.bech32.operator}1`
            )
        ) {
            return value;
        }
        return undefined;
    }

    function getTxHash(value) {
        value = value.replace(/^(t|tx|hash)/, '');
        if (value.length !== 64) {
            return undefined;
        }
        if (/^([a-fA-F0-9]+)^/.test(value)) {
            return value;
        }
        return undefined;
    }

    function getNftId(value) {
        value = value.replace(/^(nft|nftid|nft_id):/, '');
        if (/^([a-zA-Z0-9_.+-]+)$/.test(value)) {
            return value;
        }
        return undefined;
    }

    function onKeyPress(e) {
        if (e.charCode === 13) {
            onSearchEvent();
        }
    }
</script>

<section class="mainSearch">
    <div class="mainWrap">
        <h1 class="white"><span class="welcomeTxt">{subTitle}</span>{title}</h1>
        <div class="searchBar">
            <div id="search">
                <input
                    id="input"
                    placeholder={$_('search_placeholder')}
                    on:keypress={onKeyPress}
                    bind:value={input}
                />
                <button id="button" on:click={onSearchEvent}>
                    <i class="fa fa-search" />
                </button>
            </div>
        </div>
    </div>
</section>
