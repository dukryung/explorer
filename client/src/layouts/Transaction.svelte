<script>
    import TxSummaryCard from '../components/Cards/TxSummaryCard.svelte';
    import {_} from 'svelte-i18n';
    import Meta from 'svelte-meta';
    import SearchCard from '../components/Cards/SearchCard.svelte';

    export let hash;
</script>

<Meta
        title="{$_('tx_details')} {hash}"
/>

<SearchCard
        title="{$_('tx_details')}"
/>

<TxSummaryCard
        hash={hash}
/>
