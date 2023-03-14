<script>
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import NFTokenListCard from '../components/Cards/NFTokenListCard.svelte';
    import {_} from 'svelte-i18n';
    import Meta from 'svelte-meta';
</script>

<Meta
        title="{$_('meta_title')} | {$_('nftokens')}"
        description="{$_('meta_description')}"
/>

<SearchCard title={$_('nftokens')}/>


<NFTokenListCard/>

