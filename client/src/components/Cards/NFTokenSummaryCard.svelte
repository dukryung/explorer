<script>
    import { _ } from 'svelte-i18n';
    import { onMount, beforeUpdate } from 'svelte';
    import { navigate } from 'svelte-routing';

    export let token_id;
    let token_info = {};
    let onLoad = false;

    onMount(() => {
        window.scroll({
            top: 0,
            left: 0,
            behavior: 'smooth'
        });
        if (token_id) {
            klaatoo.singleRequest({
                method: 'nft.info',
                params: [token_id],
                id: klaatoo.generateRequestId(),
                success: data => {
                    token_info = data;
                    onLoad = true;
                },
                error: error => {
                    console.error(error);
                    navigate('/error', { replace: false });
                }
            });
        }
    });
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('nftoken_information')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                {#if onLoad && token_info}
                    <div class="nft_inf_wrap">
                        <div class="nft_inf_img">
                            <img
                                src={token_info.token.preview_url}
                                alt="No preview"
                            />
                        </div>
                        <div class="nft_inf_main">
                            <div>
                                <h1 class="fs24 fwB">
                                    {token_info.token.name}
                                </h1>
                            </div>
                            <div>
                                <p class="FS16 fwM Color_Dark textBreak">
                                    Owned by <a
                                        class="tablePoint linkhover"
                                        href="/account/{token_info.owner_address}"
                                        >{token_info.owner_address}</a
                                    >
                                </p>
                            </div>
                            <div class="desc_wrap">
                                <div class="desc_container">
                                    <p class="nftlabel">
                                        {$_('token_id')}
                                    </p>
                                    <p class="FS16 fwB Color_Dark">
                                        {token_id}
                                    </p>
                                </div>
                                <div class="desc_container">
                                    <p class="nftlabel">
                                        {$_('collection_id')}
                                    </p>
                                    <p class="FS16 fwB Color_Dark">
                                        <a
                                            href="/nfts/collection/{token_info.token
                                                .collection_id}/"
                                                class="linkhover"
                                        >
                                            {token_info.token.collection_id}
                                        </a>
                                    </p>
                                </div>
                                <div class="desc_container abc">
                                    <p class="nftlabel">
                                        {$_('issuer_address')}
                                    </p>
                                    <p class="FS16 fwB Color_Dark">
                                        <a
                                            class="tablePoint linkhover textBreak"
                                            href="/account/{token_info.issuer_address}"
                                            >{token_info.issuer_address}</a
                                        >
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                {:else if onLoad && !token_info}
                    <p class="FS15 fwM Color_Dark">{$_('nftoken_not_found')}</p>
                {:else}
                    <div class="LoadWrap">
                        <div id="loading" />
                    </div>
                {/if}
            </div>
        </div>
    </div>
</section>
