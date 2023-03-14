<script>
    import {_} from 'svelte-i18n';
    import {onMount} from 'svelte';
    import {Token} from '../../js/token';
    import NoData from '../Labels/NoData.svelte';
    import {Link} from 'svelte-routing';

    export let address;
    let delegations = [];
    let onLoad = false;

    onMount(() => {
      klaatoo.singleRequest({
        method: 'account.delegation',
        params: [
          address,
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          if (data !== null) {
            delegations = data;
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
                    <div><h5>{$_('delegation')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                {#if onLoad && delegations.length > 0}
                    <table class="CardTable AcInfo">
                        <colgroup>
                            <col class="LeftCol">
                            <col class="RightCol">
                        </colgroup>
                        <thead>
                        </thead>
                        <tbody>
                        {#each delegations as delegation}
                            <tr>
                                <td class="px-2 py-1 TextAL">
                                    <p class="FS15 fwM Color_Dark">
                                        <Link to="/validator/{delegation.delegation.validator_address}">
                                            {delegation.moniker}
                                        </Link>
                                    </p>
                                </td>
                                <td class="px-2 py-1 TextAR">
                                    <p class="FS15 fwM Color_Dark AcAd">
                                        {new Token({
                                            symbol: network.networkToken.symbol,
                                            amount: delegation.balance.amount,
                                            precision: network.networkToken.precision,
                                        }).string()}
                                    </p>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                {:else if onLoad && delegations.length === 0}
                    <NoData description={$_('no_delegation')}/>
                {:else}
                    <div class="LoadWrap">
                        <div id="loading"></div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>
