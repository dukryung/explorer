<script>
    import Meta from 'svelte-meta';
    import {_} from 'svelte-i18n';
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import BlockListCard from '../components/Cards/BlockListCard.svelte';
    import AvgBlockTimeCard from '../components/Cards/Stats/AvgBlockTimeCard.svelte';
    import BlockHeightCard from '../components/Cards/Stats/BlockHeightCard.svelte';
    import TotalTxCard from '../components/Cards/Stats/TotalTxCard.svelte';
    import TotalValidatorsCard from '../components/Cards/Stats/TotalValidatorsCard.svelte';
    import TxHistoryCard from '../components/Cards/Stats/TxHistoryCard.svelte';
    import TxListCard from '../components/Cards/TxListCard.svelte';
    import NfTokenListCard from '../components/Cards/NFTokenListCard.svelte';
    import NfTokenCardSlide from '../components/Cards/NFTokenCardSlide.svelte';

    const txLimit = 5;
    const txPage = 0;

    const blockLimit = 5;
    const blockPage = 0;

    const tokenPage = 0;
    let tokenLimit = window.innerWidth <= 960 ? 6 : 8 ;

    window.onresize = function(){
        tokenLimit = window.innerWidth <= 960 ? 6 : 8;
    }
    
</script>

<Meta
        title="{$_('meta_title')}"
        description="{$_('meta_description')}"
/>

<SearchCard
        subTitle={$_('welcome')}
        title={$_('chain_name')}
/>

<div class="TxRow">
    <TxHistoryCard/>
</div>

<div class="CardSection">
    <div class="CardGrid">
        <BlockHeightCard/>
        <AvgBlockTimeCard/>
        <TotalTxCard/>
        <TotalValidatorsCard/>
    </div>
</div>

<div class="CardRow">
    <div class="CardRowWrap">
        <BlockListCard
                limit={true}
                blockLimit={blockLimit}
                blockPage={blockPage}
        />
        <TxListCard
                limit={true}
                txLimit={txLimit}
                txPage={txPage}
        />
    </div>
</div>

<!-- card slide version -->
<NfTokenCardSlide/>

<!-- card version -->
<!-- <NfTokenListCard limit={true} tokenLimit={tokenLimit} tokenPage={tokenPage}/> -->
