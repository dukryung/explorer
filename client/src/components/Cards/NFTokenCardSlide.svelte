<script>
    import { _ } from 'svelte-i18n';
    import { Link, navigate } from 'svelte-routing';
    import IconRightArrow from '../../assets/IconRightArrow.svelte';
    import NoData from 'Components/Labels/NoData';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import { onMount, afterUpdate } from 'svelte';

    export let tokenLimit = 10;
    export let tokenPage = 0;
    // export let type = klaatoo.SUBSCRIBE;
    export let method = 'nft.list';

    const requestId = klaatoo.generateRequestId();

    let infos = [];
    let totalTokens = 0;
    let onLoad = false;
    let selected = 'random';
    let randomArray = [];

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
                        randomArray = [...data.infos].sort(()=> Math.random()-0.5);

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

    let currentIdx;
    let slideCount = 0;

    
    afterUpdate(()=>{
        slideCount = totalTokens;
        tokenLimit = slideCount;
    })
    
    const randomHandler = id => {
        navigate(`/nft/${id}/`);
    };
    const collectionHandler = id => {
        navigate(`/nfts/collection/${id}`);
    };


    //slider
    currentIdx = 0;
    let positionLeft = 0;
    let windowWidth = window.innerWidth
    let percent = windowWidth <= 450 ?  100 : windowWidth <= 600 ? 50 : windowWidth <= 820 ? 33.33 : windowWidth <= 960  ?  25 : 20 ;
    let slideWrap;
    let slideWidth;
    let cardWidth = (slideWidth * percent/100);
    
    afterUpdate(()=>{
        slideWidth = slideWrap && slideWrap.getBoundingClientRect().width
        cardWidth =  (slideWidth * percent/100);
    })
    
    window.onresize = function(){
        windowWidth = window.innerWidth
        slideWidth = slideWrap && slideWrap.getBoundingClientRect().width
        percent = windowWidth <= 450 ?  100 : windowWidth <= 600 ? 50 : windowWidth <= 820 ? 33.33 : windowWidth <= 960  ?  25 : 20 ;
        cardWidth =  (slideWidth * percent/100);
        positionLeft = cardWidth * currentIdx ;
    }

    const moveSlider = (num) => {
        positionLeft = (num * cardWidth); 
        currentIdx = num;
};

    const next = (e) => {
        currentIdx !== slideCount - parseInt(100/percent) && moveSlider(currentIdx + 1)
        e.stopPropagation()
    };
	
    const prev = (e) => {
        currentIdx !== 0 && moveSlider(currentIdx - 1)
        e.stopPropagation()
    };


    function dim(state) {
    return state === false ? 'buttonDim' : '';
    }

    $: preBtn = currentIdx !== 0 ?  true :  false
    $: nextBtn = currentIdx !== (slideCount - parseInt(100/percent)) ? true : false


</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('nfts')}</h5></div>
                    <Link class="ShowMore" to="/nfts">
                        {$_('show_more')}
                        <span class="NArrow">
                            <IconRightArrow />
                        </span>
                    </Link>
                </div>
            </div>

            <form class="nft_segmented_button">
                <div class="segmented-control">
                    <div
                        on:click={() => {
                            selected = 'random';
                            positionLeft = 0;
                            currentIdx = 0;
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
                            positionLeft = 0;
                            currentIdx = 0;
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
                <img on:click={prev} class="leftArrow {dim(preBtn)}" alt="leftArrow" src="/images/left-arrow.png"/>
                <img on:click={next} class="rightArrow {dim(nextBtn)}" alt="rightArrow" src="/images/right-arrow.png"/>
            <div bind:this={slideWrap} class="SlideAllWrap">
                <div class="nftSlideWrap" style="left: -{positionLeft}px;">
                    {#if selected == 'random'}
                    {#each randomArray as info}
                    <div class="cardWrap">
                        <div
                            on:click={randomHandler(info.token.id)}
                            class="nftSlideCardWrap" id="randomCard"
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
                                            <p class="FS15 fwB Color_Dark elli">
                                                {info.token.collection_id && info.token.collection_id.slice(0,-11)}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                        </div>
                    </div>
                    {/each}
                    {:else if selected == 'collection'}
                    {#each infos as info}
                    <div class="cardWrap">
                        <div
                                on:click={collectionHandler(info.token.id)}
                                class="nftSlideCardWrap" id="collectionCard"
                            >
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
                    </div>
                    {/each}
                    {/if}
                </div>
            </div>
            {:else if onLoad && infos.length === 0}
                <NoData description="No NFTs" />
            {:else} 
            <div class="LoadWrap">
                <div id="loading"></div>
            </div>
            {/if}
        </div>
    </div>
</section>