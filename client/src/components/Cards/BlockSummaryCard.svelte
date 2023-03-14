<script>
    import {_} from 'svelte-i18n';
    import {onMount} from 'svelte';
    import {convertHashToString} from 'js';
    import {Link, navigate} from 'svelte-routing';
    export let height;
    let onLoad = false;
    let block;

    onMount(() => {
      klaatoo.singleRequest({
        method: 'block.height',
        params: [
          parseInt(height),
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          block = data;
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
                    <div><h5>{$_('block_summary')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                {#if onLoad && block}
                <table class="CardTable BlSu">
                    <colgroup>
                        <col class="LeftCol">
                        <col class="RightCol">
                    </colgroup>
                    <thead>
                    </thead>

                        <tbody>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('block_height')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    {block.tm_block.header.height}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('block_hash')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS13 fwM Color_Dark BlHa">
                                    {convertHashToString(block.tm_block.header.last_block_id.hash)}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('proposer')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    {block.moniker}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('transactions')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    {block.tm_block.data.txs !== undefined ? block.tm_block.data.txs.length : 0}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('block_time')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    {new Date(block.tm_block.header.time).toLocaleString()}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('diff_time')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    {block.diff_time/1000}{$_('second')}
                                </p>
                            </td>
                        </tr>
                        </tbody>

                </table>
                {:else if onLoad && !block}
                {:else}
                    <div class="LoadWrap">
                        <div id="loading"></div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>