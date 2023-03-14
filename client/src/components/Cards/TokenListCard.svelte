<script>
    import {onDestroy, onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import {Link, navigate} from 'svelte-routing';
    import IconRightArrow from '../../assets/IconRightArrow.svelte';
    import {Token} from 'js/token';
    import NoData from 'Components/Labels/NoData';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import Paginator from 'Components/Footers/Paginator';

    // default value
    export let tokenLimit = 10;
    export let tokenPage = 0;

    // default params
    export let limit = false;
    export let type = klaatoo.SUBSCRIBE;
    export let method = 'token.list';

    const requestId = klaatoo.generateRequestId();

    let tokens = [];
    let totalTokens = 0;
    let onLoad = false;

    onMount(() => {
    });

    $: if (tokenPage !== undefined) {
      let params;
      params = [
        tokenLimit,
        tokenPage,
      ];
      klaatoo.request({
        type: type,
        method: method,
        params: params,
        id: requestId,
        success: (data) => {
          if (data !== null) {
            if (data.tokens !== null) {
              tokens = data.tokens;
            }
            totalTokens = data.total;
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
                    <div><h5>{$_('tokens')}</h5></div>
                    {#if !limit}
                        <h6 class="TW FS13">
                            <span class="TW FS13">
                                {$_('total_tokens')} :
                            </span>
                            {totalTokens}
                        </h6>
                    {:else}
                        <Link class="ShowMore" to="/tokens">
                            {$_('show_more')}
                            <span class="NArrow">
                                <IconRightArrow/>
                            </span>
                        </Link>
                    {/if}
                </div>
            </div>
            {#if onLoad && tokens.length > 0}
                <div class="CardList">
                    <table class="CardTable TokenPo">
                        <thead>
                        </thead>
                        <tbody>
                        {#each tokens as token}
                            <tr>
                                <td class="px-2 py-1 TextAL">
                                    <div class="TokenIcon">
                                        <img alt="icon" src="images/noIcon.png">
                                    </div>
                                    <p class="FS18 pT10 pL65">
                                        {token.description}
                                    </p>
                                    <p class="FS12 Color_Gray pL65 TokenDescription">
                                        Symbol: {token.symbol}, Precision: {token.precision}, Denom: {token.coin.denom}
                                    </p>
                                </td>
                                <td class="px-2 py-1 TextAR TokenMo">
                                    <p class="FS15 fwB Color_Dark">{new Token({
                                        symbol: token.symbol,
                                        amount: token.coin.amount,
                                        precision: token.precision,
                                    }).string()}
                                    </p>
                                    <Link class="FS12 Color_Gray pL65 TokenAddress" to="/account/{token.owner_address}">
                                        {token.owner_address}
                                    </Link>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                </div>

                {#if !limit}
                    <Paginator
                            bind:page={tokenPage}
                            bind:pageLimit={tokenLimit}
                            bind:maxItem={totalTokens}
                    />
                {/if}
            {:else if onLoad && tokens.length === 0}
                <NoData
                        description="No Tokens"
                />
            {:else}
                <LoadingIndicator/>
            {/if}
        </div>
    </div>
</section>