<script>
    import {_} from 'svelte-i18n';
    import Meta from 'svelte-meta';
    import ValidatorSummaryCard from 'Components/Cards/ValidatorSummaryCard';
    import BlockListCard from 'Components/Cards/BlockListCard';
    import ValidatorUptimeCard from 'Components/Cards/ValidatorUptimeCard';
    import ValidatorDelegatorsCard from 'Components/Cards/ValidatorDelegatorsCard';
    import SearchCard from '../components/Cards/SearchCard.svelte';

    export let address;

    const blockLimit = 10;
    const blockPage = 0;
</script>

<Meta
        title="{$_('meta_title')} | {$_('validator')} - {address}"
/>

<SearchCard
        title={$_('validator_details')}
/>

<ValidatorSummaryCard
        address={address}
/>

<ValidatorDelegatorsCard
        address={address}
/>

<ValidatorUptimeCard
        address={address}
/>

<BlockListCard
        method={'block.proposed'}
        blockLimit={blockLimit}
        blockPage={blockPage}
        address={address}
/>