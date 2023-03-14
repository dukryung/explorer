<script>
    import Meta from 'svelte-meta';
    import {_} from 'svelte-i18n';
    import SearchCard from 'Components/Cards/SearchCard';
    import ValidatorListCard from 'Components/Cards/ValidatorListCard';
</script>

<Meta
        title="{$_('meta_title')} | {$_('validators')}">
</Meta>

<SearchCard
        title={$_('validators')}
/>

<ValidatorListCard/>
