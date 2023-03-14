<script>
    import {onDestroy, onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import {Link, navigate} from 'svelte-routing';
    import {getUpdateTime} from 'js';
    import NoData from 'Components/Labels/NoData';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import Paginator from 'Components/Footers/Paginator';
    import IconRightArrow from '../../assets/IconRightArrow.svelte';

    // default value
    const requestId = klaatoo.generateRequestId();

    // default params
    export let limit = false;
    export let type = klaatoo.SUBSCRIBE;
    export let method = 'block.latest';
    export let blockLimit = 10;
    export let blockPage = 0;
    export let address = '';

    let blocks = [];
    let blockHeight = 0;
    let onLoad = false;

    onMount(() => {

    });

    $: if (blockPage !== undefined) {
      let params;
      if (method === 'block.latest') {
        params = [
          blockLimit,
          blockPage,
        ];
      } else {
        params = [
          blockLimit,
          blockPage,
          address,
        ];
      }
      klaatoo.request({
        type: type,
        method: method,
        params: params,
        id: requestId,
        success: (data) => {
          if (data !== null) {
            if (data.blocks !== null) blocks = data.blocks;
            blockHeight = data.total;
          }
          onLoad = true;
        },
        error: (error) => {
          console.error(error);
          navigate('/error', {replace: false});
        },
      });
    }

    onDestroy(() => {
      if (type === klaatoo.SUBSCRIBE) {
        klaatoo.unsubscribe({
          method: method,
          id: requestId,
        });
      }
    });
</script>
<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('blocks')}</h5></div>
                    {#if !limit}
                        <h6 class="TW FS13">
                            <span class="TW FS13">
                                {$_('total_blocks')} :
                            </span>
                            {blockHeight}
                        </h6>
                    {:else}
                        <Link class="ShowMore" to="/blocks">
                            {$_('show_more')}
                            <span class="NArrow">
                                <IconRightArrow/>
                            </span>
                        </Link>
                    {/if}
                </div>
            </div>
            {#if onLoad && blocks.length > 0}
                <div class="CardList">
                    <table class="CardTable BlPo">
                        <thead>
                        </thead>
                        <tbody>
                        {#each blocks as block}
                            <tr>
                                <td class="px-2 py-1 TextAL">
                                    <Link class="fs24 fwB Color_Dark" to="/block/{block.tm_block.header.height}">
                                        #{block.tm_block.header.height}
                                    </Link>
                                    <p class="FS12 fwR Color_Gray">
                                        <!-- Txs Count -->
                                        {$_('includes')} {block.data === undefined ? 0 : block.data.length} {$_('txs')}

                                        {$_('proposer')} :
                                        <Link class="FS12 fwM Color_Dark" to="/validator/{block.val_address}">
                                            {block.moniker}
                                        </Link>
                                    </p>
                                </td>
                                <td class="px-2 py-1 TextAR">
                                    <p class="FS15 fwB Color_Dark">
                                        {getUpdateTime(block.tm_block.header.time)}
                                    </p>
                                    <p class="FS12 fwR Color_Gray">
                                        {new Date(block.tm_block.header.time).toDateString()}
                                    </p>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                </div>
                {#if !limit}
                    <Paginator
                            bind:page={blockPage}
                            bind:pageLimit={blockLimit}
                            bind:maxItem={blockHeight}
                    />
                {:else}
                {/if}
            {:else if onLoad && blocks.length === 0}
                <NoData/>
            {:else}
                <LoadingIndicator/>
            {/if}
        </div>
    </div>
</section>
