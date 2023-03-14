<script>
    import {_} from 'svelte-i18n';
    import {onDestroy} from 'svelte';
    import Paginator from 'Components/Footers/Paginator';
    import {Link, navigate} from 'svelte-routing';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import NoData from 'Components/Labels/NoData';
    import {Token} from 'js/token';
    import {fixedCommisionRate, getVotingPower, stats} from 'js';

    // default value
    const requestId = klaatoo.generateRequestId();

    // default params
    export let limit = false;
    export let type = klaatoo.SUBSCRIBE;
    export let method = 'validator.list';
    export let validatorLimit = 10;
    export let validatorPage = 0;

    let validators = [];
    let validatorTotal = 0;
    let onLoad = false;

    $: if (validatorPage !== undefined) {
      let params;
      params = [
        validatorLimit,
        validatorPage,
      ];
      klaatoo.request({
        type: type,
        method: method,
        params: params,
        id: requestId,
        success: (data) => {
          if (data !== null) {
            if (data.validators !== null) validators = data.validators;
            validatorTotal = data.total;
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
                <div class="CardTitle Validator">
                    <h5>{$_('validators')}</h5>
                    <h6 class="TW FS13 BondedTokens">
                        <span class="TW FS13">
                            {$_('total_bonded_tokens')} :
                        </span>
                        {new Token({amount: $stats.total_bonded_tokens}).string()}
                    </h6>
                </div>
            </div>
            {#if onLoad && validators.length > 0}
                <div class="CardList">
                        <table class="CardTable VaPo">
                            <tbody>
                            {#each validators as validator, i}
                                <tr>
                                    <td class="px-2 py-1 TextAL">
                                            <Link class="Color_Dark fwB" to="/validator/{validator.detail.operator_address}"><span>{validator.rank}.</span>
                                                {validator.detail.description.moniker}
                                            </Link>
                                        <p class="FS15 fwM Color_Dark">
                                            {$_('delegated_tokens')} : {new Token({amount: validator.detail.tokens}).string()}
                                        </p>
                                    </td>
                                    <td class="px-2 py-1 TextAR">
                                        <p class="FS15 fwM Color_Dark">
                                            {$_('voting_power')} : {getVotingPower($stats.total_bonded_tokens, validator.detail.tokens)}%
                                        </p>
                                        <p class="FS13 fwM Color_Gray">
                                            {$_('commission_rate')} : {fixedCommisionRate(validator.detail.commission.commission_rates.rate)}
                                        </p>
                                    </td>
                                </tr>
                            {/each}
                            </tbody>
                        </table>
                    {#if !limit}
                        <Paginator
                                bind:page={validatorPage}
                                bind:pageLimit={validatorLimit}
                                bind:maxItem={validatorTotal}
                        />
                    {:else}
                    {/if}
                </div>
            {:else if onLoad && validators.length === 0}
                <NoData/>
            {:else}
                <LoadingIndicator/>
            {/if}
        </div>
    </div>
</section>