<script>
    import { onDestroy, onMount } from 'svelte';
    import { _ } from 'svelte-i18n';
    import { Link, navigate } from 'svelte-routing';
    import { getUpdateTime, truncateHex, txType } from 'js';
    import TxResultIcon from '../Labels/TxResultIcon.svelte';
    import IconRightArrow from '../../assets/IconRightArrow.svelte';
    import { Token } from '../../js/token';
    import NoData from 'Components/Labels/NoData';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import Paginator from 'Components/Footers/Paginator';

    // default value
    export let txLimit = 10;
    export let txPage = 0;
    export let address = '';
    export let height = 0;

    // default params
    export let limit = false;
    export let type = klaatoo.SUBSCRIBE;
    export let method = 'tx.latest';

    const requestId = klaatoo.generateRequestId();
    let txs = [];
    let totalTxs = 0;
    let onLoad = false;

    onMount(() => {});

    $: if (txPage !== undefined) {
        window.scroll({
            top: 0,
            left: 0,
            behavior: 'smooth'
        });
        let params;
        if (method === 'tx.latest') {
            params = [txLimit, txPage];
        } else if (method === 'tx.address') {
            params = [address, txLimit, txPage];
        } else if (method === 'tx.height') {
            params = [height, txLimit, txPage];
        }
        klaatoo.request({
            type: type,
            method: method,
            params: params,
            id: requestId,
            success: data => {
                if (data !== null) {
                    if (data.txs !== null) txs = data.txs;
                    totalTxs = data.total;
                }
                onLoad = true;
            },
            error: error => {
                console.error(error);
                navigate('/error', { replace: false });
            }
        });
    }

    onDestroy(() => {
        if (type === klaatoo.SUBSCRIBE) {
            klaatoo.unsubscribe({
                method: method,
                id: requestId
            });
        }
    });

    function getSendAmount(messages) {
        if (messages[0]['@type'] === '/nikto.bankz.v1.MsgTransfer') {
            return `${new Token({
                symbol: network.networkToken.symbol,
                amount: messages[0].coin.amount,
                precision: network.networkToken.precision
            }).string()}`;
        } else if (messages[0]['@type'] === '/cosmos.bank.v1beta1.MsgSend') {
            return `${new Token({
                symbol: network.networkToken.symbol,
                amount: messages[0].amount[0].amount,
                precision: network.networkToken.precision
            }).string()}`;
        } else {
            return `0 ${network.networkToken.symbol}`;
        }
    }

    function getTxType(messages) {
        const type = txType[messages[0]['@type']];
        return type === undefined ? messages[0]['@type'] : $_(type);
    }
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('transactions')}</h5></div>
                    {#if !limit}
                        <h6 class="TW FS13">
                            <span class="TW FS13">
                                {$_('total_transactions')} :
                            </span>
                            {totalTxs}
                        </h6>
                    {:else}
                        <Link class="ShowMore" to="/txs">
                            {$_('show_more')}
                            <span class="NArrow">
                                <IconRightArrow />
                            </span>
                        </Link>
                    {/if}
                </div>
            </div>
            {#if onLoad && txs.length > 0}
                <div class="CardList">
                    <table class="CardTable TxPo">
                        <thead />
                        <tbody>
                            {#each txs as tx}
                                <tr>
                                    <td class="px-2 py-1 TextAL">
                                        <div class="TxIcon">
                                            <TxResultIcon code={tx.code} />
                                        </div>
                                        <p class="FS18 pT10 pL65">
                                            {getSendAmount(tx.tx.body.messages)}
                                            <span class="FS12 Color_Point pL4">
                                                {getTxType(tx.tx.body.messages)}
                                            </span>
                                        </p>
                                        <Link
                                            class="FS12 Color_Gray pL65 TxAddress"
                                            to="/account/{tx.sender}"
                                        >
                                            {tx.sender}
                                        </Link>
                                    </td>
                                    <td class="px-2 py-1 TextAR MoPosition">
                                        <p class="FS15 fwB Color_Dark">
                                            {getUpdateTime(tx.timestamp)}
                                        </p>
                                        <p class="FS12 fwR Color_Gray">
                                            <Link to="/tx/{tx.txhash}">
                                                <span
                                                    class="FS12 fwR Color_Gray"
                                                >
                                                    {truncateHex(tx.txhash)}
                                                </span>
                                            </Link>
                                        </p>
                                        <p class="FS12 fwR Color_Gray">
                                            {new Date(
                                                tx.timestamp
                                            ).toDateString()}
                                        </p>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>

                {#if !limit}
                    <Paginator
                        bind:page={txPage}
                        bind:pageLimit={txLimit}
                        bind:maxItem={totalTxs}
                    />
                {:else}{/if}
            {:else if onLoad && txs.length === 0}
                <NoData />
            {:else}
                <LoadingIndicator />
            {/if}
        </div>
    </div>
</section>
