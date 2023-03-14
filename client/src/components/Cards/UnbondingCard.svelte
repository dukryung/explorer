<script>
    import {_} from 'svelte-i18n';
    import {onMount} from 'svelte';
    import NoData from '../Labels/NoData.svelte';
    import {Link} from 'svelte-routing';
    import {Token} from 'js/token';

    export let address;
    export let test;
    let unbonding = [];
    let onLoad = false;


    onMount(() => {
      klaatoo.singleRequest({
        method: 'account.unbonding',
        params: [
          address,
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          if (data !== null) {
            unbonding = data;
          }
          onLoad = true;
        },
        error: (error) => {
          console.error(error);
        },
      });
    });
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('unbonding')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                {#if onLoad && unbonding.length > 0}
                    <table class="CardTable AcInfo">
                        <colgroup>
                            <col class="LeftCol">
                            <col class="RightCol">
                        </colgroup>
                        <thead>
                        </thead>
                        <tbody>
                        {#each unbonding as unbonding}
                            <tr>
                                <td class="px-2 py-1 TextAL">
                                    <p class="FS15 fwM Color_Dark">
                                        <Link to="/validator/{unbonding.unbonding.validator_address}">
                                            {unbonding.moniker}
                                        </Link>
                                    </p>
                                </td>
                                <td class="px-2 py-1 TextAR">
                                    <p class="FS15 fwM Color_Dark AcAd">
                                        {new Token({amount: unbonding.unbonding.entries[0].balance}).string()}
                                    </p>
                                    <p class="FS15 fwM Color_Dark AcAd">
                                        {new Date(unbonding.unbonding.entries[0].completion_time).toDateString()}
                                    </p>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                {:else if onLoad && unbonding.length === 0}
                    <NoData description={$_('no_unbonding')}/>
                {:else}
                    <div class="LoadWrap">
                        <div id="loading"></div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>
