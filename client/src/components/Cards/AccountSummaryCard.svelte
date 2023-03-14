<script>
    import {_} from 'svelte-i18n';
    import {onMount} from 'svelte';
    import {Token} from 'js/token';
    import {navigate} from 'svelte-routing';

    export let address;
    let token = {};
    let onLoad = false;
    $:if (address) {
      klaatoo.singleRequest({
        method: 'account.balance',
        params: [
          address,
          network.networkToken.base,
        ],
        id: klaatoo.generateRequestId(),
        success: (data) => {
          token = data;
          onLoad = true;
        },
        error: (error) => {
          console.error(error);
        },
      });
    };
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('account_information')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                {#if onLoad && token}
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
                                <p class="FS15 fwM Color_Dark">{$_('address')}
                                </p>

                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark AcAd">{address}
                                </p>
                            </td>
                        </tr>
                        <tr>
                            <td class="px-2 py-1 TextAL">
                                <p class="FS15 fwM Color_Dark">
                                    {$_('network_token')}
                                </p>
                            </td>
                            <td class="px-2 py-1 TextAR">
                                <p class="FS20 fwM Color_Dark">
                                    {new Token({amount: token.coin.amount}).string()}
                                </p>
                            </td>
                        </tr>

                        </tbody>
                    </table>
                {:else if onLoad && !token}
                {:else}
                    <div class="LoadWrap">
                        <div id="loading"></div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>
