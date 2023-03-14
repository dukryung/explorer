<script>
    import { _ } from 'svelte-i18n';
    import Meta from 'svelte-meta';
    import { onMount, beforeUpdate } from 'svelte';
    import NFTokenSummaryCard from '../components/Cards/NFTokenSummaryCard.svelte';
    import TxListCard from '../components/Cards/TxListCard.svelte';
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import AccountTokenListCard from 'Components/Cards/AccountTokenListCard';

    export let token_id;

    const urlParams = new URLSearchParams(window.location.search);
    const hasId = urlParams.has('id');
    if (hasId) {
        token_id = urlParams.get('id');
    }

    const txLimit = 10;
    const txPage = 0;
</script>

<Meta title="{$_('meta_title')} | {$_('nftoken')} - {token_id}" />

<SearchCard title={$_('nftoken_details')} />

<NFTokenSummaryCard {token_id} />
<!--
<TxListCard
        type={klaatoo.SUBSCRIBE}
        method={'tx.address'}
        address={id}
        txLimit={txLimit}
        txPage={txPage}
/>
-->
