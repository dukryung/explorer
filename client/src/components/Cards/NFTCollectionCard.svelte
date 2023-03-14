<script>
    import { onMount } from 'svelte';
    import { _ } from 'svelte-i18n';
    import { Link, navigate } from 'svelte-routing';
    import IconRightArrow from '../../assets/IconRightArrow.svelte';
    import NoData from 'Components/Labels/NoData';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import Paginator from 'Components/Footers/Paginator';

    // default value
    export let tokenLimit = 12;
    export let tokenPage = 0;
    export let col_id;
    export let address;
    // default params
    export let limit = false;
    // export let type = klaatoo.SUBSCRIBE;
    export let method = 'nft.listbycollection';

    const requestId = klaatoo.generateRequestId();

    let infos = [];
    let totalTokens = 0;
    let onLoad = false;
    let nft_name = '';
    onMount(() => {
        window.scroll({
            top: 0,
            left: 0,
            behavior: 'smooth'
        });
    });
    $: if (tokenPage !== undefined) {
        let params;
        params = [col_id, tokenLimit, tokenPage];
        if (method == 'nft.listbyowner' || method == 'nft.listbyissuer') {
            params = [address, tokenLimit, tokenPage];
        }
        if (method == 'nft.listbyowner') {
            nft_name = 'Owned ';
        } else if (method == 'nft.listbyissuer') {
            nft_name = 'Issued ';
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

    const randomHandler = id => {
        navigate(`/nft/${id}/`);
    };
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{nft_name}{$_('nfts')}</h5></div>
                    {#if !limit}
                        <h6 class="TW FS13">
                            <span class="TW FS13"> {$_('total_nfts')} </span>
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

            {#if onLoad && infos.length > 0}
            <div class="nftWrap">
                {#each infos as info}
                <div
                on:click={randomHandler(info.token.id)}
                class="nftCardWrap"
            >
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
                                <p class="FS15 fwB Color_Dark">
                                    {info.token.collection_id && info.token.collection_id.slice(0,-11)}
                                </p>
                            </div>
                        </div>
                    </div>
            </div>
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
