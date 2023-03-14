<script>
    import { onDestroy, onMount } from 'svelte';
    import { _ } from 'svelte-i18n';
    import { Link, navigate } from 'svelte-routing';
    import IconRightArrow from '../../assets/IconRightArrow.svelte';
    import NoData from 'Components/Labels/NoData';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import Paginator from 'Components/Footers/Paginator';

    // default value
    export let tokenLimit = 12;
    export let tokenPage = 0;

    // default params
    export let limit = false;
    // export let type = klaatoo.SUBSCRIBE;
    export let method = 'nft.list';

    const requestId = klaatoo.generateRequestId();

    let infos = [];
    let totalTokens = 0;
    let onLoad = false;
    let selected = 'random';

    $: if (tokenPage !== undefined) {
        let params;
        params = [tokenLimit, tokenPage];
        if (selected == 'random') {
            method = 'nft.list';
            onLoad = false;
        } else if (selected == 'collection') {
            method = 'nft.listcollections';
            onLoad = false;
        }

        klaatoo.singleRequest({
            method: method,
            params: params,
            id: requestId,
            success: data => {
                if (data !== null) {
                    if (data.infos !== null) {
                        infos = data.infos;
                    }
                    totalTokens = data.total;
                }

                onLoad = true;
            },
            error: error => {
                console.error(error);
                navigate('/error', { replace: false });
            }
        });
    }

    // onDestroy(() => {
    //   if (type === klaatoo.SUBSCRIBE) {
    //     klaatoo.unsubscribe({
    //       method: method,
    //       id: requestId,
    //     });
    //   }
    // });
    const randomHandler = id => {
        navigate(`/nft/${id}/`);
    };
    const collectionHandler = id => {
        navigate(`/nfts/collection/${id}`);
    };
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('nfts')}</h5></div>
                    {#if !limit}
                        <h6 class="TW FS13">
                            <span class="TW FS13">
                                {selected == 'random'
                                    ? `${$_('total_nfts')}`
                                    : `${$_('total_collection')}`}
                            </span>
                            {totalTokens}
                        </h6>
                    {:else}
                        <Link class="ShowMore" to="/nfts">
                            {$_('show_more')}
                            <span class="NArrow">
                                <IconRightArrow />
                            </span>
                        </Link>
                    {/if}
                </div>
            </div>

            <form class="nft_segmented_button">
                <div class="segmented-control">
                    <div
                        on:click={() => {
                            selected = 'random';
                            tokenPage = 0;

                        }}
                        class="segmented-control-btn"
                    >
                        <input
                            type="radio"
                            id="random"
                            name="collection_random"
                            checked
                        />
                        <label for="random">{$_('random_nft')}</label>
                    </div>
                    <div
                        on:click={() => {
                            selected = 'collection';
                            tokenPage = 0;
                        }}
                        class="segmented-control-btn"
                    >
                        <input
                            type="radio"
                            id="collection"
                            name="collection_random"
                        />
                        <label for="collection">{$_('nft_collection')}</label>
                    </div>
                </div>
            </form>

            {#if onLoad && infos.length > 0}
            <div class="nftWrap">
                {#each infos as info}
                    {#if selected == 'random'}
                    <div
                    on:click={randomHandler(info.token.id)}
                    class="nftCardWrap">
                        <img
                            src={info.token.preview_url}
                            alt="nft preview"
                        />
                        <div class="nftTextWrap">
                            
                            <p class="infoName">
                                {info.token.name}
                            </p>
                            <div class="infoTextWrap">
                                <div class="infoLeft">
                                    <p class="infoLabel">ID</p>
                                    <p class="FS15 fwB Color_Dark elli">
                                        {info.token.id}
                                    </p>
                                </div>

                                <div class="infoRight">
                                    <p class="infoLabel">Collection</p>
                                    <p class="FS15 fwB Color_Dark elli">
                                        {info.token.collection_id && info.token.collection_id.slice(0,-11)}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                    {:else if selected == 'collection'}
                    <div
                    on:click={collectionHandler(info.token.id)}
                    class="nftCardWrap">
                        <img
                        src={info.token.preview_url}
                        alt="nft preview"/>
                        <div class="nftTextWrap">
                            <p class="infoName">
                                {info.token.id}
                            </p>

                            <div class="infoTextWrap">
                                <div class="TextAL">
                                    <p class="infoLabel">Name</p>
                                    <p class="infotext elli">
                                        {info.token.name}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                    {/if}
                {/each}
            </div>
                {#if !limit}
                    <Paginator
                        bind:page={tokenPage}
                        bind:pageLimit={tokenLimit}
                        bind:maxItem={totalTokens}
                    />
                {/if}
            {:else if onLoad && infos.length === 0}
                <NoData description="No NFTs" />
            {:else}
                <LoadingIndicator />
            {/if}
        </div>
    </div>
</section>
