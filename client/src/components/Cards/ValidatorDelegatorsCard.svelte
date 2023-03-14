<script>
    import {_} from 'svelte-i18n';
    import {onMount} from 'svelte';
    import {Link, navigate} from 'svelte-routing';
    import {Token} from '../../js/token';
    import Paginator from 'Components/Footers/Paginator';

    export let address;
    let delegatorPage = 0;
    let delegatorLimit = 10;
    let totalDelegations = 0;

    let delegations = [];
    let onLoad = false;
    onMount(() => {
      klaatoo.singleRequest({
        method: 'validator.delegation',
        params: [
          address,
          delegatorLimit,
          delegatorPage,
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          delegations = data.delegation;
          totalDelegations = data.total;
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
                    <div><h5>{$_('delegators')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                {#if onLoad}
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

                                        <Link to="/account/{delegation.delegation.delegator_address}">
                                            {delegation.delegation.delegator_address}
                                        </Link>
                                    </p>
                                </td>
                                <td class="px-2 py-1 TextAR TokenMo">
                                    <p class="FS15 fwB Color_Dark">
                                        {new Token({
                                            symbol: network.networkToken.symbol,
                                            denom: delegation.balance.denom,
                                            amount: delegation.balance.amount,
                                            precision: network.networkToken.precision,
                                        }).string()}
                                    </p>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                    <Paginator
                            bind:pageLimit={delegatorLimit}
                            bind:page={delegatorPage}
                            bind:limit={totalDelegations}
                    />
                {:else}
                    <div class="LoadWrap">
                        <div id="loading"></div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>