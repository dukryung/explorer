<script>
    import BlockListCard from 'Components/Cards/BlockListCard.svelte';
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import {_} from 'svelte-i18n';
    import Meta from 'svelte-meta';
</script>

<Meta
        title="{$_('meta_title')} | {$_('blocks')}"
/>

<SearchCard
        title={$_('blocks')}
/>

<BlockListCard/>
