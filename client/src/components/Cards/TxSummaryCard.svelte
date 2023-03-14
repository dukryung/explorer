<script>
    import {_} from 'svelte-i18n';
    import {onMount} from 'svelte';
    import {Link, navigate} from 'svelte-routing';
    import {getUpdateTime} from 'js';
    import TxResultText from '../Labels/TxResultText.svelte';
    import TxMessage from '../Messages/TxMessage.svelte';

    export let hash;
    let onLoad = false;
    let tx;

    onMount(() => {
      klaatoo.singleRequest({
        method: 'tx.hash',
        params: [
          hash,
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          tx = data;
          onLoad = true;
        },
        error: (error) => {
          console.error(error);
          navigate('/error', {replace: false});
        },
      });
    });
</script>


<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('tx_summary')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                <table class="CardTable TxDview">
                    <colgroup>
                        <col>
                    </colgroup>
                    <thead>
                    </thead>
                    <tbody>
                    {#if onLoad && tx}
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark ">
                                    {$_('tx_hash')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS13 fwM Color_Dark TxHa">
                                    {tx.txhash}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('tx_status')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    <TxResultText code={tx.code}/>
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('tx_height')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <Link class="FS15 fwM Color_Dark" to="/block/{tx.height}"
                                      replace={false}>{tx.height}
                                </Link>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('tx_time')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    {new Date(tx.timestamp).toLocaleString()} ({getUpdateTime(tx.timestamp)})
                                </p>
                            </td>
                        </tr>
                    {:else if onLoad && !tx}
                    {:else}
                        <div class="LoadWrap">
                            <div id="loading"></div>
                        </div>
                    {/if}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</section>

{#if tx}
    <TxMessage messages={tx.tx.body.messages}/>
{/if}