<script>
    import TxListCard from '../components/Cards/TxListCard.svelte';
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import {_} from 'svelte-i18n';
    import Meta from 'svelte-meta';
</script>

<Meta
        title="{$_('meta_title')} | {$_('transactions')}"
/>

<SearchCard
        title={$_('transactions')}
/>

<TxListCard/>