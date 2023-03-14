<script>
    import {_} from 'svelte-i18n';
    import IconSearch from '../../assets/IconSearch.svelte';
    import TokenItem from 'Components/TableItems/TokenItem';
    import {navigate} from 'svelte-routing';

    export let address;

    let searchToken = '';

    let tokens = [];
    let onLoad = false;
    klaatoo.singleRequest({
      method: 'account.allbalances',
      params: [
        address,
        network.networkToken.base,
      ],
      id: klaatoo.generateRequestId(),
      success: (data) => {
        tokens = data;
        onLoad = true;
      },
      error: (error) => {
        console.error(error);
        navigate('/error', {replace: false});
      },
    });
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTokenTitle">
                    <div class="TokenTitle">
                        <h5>
                            {$_('account_token')}
                        </h5>
                    </div>
                    <div class="searchBarToken">
                        <div id="searchToken">
                            <input id="input" placeholder="{$_('token_search_placeholder')}" bind:value={searchToken}>
                            <a>
                                <IconSearch/>
                            </a>
                        </div>
                    </div>

                </div>
            </div>
            <div class="CardList ScrollY">
                <table class="CardTable TokenPo">
                    <thead>
                    </thead>
                    <tbody>
                    {#if onLoad}
                        {#each tokens as token}
                            {#if searchToken === ''}
                                <TokenItem token={token}/>
                            {:else if searchToken !== '' && token.symbol.startsWith(searchToken.toUpperCase())}
                                <TokenItem token={token}/>
                            {/if}
                        {/each}
                    {:else}
                        <div class="LoadWrap">
                            <div id="loading"></div>
                        </div>
                    {/if}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</section>

