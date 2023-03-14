<script>
    import BlockSummaryCard from '../components/Cards/BlockSummaryCard.svelte';
    import {_} from 'svelte-i18n';
    import TxListCard from '../components/Cards/TxListCard.svelte';
    import Meta from 'svelte-meta';
    import SearchCard from '../components/Cards/SearchCard.svelte';

    export let height;
</script>

<Meta
        title="{$_('meta_title')} | {$_('block')} # {height}"
/>

<SearchCard
        title="{$_('block')} # {height}"
/>

<BlockSummaryCard
        height={height}
/>

<TxListCard
        type={klaatoo.SINGLE_REQUEST}
        method={'tx.height'}
        height={parseInt(height)}
/>
