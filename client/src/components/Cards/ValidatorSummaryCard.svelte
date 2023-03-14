<script>
    import {onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import {Token} from 'js/token';
    import {navigate} from 'svelte-routing';
    import {fixedCommisionRate, getVotingPower, stats} from 'js';

    export let address;

    let onLoad = false;
    let validator = {};
    onMount(() => {
      klaatoo.singleRequest({
        method: 'validator.address',
        params: [
          address,
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          validator = data;
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
                    <div><h5>{$_('validator_information')}</h5></div>
                    {#if onLoad}
                        {#if validator.detail.jailed !== undefined}
                            <p style="background-color: red; color: white">{$_('jailed')}</p>
                        {/if}
                    {/if}
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
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">{$_('val_address')}</p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark AcAd">
                                    {validator.val_address}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">{$_('moniker')}</p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark AcAd">
                                    {validator.moniker}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">{$_('rank')}</p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark AcAd">
                                    #{validator.rank}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">{$_('voting_power')}</p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark AcAd">
                                    {getVotingPower($stats.total_bonded_tokens, validator.detail.tokens)}%
                                    ({new Token({amount: validator.detail.tokens}).string()})
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">{$_('commision_rate')}</p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark AcAd">
                                    {fixedCommisionRate(validator.detail.commission.commission_rates.rate)}
                                </p>
                            </td>
                        </tr>
                        {#if validator.detail.description.details !== undefined}
                            <tr>
                                <td class="px-2 py-1 TextAL">
                                    <p class="FS15 fwM Color_Dark">{$_('description_details')}</p>
                                </td>
                                <td class="px-2 py-1 TextAR">
                                    <p class="FS15 fwM Color_Dark AcAd">
                                        {validator.detail.description.details}
                                    </p>
                                </td>
                            </tr>
                        {/if}

                        {#if validator.detail.description.website !== undefined}
                            <tr>
                                <td class="px-2 py-1 TextAL">
                                    <p class="FS15 fwM Color_Dark">{$_('website')}</p>
                                </td>
                                <td class="px-2 py-1 TextAR">
                                    <p class="FS15 fwM Color_Dark AcAd">
                                        <a href={validator.detail.description.website}>
                                            {validator.detail.description.website}
                                        </a>
                                    </p>
                                </td>
                            </tr>
                        {/if}
                        </tbody>
                    </table>
                {:else if onLoad}
                {:else}
                    <div class="LoadWrap">
                        <div id="loading"></div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>
