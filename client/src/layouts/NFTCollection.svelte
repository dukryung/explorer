<script>
    import {_} from 'svelte-i18n';
    import Meta from 'svelte-meta';
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import NftCollectionCard from '../components/Cards/NFTCollectionCard.svelte';

    export let col_id;;
   
</script>

<Meta
        title="{$_('meta_title')} | {$_('nftokens')} - {col_id}"
/>

<SearchCard
        title={col_id}
/>

<NftCollectionCard col_id={col_id} />