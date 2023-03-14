<script>
    import {_} from 'svelte-i18n';
    import {onMount} from 'svelte';
    import NoData from '../Labels/NoData.svelte';
    import {Link} from 'svelte-routing';
    import {Token} from 'js/token';
    import IconRightArrow from '../../assets/IconRightArrow.svelte';

    export let address;
    let redelegations = [];
    let onLoad = false;


    onMount(() => {
      klaatoo.singleRequest({
        method: 'account.redelegation',
        params: [
          address,
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          if (data !== null) {
            redelegations = data;
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
                    <div><h5>{$_('redelegation')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                {#if onLoad && redelegations.length > 0}
                    <table class="CardTable AcInfo">
                        <colgroup>
                            <col class="LeftCol">
                            <col class="RightCol">
                        </colgroup>
                        <thead>
                        </thead>
                        <tbody>
                        {#each redelegations as redelegation}
                            <tr>
                                <td class="px-2 py-1 TextAL">
                                    <p class="FS15 fwM Color_Dark">
                                        <Link to="/validator/{redelegation.redelegation.validator_src_address}">
                                            {redelegation.moniker_src}
                                        </Link>
                                        <IconRightArrow/>
                                        <Link to="/validator/{redelegation.redelegation.validator_dst_address}">
                                            {redelegation.moniker_dst}
                                        </Link>
                                    </p>
                                </td>
                                <td class="px-2 py-1 TextAR">
                                    <p class="FS15 fwM Color_Dark AcAd">{
                                        new Token({amount: redelegation.entries[0].balance}).string()}
                                    </p>
                                    <p class="FS15 fwM Color_Dark AcAd">
                                        {new Date(redelegation.entries[0].redelegation_entry.completion_time).toDateString()}
                                    </p>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                {:else if onLoad && redelegations.length === 0}
                    <NoData description={$_('no_redelegation')}/>
                {:else}
                    <div class="LoadWrap">
                        <div id="loading"></div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>
